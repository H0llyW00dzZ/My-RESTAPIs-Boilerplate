// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"io"

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

// processLargeData is a higher-order function that processes large data using the provided processor function.
// It reads the data from the provided io.Reader and writes the processed data to the provided io.Writer.
// The processor function is responsible for encrypting or decrypting the data.
// It generates a signature for the processed data and appends it to the output.
func processLargeData(src io.Reader, dst io.Writer, useArgon2 bool, secryptKey, signKey string, processor func([]byte, []byte) ([]byte, error)) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(salt, useArgon2, secryptKey)

	if _, err := dst.Write(salt); err != nil {
		return err
	}

	hash := sha256.New()
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		processed, err := processor(buf[:n], key)
		if err != nil {
			return err
		}

		if _, err := dst.Write(processed); err != nil {
			return err
		}

		hash.Write(processed)
	}

	signature := signData(hash.Sum(nil), signKey)
	if _, err := dst.Write(signature); err != nil {
		return err
	}

	return nil
}
