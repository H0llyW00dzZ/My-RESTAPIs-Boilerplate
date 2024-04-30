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
	"os"

	// Import the godotenv package for loading environment variables from a .env file
	// The "_" blank identifier is used to import the package for its side effects (auto-loading .env file)
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/argon2"
)

var (
	// secryptkey holds the secret encryption key.
	//
	// NOTE: In production, this key should be kept secret and not stored in an environment variable.
	// The reason it is set from an environment variable here is for ease of testing.
	// Retrieve the secret encryption key from the environment variable "SECRETCRYPT_KEY".
	secryptkey = os.Getenv("SECRETCRYPT_KEY")

	// signkey holds the secret signing key.
	//
	// NOTE: In production, this key should be kept secret and not stored in an environment variable.
	// The reason it is set from an environment variable here is for ease of testing.
	// Retrieve the secret signing key from the environment variable "SIGN_KEY".
	signkey = os.Getenv("SIGN_KEY")
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
func signData(data []byte) []byte {
	mac := hmac.New(sha256.New, []byte(signkey))
	mac.Write(data)
	return mac.Sum(nil)
}

// verifySignature verifies the HMAC signature of the given data using the signing key.
func verifySignature(data, signature []byte) bool {
	expectedMAC := signData(data)
	return subtle.ConstantTimeCompare(signature, expectedMAC) == 1
}

// deriveKey derives an encryption key using Argon2 key derivation function or returns the secryptkey directly.
func deriveKey(salt []byte, useArgon2 bool) []byte {
	// Note: Using Argon2 is expensive (100MB+ per encrypt/decrypt) the cost, which is not recommended.
	// I might try to introduce cryptographic techniques to implement a similar but cheaper approach and a new cipher from scratch later.
	if useArgon2 {
		return argon2.IDKey([]byte(secryptkey), salt, 1, 64*1024, 4, 32)
	}
	return []byte(secryptkey)
}

// processLargeData is a higher-order function that processes large data using the provided processor function.
// It reads the data from the provided io.Reader and writes the processed data to the provided io.Writer.
// The processor function is responsible for encrypting or decrypting the data.
// It generates a signature for the processed data and appends it to the output.
func processLargeData(src io.Reader, dst io.Writer, useArgon2 bool, processor func([]byte, []byte) ([]byte, error)) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(salt, useArgon2)

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

	signature := signData(hash.Sum(nil))
	if _, err := dst.Write(signature); err != nil {
		return err
	}

	return nil
}
