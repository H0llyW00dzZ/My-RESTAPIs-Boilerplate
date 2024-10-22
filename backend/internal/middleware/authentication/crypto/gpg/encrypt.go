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

// EncryptFile encrypts the given file using the provided PGP public key.
func (e *Encryptor) EncryptFile(inputFile, outputFile string) error {
	// Read the public key
	key, err := crypto.NewKeyFromArmored(e.publicKey)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	keyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		return fmt.Errorf("failed to create key ring: %w", err)
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
	defer outFile.Close()

	// Create metadata for the encryption
	metadata := &crypto.PlainMessageMetadata{
		IsBinary: true,
	}

	// Create a writer for the encrypted output
	encryptWriter, err := keyRing.EncryptStream(outFile, metadata, nil)
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

// EncryptStream encrypts data from an input stream and writes to an output stream using the provided PGP public key.
//
// TODO: Improve this Use Struct
func EncryptStream(input io.Reader, output io.Writer, publicKey string) error {
	// Read the public key
	key, err := crypto.NewKeyFromArmored(publicKey)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	keyRing, err := crypto.NewKeyRing(key)
	if err != nil {
		return fmt.Errorf("failed to create key ring: %w", err)
	}

	// Note: The buffer size of 4096 bytes is suitable for streaming encryption.
	// It allows processing of large files or whole disk efficiently without loading the entire file into memory.
	buffer := make([]byte, 4096) // Define a buffer size
	for {
		n, err := input.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read input: %w", err)
		}
		if n == 0 {
			break
		}

		// Encrypt the data chunk
		message := crypto.NewPlainMessage(buffer[:n])
		encryptedMessage, err := keyRing.Encrypt(message, nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt data: %w", err)
		}

		// Write the encrypted data to the output
		if _, err := output.Write(encryptedMessage.GetBinary()); err != nil {
			return fmt.Errorf("failed to write encrypted data: %w", err)
		}
	}

	return nil
}
