// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
)

func BenchmarkHybridEncryptDecryptStream(b *testing.B) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		b.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		b.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		b.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate plaintext data.
	plaintextSize := 1024 * 1024 // 1 MB
	plaintext := make([]byte, plaintextSize)
	_, err = rand.Read(plaintext)
	if err != nil {
		b.Fatalf("Failed to generate plaintext: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Encrypt the data.
		inputBuffer := bytes.NewBuffer(plaintext)
		encryptedBuffer := new(bytes.Buffer)
		err = s.Encrypt(inputBuffer, encryptedBuffer)
		if err != nil {
			b.Fatalf("Failed to encrypt data: %v", err)
		}

		// Decrypt the data.
		encryptedData := encryptedBuffer.Bytes()
		encryptedBuffer = bytes.NewBuffer(encryptedData)
		decryptedBuffer := new(bytes.Buffer)
		err = s.Decrypt(encryptedBuffer, decryptedBuffer)
		if err != nil {
			b.Fatalf("Failed to decrypt data: %v", err)
		}
	}
}

func BenchmarkHybridEncryptDecryptStreamWithHMAC(b *testing.B) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		b.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		b.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		b.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate a random HMAC key.
	hmacKey := make([]byte, 32)
	_, err = rand.Read(hmacKey)
	if err != nil {
		b.Fatalf("Failed to generate HMAC key: %v", err)
	}

	// Enable HMAC authentication.
	s.EnableHMAC(hmacKey)

	// Generate plaintext data.
	plaintextSize := 1024 * 1024 // 1 MB
	plaintext := make([]byte, plaintextSize)
	_, err = rand.Read(plaintext)
	if err != nil {
		b.Fatalf("Failed to generate plaintext: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Encrypt the data.
		inputBuffer := bytes.NewBuffer(plaintext)
		encryptedBuffer := new(bytes.Buffer)
		err = s.Encrypt(inputBuffer, encryptedBuffer)
		if err != nil {
			b.Fatalf("Failed to encrypt data: %v", err)
		}

		// Calculate the HMAC digest.
		//
		// Note: This Calculate HMAC digest does not actually calculate, its a trick
		encryptedData := encryptedBuffer.Bytes()
		_, err = s.Digest(bytes.NewReader(encryptedData))
		if err != nil {
			b.Fatalf("Failed to calculate HMAC digest: %v", err)
		}

		// Decrypt the data.
		encryptedBuffer = bytes.NewBuffer(encryptedData)
		decryptedBuffer := new(bytes.Buffer)
		err = s.Decrypt(encryptedBuffer, decryptedBuffer)
		if err != nil {
			b.Fatalf("Failed to decrypt data: %v", err)
		}
	}
}
