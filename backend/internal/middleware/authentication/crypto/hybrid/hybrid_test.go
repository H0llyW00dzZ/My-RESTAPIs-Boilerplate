// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package hybrid_test

import (
	"bytes"
	"crypto/rand"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
	"testing"

	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

func TestEncryptDecryptCookie(t *testing.T) {
	// Generate a random encryption key
	key := encryptcookie.GenerateKey()

	// Test cases
	testCases := []struct {
		name     string
		value    string
		encoding string
	}{
		{
			name:     "Simple cookie value",
			value:    "hello world",
			encoding: "base64",
		},
		{
			name:     "Complex cookie value",
			value:    "!@#$%^&*()_+=-`~[]{}|;':\"<>,.?/\\",
			encoding: "base64",
		},
		{
			name:     "Empty cookie value",
			value:    "",
			encoding: "base64",
		},
		{
			name:     "Simple cookie value with hex encoding",
			value:    "hello world",
			encoding: "hex",
		},
		{
			name:     "Complex cookie value with hex encoding",
			value:    "!@#$%^&*()_+=-`~[]{}|;':\"<>,.?/\\",
			encoding: "hex",
		},
		{
			name:     "Empty cookie value with hex encoding",
			value:    "",
			encoding: "hex",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an instance of the hybrid encryption service
			service := hybrid.New(key, tc.encoding)

			// Encrypt the cookie value
			encryptedCookie, err := service.EncryptCookie(tc.value)
			if err != nil {
				t.Fatalf("Failed to encrypt cookie: %v", err)
			}

			// Decrypt the encrypted cookie value
			decryptedValue, err := service.DecryptCookie(encryptedCookie)
			if err != nil {
				t.Fatalf("Failed to decrypt cookie: %v", err)
			}

			// Compare the decrypted value with the original value
			if decryptedValue != tc.value {
				t.Errorf("Decrypted value does not match the original value")
				t.Errorf("Expected: %s", tc.value)
				t.Errorf("Got: %s", decryptedValue)
			}
		})
	}

	t.Run("Invalid cookie format", func(t *testing.T) {
		// Create an instance of the hybrid encryption service
		service := hybrid.New(key, "base64")

		// Create an invalid cookie format
		invalidCookie := "invalid-cookie-format"

		// Try to decrypt the invalid cookie
		_, err := service.DecryptCookie(invalidCookie)
		if err != hybrid.ErrorInvalidCookie {
			t.Errorf("Expected error: %v", hybrid.ErrorInvalidCookie)
			t.Errorf("Got: %v", err)
		}
	})

	t.Run("Invalid base64-encoded key", func(t *testing.T) {
		// Create an invalid base64-encoded key
		invalidKey := "invalid-key"

		// Create an instance of the hybrid encryption service with the invalid key
		service := hybrid.New(invalidKey, "base64")

		// Try to encrypt a cookie value with the invalid key
		_, err := service.EncryptCookie("test")
		if err != hybrid.ErrorInvalidKey {
			t.Errorf("Expected error: %v", hybrid.ErrorInvalidKey)
			t.Errorf("Got: %v", err)
		}

		// Try to decrypt a cookie value with the invalid key
		_, err = service.DecryptCookie("test")
		if err != hybrid.ErrorInvalidKey {
			t.Errorf("Expected error: %v", hybrid.ErrorInvalidKey)
			t.Errorf("Got: %v", err)
		}
	})
}

func TestStreamEncryptDecrypt(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // ChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate ChaCha20-Poly1305 key: %v", err)
	}

	// Test cases
	testCases := []struct {
		name  string
		value string
	}{
		{
			name:  "Simple data value",
			value: "hello world",
		},
		{
			name:  "Complex data value",
			value: "!@#$%^&*()_+=-`~[]{}|;':\"<>,.?/\\",
		},
		{
			name:  "Empty data value",
			value: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an instance of the stream encryption service
			service := hybrid.NewStreamService(aesKey, chachaKey)

			// Encrypt the data value
			encryptedData, err := service.Encrypt(tc.value)
			if err != nil {
				t.Fatalf("Failed to encrypt data: %v", err)
			}

			// Decrypt the encrypted data value
			decryptedValue, err := service.Decrypt(encryptedData)
			if err != nil {
				t.Fatalf("Failed to decrypt data: %v", err)
			}

			// Compare the decrypted value with the original value
			if decryptedValue != tc.value {
				t.Errorf("Decrypted value does not match the original value")
				t.Errorf("Expected: %s", tc.value)
				t.Errorf("Got: %s", decryptedValue)
			}
		})
	}
}

func TestHybridEncryptDecryptStream(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // ChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate ChaCha20-Poly1305 key: %v", err)
	}

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = stream.EncryptStream(inputBuffer, encryptedBuffer, aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = stream.DecryptStream(encryptedBuffer, decryptedBuffer, aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext. Got: %s, Want: %s", decryptedData, plaintext)
	}
}

func TestHybridEncryptDecryptStreamWithApiKey(t *testing.T) {
	// Predefined API keys or secret keys, which should be securely stored and retrieved.
	aesKey := []byte("gopher-testing-testing-testinggg")
	chachaKey := []byte("gopher-testing-testing-testinggg")

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err := stream.EncryptStream(inputBuffer, encryptedBuffer, aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = stream.DecryptStream(encryptedBuffer, decryptedBuffer, aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext. Got: %s, Want: %s", decryptedData, plaintext)
	}
}

func TestHybridEncryptDecryptStreamLargeData(t *testing.T) {
	// Note: Works well testing on AMD Ryzen 9 3900x 12-Core Processor (24 CPUs) RAM 32GB
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // ChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate ChaCha20-Poly1305 key: %v", err)
	}

	// Generate a large plaintext data.
	plaintextSize := 10 * 1024 * 1024 // 10 MB
	plaintext := make([]byte, plaintextSize)
	_, err = rand.Read(plaintext)
	if err != nil {
		t.Fatalf("Failed to generate plaintext: %v", err)
	}

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = stream.EncryptStream(inputBuffer, encryptedBuffer, aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = stream.DecryptStream(encryptedBuffer, decryptedBuffer, aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext.")
	}
}
