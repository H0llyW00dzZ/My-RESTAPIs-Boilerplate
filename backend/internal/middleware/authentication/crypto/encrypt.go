// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"

	// Import the godotenv package for loading environment variables from a .env file
	// The "_" blank identifier is used to import the package for its side effects (auto-loading .env file)
	_ "github.com/joho/godotenv/autoload"
)

var (
	secryptkey = os.Getenv("SECRETCRYPT_KEY")
)

var (
	// ErrorInvalidCipherText is a custom error variable that represents an error
	// which occurs when the ciphertext (encrypted text) is invalid or malformed.
	ErrorInvalidCipherText = errors.New("invalid ciphertext")
)

// EncryptData encrypts the given token using AES encryption with the provided encryption key.
// It returns the base64-encoded ciphertext, which consists of the nonce concatenated with the encrypted data.
func EncryptData(token string) (string, error) {
	block, err := aes.NewCipher([]byte(secryptkey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptData decrypts the given encrypted token using AES decryption with the same encryption key used during encryption.
// It expects the encrypted token to be base64-encoded and returns the decrypted plaintext data.
func DecryptData(encryptedToken string) (string, error) {
	// Decode the base64-encoded ciphertext to obtain the original ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block using the same encryption key used during encryption
	block, err := aes.NewCipher([]byte(secryptkey))
	if err != nil {
		return "", err
	}

	// Create a new GCM mode instance for decryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce from the ciphertext
	// The nonce is prepended to the ciphertext during encryption and is required for decryption
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrorInvalidCipherText
	}
	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	// Decrypt the ciphertext using the nonce and the same encryption key
	// The Open function returns the decrypted plaintext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	// Convert the decrypted plaintext bytes to a string and return it
	return string(plaintext), nil
}
