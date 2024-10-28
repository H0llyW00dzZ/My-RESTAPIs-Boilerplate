// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package gpg provides functionality for encrypting data using OpenPGP/GPG public keys.
//
// This package includes utilities to create and manage key rings, encrypt files,
// and handle streaming encryption for efficient data transmission or other.
//
// Unlike GPG Proton built on top (fork) the standard library, this package
// is designed with a modern approach, focusing on top of I/O operations.
//
// # Example TCP (transmission over a network):
//
//	// Server-side example
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
//
//	// Client-side example
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
// # Example WebSocket (transmission over a network):
//
//	// Server-side example
//	func handleWebSocket(conn *websocket.Conn, encryptor *gpg.Encryptor) error {
//	    // Prepare the data to be encrypted
//	    inputData := bytes.NewReader([]byte("This is the data to be encrypted."))
//
//	    // Encrypt and send data over the WebSocket connection
//	    writer, err := conn.NextWriter(websocket.BinaryMessage)
//	    if err != nil {
//	        return fmt.Errorf("failed to get writer: %w", err)
//	    }
//	    defer writer.Close()
//
//	    if err := encryptor.EncryptStream(inputData, writer); err != nil {
//	        return fmt.Errorf("failed to encrypt and send data: %w", err)
//	    }
//
//	    return nil
//	}
//
//	// Client-side example
//	func readWebSocket(conn *websocket.Conn) ([]byte, error) {
//	    _, reader, err := conn.NextReader()
//	    if err != nil {
//	        return nil, fmt.Errorf("failed to get reader: %w", err)
//	    }
//
//	    var encryptedData bytes.Buffer
//	    if _, err := io.Copy(&encryptedData, reader); err != nil {
//	        return nil, fmt.Errorf("failed to read encrypted data: %w", err)
//	    }
//
//	    // Decrypt the data as needed
//	    // decryptedData, err := decryptData(encryptedData.Bytes(), "your-armored-private-key", "your-passphrase")
//	    // if err != nil {
//	    //     return nil, fmt.Errorf("failed to decrypt data: %w", err)
//	    // }
//
//	    // return decryptedData, nil
//	    return encryptedData.Bytes(), nil
//	}
//
// # Example UDP (transmission over a network):
//
//	// Server-side example
//	func main() {
//		addr := net.UDPAddr{
//			Port: 8080,
//			IP:   net.ParseIP("0.0.0.0"),
//		}
//		conn, err := net.ListenUDP("udp", &addr)
//		if err != nil {
//			log.Fatalf("Failed to listen on UDP port 8080: %v", err)
//		}
//		defer conn.Close()
//
//		log.Println("Server listening on UDP port 8080")
//
//		buffer := make([]byte, 4096)
//		for {
//			n, clientAddr, err := conn.ReadFromUDP(buffer)
//			if err != nil {
//				log.Printf("Failed to read from UDP: %v", err)
//				continue
//			}
//
//			go handleUDPConnection(conn, clientAddr, buffer[:n])
//		}
//	}
//
//	func handleUDPConnection(conn *net.UDPConn, clientAddr *net.UDPAddr, data []byte) {
//		// Create a new Encryptor with the public key
//		encryptor, err := gpg.NewEncryptor([]string{"your-armored-public-key"})
//		if err != nil {
//			log.Printf("Failed to create encryptor: %v", err)
//			return
//		}
//
//		inputBuffer := bytes.NewReader(data)
//		var outputBuffer bytes.Buffer
//
//		// Encrypt the data
//		if err := encryptor.EncryptStream(inputBuffer, &outputBuffer); err != nil {
//			log.Printf("Failed to encrypt data: %v", err)
//			return
//		}
//
//		// Send the encrypted data back to the client
//		_, err = conn.WriteToUDP(outputBuffer.Bytes(), clientAddr)
//		if err != nil {
//			log.Printf("Failed to send data: %v", err)
//		}
//	}
//
//	// Client-side example
//	func main() {
//		serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
//		if err != nil {
//			log.Fatalf("Failed to resolve server address: %v", err)
//		}
//
//		conn, err := net.DialUDP("udp", nil, serverAddr)
//		if err != nil {
//			log.Fatalf("Failed to dial UDP server: %v", err)
//		}
//		defer conn.Close()
//
//		// Prepare data to be sent
//		data := []byte("This is the data to be encrypted.")
//		_, err = conn.Write(data)
//		if err != nil {
//			log.Fatalf("Failed to send data: %v", err)
//		}
//
//		// Read the encrypted response
//		buffer := make([]byte, 4096)
//		n, _, err := conn.ReadFromUDP(buffer)
//		if err != nil {
//			log.Fatalf("Failed to read response: %v", err)
//		}
//
//		encryptedData := buffer[:n]
//		fmt.Println("Received encrypted data:", encryptedData)
//
//		// Decrypt the data
//		decryptedData, err := decryptData(encryptedData, "your-armored-private-key", "your-passphrase")
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
// Additionally, any network or other mechanism (e.g., Implement own storage mechanism over the network, such as using buckets better than any (e.g., aws),
// which is ideal for Kubernetes environments.) that supports I/O operations will work well with this.
//
// # Compatibility Considerations:
//
// The gopenpgp/gpgproton library and this package rely on specific cipher algorithms and
// Go language features. It's important to consider compatibility when updating
// Go versions or changing cipher settings.
//
// # Cipher Compatibility:
//
//   - The default cipher mode for OpenPGP is typically CFB (Cipher Feedback). While
//     it is secure for most applications, ensure that any changes to the cipher
//     algorithm are supported by all systems involved in the encryption/decryption process.
//   - If you modify the cipher settings, verify that these changes are compatible
//     with the gopenpgp/gpgproton library's capabilities and the OpenPGP standard
//
// # Go Version Compatibility:
//
//   - Always test your application when upgrading to a new Go version. Changes in
//     the Go runtime or standard library might affect cryptographic operations.
//   - If a new Go version introduces changes that impact cipher operations or
//     cryptographic libraries, evaluate whether these changes are necessary for
//     your application. If the changes are not relevant/not make any sense (e.g., deprecating certain ciphers that are still secure),
//     consider postponing the upgrade until compatibility is ensured. Alternatively, just fork the old version and maintain it yourself,
//     as Go programming makes this manageable.
package gpg
