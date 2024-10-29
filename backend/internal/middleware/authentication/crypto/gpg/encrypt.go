// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"
	"io"
	"os"

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
		IsBinary: e.config.isBinary,
		Filename: inFile.Name(),
		ModTime:  e.config.modTime,
	}

	// Choose the appropriate encryption function based on the compress option
	encryptFunc := keyRing.EncryptStream
	if e.config.compress {
		encryptFunc = keyRing.EncryptStreamWithCompression
	}

	// Create a writer for the encrypted output
	encryptWriter, err := encryptFunc(outFile, metadata, nil)
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

	// Determine if the input and output I/O is a file and set the filename.
	filename, err := extractFilename(i, o, newGPGModern)
	if err != nil {
		return err
	}

	// Create metadata (header) for the encryption
	metadata := &crypto.PlainMessageMetadata{
		IsBinary: e.config.isBinary,
		Filename: filename,
		ModTime:  e.config.modTime,
	}

	// Choose the appropriate encryption function based on the compress option
	encryptFunc := keyRing.EncryptStream
	if e.config.compress {
		encryptFunc = keyRing.EncryptStreamWithCompression
	}

	// Start a goroutine to handle encryption
	go func() {
		defer w.Close()
		// Create a writer for the encrypted output
		//
		// Note: When encrypting data to send over the network, whether on a secure or insecure network,
		// additional compression is optional (e.g., disable compression). This encryption process already includes built-in compression,
		// which can help reduce bandwidth costs and also helps reduce memory usage.
		encryptWriter, err := encryptFunc(w, metadata, nil)
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

func (e *Encryptor) encryptArmored(armoredKey string) (string, error) {
	keyRing, err := e.createKeyRing()
	if err != nil {
		return "", fmt.Errorf("failed to create key ring: %w", err)
	}

	message := crypto.NewPlainMessageFromString(armoredKey)
	encryptedMessage, err := keyRing.Encrypt(message, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt armored key: %w", err)
	}

	// Use custom headers to remove the comment
	armored, err := encryptedMessage.GetArmoredWithCustomHeaders(customHeader, keyBoxVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get armored message: %w", err)
	}

	return armored, nil
}
