// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto_test

import (
	"crypto/rand"
	"encoding/base64"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto"
	"testing"

	// Import the godotenv package for loading environment variables from a .env file
	// The "_" blank identifier is used to import the package for its side effects (auto-loading .env file)
	_ "github.com/joho/godotenv/autoload"
)

func TestVerifyCiphertext(t *testing.T) {
	// Test case 1: Valid ciphertext and signature
	plaintext := "Hello, World!"
	encryptedData, signature, err := crypto.EncryptData(plaintext, false)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	if !crypto.VerifyCiphertext(encryptedData, signature) {
		t.Error("Expected valid ciphertext, but got invalid")
	}

	// Test case 2: Invalid signature
	invalidSignature := base64.StdEncoding.EncodeToString([]byte("invalid_signature"))
	if crypto.VerifyCiphertext(encryptedData, invalidSignature) {
		t.Error("Expected invalid ciphertext due to invalid signature, but got valid")
	}

	// Test case 3: Corrupted ciphertext
	corruptedData := encryptedData[:len(encryptedData)-5] + "12345"
	if crypto.VerifyCiphertext(corruptedData, signature) {
		t.Error("Expected invalid ciphertext due to corruption, but got valid")
	}

	// Test case 4: Ciphertext too short
	shortData := encryptedData[:10]
	if crypto.VerifyCiphertext(shortData, signature) {
		t.Error("Expected invalid ciphertext due to short length, but got valid")
	}

	// Test case 5: Invalid nonce length
	invalidNonceData := encryptedData[:5] + encryptedData[10:]
	if crypto.VerifyCiphertext(invalidNonceData, signature) {
		t.Error("Expected invalid ciphertext due to invalid nonce length, but got valid")
	}

	// Test case 6: Invalid ciphertext length
	invalidCiphertextData := encryptedData[:len(encryptedData)-3]
	if crypto.VerifyCiphertext(invalidCiphertextData, signature) {
		t.Error("Expected invalid ciphertext due to invalid ciphertext length, but got valid")
	}

	// Test case 7: Random data
	randomData := make([]byte, 100)
	_, err = rand.Read(randomData)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}
	randomEncryptedData := base64.StdEncoding.EncodeToString(randomData)
	if crypto.VerifyCiphertext(randomEncryptedData, signature) {
		t.Error("Expected invalid ciphertext for random data, but got valid")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Test case 1: Encrypt and decrypt data using Argon2
	plaintext := "Hello, World!"
	encryptedData, signature, err := crypto.EncryptData(plaintext, true)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	decryptedText, err := crypto.DecryptData(encryptedData, signature, true)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}
	if decryptedText != plaintext {
		t.Errorf("Decrypted text does not match the original plaintext")
	}

	// Test case 2: Encrypt and decrypt data without Argon2
	encryptedData, signature, err = crypto.EncryptData(plaintext, false)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	decryptedText, err = crypto.DecryptData(encryptedData, signature, false)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}
	if decryptedText != plaintext {
		t.Errorf("Decrypted text does not match the original plaintext")
	}

	// Test case 3: Decrypt with invalid signature
	invalidSignature := base64.StdEncoding.EncodeToString([]byte("invalid_signature"))
	_, err = crypto.DecryptData(encryptedData, invalidSignature, false)
	if err != crypto.ErrorInvalidSignature {
		t.Errorf("Expected ErrorInvalidSignature, but got: %v", err)
	}

	// Test case 4: Decrypt with corrupted ciphertext
	corruptedData := encryptedData[:len(encryptedData)-5] + "12345"
	_, err = crypto.DecryptData(corruptedData, signature, false)
	if err == nil {
		t.Error("Expected decryption error, but got nil")
	}
}
