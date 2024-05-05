// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package hybrid_test

import (
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid"
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
