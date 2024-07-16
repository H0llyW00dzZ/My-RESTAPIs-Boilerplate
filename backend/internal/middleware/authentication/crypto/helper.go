// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrorInvalidCipherText is a custom error variable that represents an error
	// which occurs when the ciphertext (encrypted text) is invalid or malformed.
	ErrorInvalidCipherText = errors.New("invalid ciphertext")
	// ErrorInvalidSignature is a custom error variable that represents an error
	// which occurs when the signature is invalid or does not match the expected signature.
	ErrorInvalidSignature = errors.New("invalid signature")
)

// signData generates an HMAC signature for the given data using the signing key.
func signData(data []byte, signKey string) []byte {
	mac := hmac.New(sha256.New, []byte(signKey))
	mac.Write(data)
	return mac.Sum(nil)
}

// verifySignature verifies the HMAC signature of the given data using the signing key.
func verifySignature(data, signature []byte, signKey string) bool {
	expectedMAC := signData(data, signKey)
	return subtle.ConstantTimeCompare(signature, expectedMAC) == 1
}

// deriveKey derives an encryption key using Argon2 key derivation function or returns the secryptkey directly.
func deriveKey(salt []byte, useArgon2 bool, secryptKey string) []byte {
	// Note: Using Argon2 is expensive (100MB+ per encrypt/decrypt) the cost, which is not recommended.
	// I might try to introduce cryptographic techniques to implement a similar but cheaper approach and a new cipher from scratch later.
	if useArgon2 {
		return argon2.IDKey([]byte(secryptKey), salt, 1, 64*1024, 4, 32)
	}
	return []byte(secryptKey)
}
