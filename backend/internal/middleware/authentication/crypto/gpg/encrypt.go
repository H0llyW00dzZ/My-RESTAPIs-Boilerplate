// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// EncryptFile encrypts the given file using the provided GPG/OpenPGP public key.
//
// Note: Performance may depends on CPU and disk speed. However In K8s, this can be challenging for HPA
// because the encryption writes to disk, not to an object that can be easily stored elsewhere (e.g., storage mechanisms by Fiber, buckets, etc.).
func (e *Encryptor) EncryptFile(inputFile, outputFile string) (err error) {
	// Create a key ring from the public key
	keyRing, err := e.createKeyRing()
	if err != nil {
		return err
	}

	// Open the input file
	inFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	// Create the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() {
		if cerr := outFile.Close(); cerr != nil || err != nil {
			os.Remove(outputFile)
		}
	}()

	// Create metadata (header) for the encryption
	metadata := &crypto.PlainMessageMetadata{
		IsBinary: true,
		Filename: inFile.Name(),
		ModTime:  crypto.GetUnixTime(),
	}

	// Create a writer for the encrypted output
	encryptWriter, err := keyRing.EncryptStreamWithCompression(outFile, metadata, nil)
	if err != nil {
		return fmt.Errorf("failed to create encryption stream: %w", err)
	}
	defer encryptWriter.Close()

	// Stream the data
	//
	// This differs from "EncryptStream" which to a object because it stream writes directly to a file, not an object.
	buf := make([]byte, 4096) // Buffer size of 4KB
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from input file: %w", err)
		}
		if n == 0 {
			break
		}

		if _, err := encryptWriter.Write(buf[:n]); err != nil {
			return fmt.Errorf("failed to write encrypted data: %w", err)
		}
	}

	return nil
}

// EncryptStream (Object) encrypts data from an input stream and writes it to an output stream using the Encryptor's public key.
// This method is efficient for sending data over a network (e.g., TCP not only HTTP or GRPC whatever it is) or writing to a file.
//
// For Example Server-Side (E2E Logic transmission over a network):
//
//	func main() {
//		// Listen on TCP port 8080
//		listener, err := net.Listen("tcp", ":8080")
//		if err != nil {
//			log.Fatalf("Failed to listen on port 8080: %v", err)
//		}
//		defer listener.Close()
//
//		log.Println("Server listening on port 8080")
//
//		for {
//			conn, err := listener.Accept()
//			if err != nil {
//				log.Printf("Failed to accept connection: %v", err)
//				continue
//			}
//
//			go handleConnection(conn)
//		}
//	}
//
//	func handleConnection(conn net.Conn) {
//		defer func() {
//			if err := conn.Close(); err != nil {
//				log.Printf("Failed to close connection: %v", err)
//			}
//		}()
//
//		// Create a new Encryptor with the public key
//		encryptor, err := gpg.NewEncryptor([]string{"your-armored-public-key"})
//		if err != nil {
//			log.Printf("Failed to create encryptor: %v", err)
//			return
//		}
//
//		// Prepare the data to be encrypted
//		inputData := []byte("This is the data to be encrypted.")
//		inputBuffer := bytes.NewReader(inputData)
//
//		// Encrypt the data and send it over the network connection
//		if err := encryptor.EncryptStream(inputBuffer, conn); err != nil {
//			log.Printf("Failed to encrypt and send data: %v", err)
//			return
//		}
//
//			log.Println("Data encrypted and sent successfully")
//		}
//
// For Example Client-Side (E2E Logic transmission over a network):
//
//	func main() {
//		// Connect to the server
//		conn, err := net.Dial("tcp", "localhost:8080")
//		if err != nil {
//			log.Fatalf("Failed to connect to server: %v", err)
//		}
//		defer conn.Close()
//
//		// Read the encrypted data from the connection
//		var encryptedData bytes.Buffer
//		if _, err := io.Copy(&encryptedData, conn); err != nil {
//			log.Fatalf("Failed to read encrypted data: %v", err)
//		}
//
//		// Decrypt the data
//		decryptedData, err := decryptData(encryptedData.Bytes(), "your-armored-private-key", "your-passphrase")
//		if err != nil {
//			log.Fatalf("Failed to decrypt data: %v", err)
//		}
//
//		fmt.Println("Decrypted data:", string(decryptedData))
//	}
//
//	func decryptData(encryptedData []byte, armoredPrivateKey, passphrase string) ([]byte, error) {
//		// Unlock the private key
//		privateKey, err := crypto.NewKeyFromArmored(armoredPrivateKey)
//		if err != nil {
//			return nil, fmt.Errorf("failed to parse private key: %w", err)
//		}
//
//		unlockedKey, err := privateKey.Unlock([]byte(passphrase))
//		if err != nil {
//			return nil, fmt.Errorf("failed to unlock private key: %w", err)
//		}
//
//		// Decrypt the message
//		message := crypto.NewPGPMessage(encryptedData)
//		plainMessage, err := helper.DecryptMessage(unlockedKey, message)
//		if err != nil {
//			return nil, fmt.Errorf("failed to decrypt message: %w", err)
//		}
//
//		return plainMessage.GetBinary(), nil
//	}
//
// For enhanced traffic, consider using TLS over protocols like HTTPS or gRPC.
// Additionally, any network or other mechanism (e.g., Implement own storage mechanism, such as using buckets better than any (e.g., aws),
// which is ideal for Kubernetes environments.) that supports I/O operations will work well with this.
//
// Note: Memory allocations may vary depending on the input and output types.
// If writing to a file (file disk not a memory again), the allocations are minimal.
func (e *Encryptor) EncryptStream(i io.Reader, o io.Writer) error {
	// Create a key ring from the public key
	keyRing, err := e.createKeyRing()
	if err != nil {
		return err
	}

	// Create a pipe to handle streaming encryption
	r, w := io.Pipe()

	// Determine if the input is a file and set the filename.
	//
	// This effectively Go detects actual files in I/O.
	//
	// # Result:
	// 	Decrypt Operation - Success
	//
	// # General State:
	// 	- File Name: test_output_1960559248.txt
	// 	- MIME: false
	// 	- Message Integrity Protection: true
	// 	- Symmetric Encryption Algorithm: AES256.CFB
	// 	- German Encryption Standards: false
	var filename string
	if file, ok := i.(*os.File); ok {
		filename = filepath.Base(file.Name())
	} else {
		// If input is not a file, check if output is a file
		if outFile, ok := o.(*os.File); ok {
			outName := filepath.Base(outFile.Name())
			if strings.HasSuffix(outName, newGPGModern) {
				filename = strings.TrimSuffix(outName, newGPGModern)
			}
		}
	}

	// Create metadata (header) for the encryption
	metadata := &crypto.PlainMessageMetadata{
		IsBinary: true,
		Filename: filename,
		ModTime:  crypto.GetUnixTime(),
	}

	// Start a goroutine to handle encryption
	go func() {
		defer w.Close()
		// Create a writer for the encrypted output
		//
		// Note: When encrypting data then send over the network, whether secure network or insecure network,
		// additional compression is optional. This encryption process already includes built-in compression,
		// which can help reduce bandwidth costs.
		encryptWriter, err := keyRing.EncryptStreamWithCompression(w, metadata, nil)
		if err != nil {
			w.CloseWithError(fmt.Errorf("failed to create encryption stream: %w", err))
			return
		}
		defer encryptWriter.Close()

		// Note: The buffer size of 4096 bytes is suitable for streaming encryption.
		// It allows processing of large files or whole disk efficiently without loading the entire file into memory.
		buffer := make([]byte, 4096) // Define a buffer size
		for {
			n, err := i.Read(buffer)
			if err != nil && err != io.EOF {
				w.CloseWithError(fmt.Errorf("failed to read input: %w", err))
				return
			}
			if n == 0 {
				break
			}

			if _, err := encryptWriter.Write(buffer[:n]); err != nil {
				w.CloseWithError(fmt.Errorf("failed to write encrypted data: %w", err))
				return
			}
		}
	}()

	// Copy the encrypted data from the pipe reader to the output
	if _, err := io.Copy(o, r); err != nil {
		return fmt.Errorf("failed to write encrypted data to output: %w", err)
	}

	return nil
}
