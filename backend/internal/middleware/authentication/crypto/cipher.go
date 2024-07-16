// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package crypto

import (
	"crypto/subtle"
	"encoding/base64"
)

// encryptWithCipher is a higher-order function that encrypts the given data using the provided cipher.
func encryptWithCipher(data []byte, cipher func([]byte) ([]byte, error)) ([]byte, error) {
	ciphertext, err := cipher(data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// decryptWithCipher is a higher-order function that decrypts the given ciphertext using the provided cipher.
func decryptWithCipher(ciphertext []byte, cipher func([]byte) ([]byte, error)) ([]byte, error) {
	plaintext, err := cipher(ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// VerifyCiphertext verifies the integrity of the ciphertext without decrypting it.
// It checks if the ciphertext has a valid structure and matches the expected format.
// It expects the encrypted data and signature to be base64-encoded.
// It returns true if the ciphertext is valid, false otherwise.
func VerifyCiphertext(encryptedData, signature, signKey string) bool {
	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return false
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(decodedData) < 16 {
		return false
	}

	expectedSignature := signData(decodedData, signKey)

	return subtle.ConstantTimeCompare(decodedSignature, expectedSignature) == 1
}
