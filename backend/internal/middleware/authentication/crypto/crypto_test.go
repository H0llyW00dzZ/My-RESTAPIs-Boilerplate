// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto_test

import (
	"crypto/rand"
	"encoding/base64"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto"
	"testing"
)

func TestVerifyCiphertext(t *testing.T) {
	// Create an instance of the cryptoService
	secryptKey := "gopher-testing-testing-testinggg"
	signKey := "gopher-testing-testing-testing"
	// Test case 1: Valid ciphertext and signature
	plaintext := "Hello, World!"
	encryptedData, signature, err := crypto.EncryptData(plaintext, false, secryptKey, signKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	if !crypto.VerifyCiphertext(encryptedData, signature, signKey) {
		t.Error("Expected valid ciphertext, but got invalid")
	}

	// Test case 2: Invalid signature
	invalidSignature := base64.StdEncoding.EncodeToString([]byte("invalid_signature"))
	if crypto.VerifyCiphertext(encryptedData, invalidSignature, signKey) {
		t.Error("Expected invalid ciphertext due to invalid signature, but got valid")
	}

	// Test case 3: Corrupted ciphertext
	corruptedData := encryptedData[:len(encryptedData)-5] + "12345"
	if crypto.VerifyCiphertext(corruptedData, signature, signKey) {
		t.Error("Expected invalid ciphertext due to corruption, but got valid")
	}

	// Test case 4: Ciphertext too short
	shortData := encryptedData[:10]
	if crypto.VerifyCiphertext(shortData, signature, signKey) {
		t.Error("Expected invalid ciphertext due to short length, but got valid")
	}

	// Test case 5: Invalid nonce length
	invalidNonceData := encryptedData[:5] + encryptedData[10:]
	if crypto.VerifyCiphertext(invalidNonceData, signature, signKey) {
		t.Error("Expected invalid ciphertext due to invalid nonce length, but got valid")
	}

	// Test case 6: Invalid ciphertext length
	invalidCiphertextData := encryptedData[:len(encryptedData)-3]
	if crypto.VerifyCiphertext(invalidCiphertextData, signature, signKey) {
		t.Error("Expected invalid ciphertext due to invalid ciphertext length, but got valid")
	}

	// Test case 7: Random data
	randomData := make([]byte, 100)
	_, err = rand.Read(randomData)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}
	randomEncryptedData := base64.StdEncoding.EncodeToString(randomData)
	if crypto.VerifyCiphertext(randomEncryptedData, signature, signKey) {
		t.Error("Expected invalid ciphertext for random data, but got valid")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Create an instance of the cryptoService
	secryptKey := "gopher-testing-testing-testinggg"
	signKey := "gopher-testing-testing-testing"
	// Test case 1: Encrypt and decrypt data using Argon2
	plaintext := "Hello, World!"
	encryptedData, signature, err := crypto.EncryptData(plaintext, true, secryptKey, signKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	decryptedText, err := crypto.DecryptData(encryptedData, signature, true, secryptKey, signKey)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}
	if decryptedText != plaintext {
		t.Errorf("Decrypted text does not match the original plaintext")
	}

	// Test case 2: Encrypt and decrypt data without Argon2
	encryptedData, signature, err = crypto.EncryptData(plaintext, false, secryptKey, signKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	decryptedText, err = crypto.DecryptData(encryptedData, signature, false, secryptKey, signKey)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}
	if decryptedText != plaintext {
		t.Errorf("Decrypted text does not match the original plaintext")
	}

	// Test case 3: Decrypt with invalid signature
	invalidSignature := base64.StdEncoding.EncodeToString([]byte("invalid_signature"))
	_, err = crypto.DecryptData(encryptedData, invalidSignature, false, secryptKey, signKey)
	if err != crypto.ErrorInvalidSignature {
		t.Errorf("Expected ErrorInvalidSignature, but got: %v", err)
	}

	// Test case 4: Decrypt with corrupted ciphertext
	corruptedData := encryptedData[:len(encryptedData)-5] + "12345"
	_, err = crypto.DecryptData(corruptedData, signature, false, secryptKey, signKey)
	if err == nil {
		t.Error("Expected decryption error, but got nil")
	}
}

func TestCryptoService(t *testing.T) {
	// Create an instance of the cryptoService
	useArgon2 := false
	secryptKey := "gopher-testing-testing-testinggg"
	signKey := "gopher-testing-testing-testing"
	service := crypto.New(useArgon2, secryptKey, signKey)

	// Test case 1: Encrypt and decrypt data
	plaintext := "Hello, World!"
	encryptedData, signature, err := service.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}
	decryptedText, err := service.Decrypt(encryptedData, signature)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}
	if decryptedText != plaintext {
		t.Errorf("Decrypted text does not match the original plaintext")
	}

	// Test case 2: Verify valid ciphertext
	if !service.VerifyCiphertext(encryptedData, signature) {
		t.Error("Expected valid ciphertext, but got invalid")
	}

	// Test case 3: Verify invalid signature
	invalidSignature := "invalid-signature"
	if service.VerifyCiphertext(encryptedData, invalidSignature) {
		t.Error("Expected invalid ciphertext due to invalid signature, but got valid")
	}

	// Test case 4: Verify corrupted ciphertext
	corruptedData := encryptedData[:len(encryptedData)-5] + "12345"
	if service.VerifyCiphertext(corruptedData, signature) {
		t.Error("Expected invalid ciphertext due to corruption, but got valid")
	}
}
