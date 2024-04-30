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
	"golang.org/x/crypto/argon2"
)

var (
	// secryptkey holds the secret encryption key retrieved from the environment variable "SECRETCRYPT_KEY"
	secryptkey = os.Getenv("SECRETCRYPT_KEY")
)

var (
	// ErrorInvalidCipherText is a custom error variable that represents an error
	// which occurs when the ciphertext (encrypted text) is invalid or malformed.
	ErrorInvalidCipherText = errors.New("invalid ciphertext")
)

// EncryptData encrypts the given data using AES encryption with a derived encryption key.
// It returns the base64-encoded ciphertext, which consists of the salt, nonce, and encrypted data.
func EncryptData(data string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Derive a secure encryption key using Argon2 key derivation function
	key := argon2.IDKey([]byte(secryptkey), salt, 1, 64*1024, 4, 32)

	// Create a new AES cipher block using the derived encryption key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM mode instance for encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a random nonce (number used once) for each encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data using the GCM mode and the generated nonce
	ciphertext := gcm.Seal(nil, nonce, []byte(data), nil)

	// Allocate a buffer to store the salt, nonce, and ciphertext
	encryptedData := make([]byte, 16+len(nonce)+len(ciphertext))
	copy(encryptedData[:16], salt)
	copy(encryptedData[16:16+len(nonce)], nonce)
	copy(encryptedData[16+len(nonce):], ciphertext)

	// Encode the encrypted data to base64 and return it
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// DecryptData decrypts the given encrypted data using AES decryption with the same derived encryption key used during encryption.
// It expects the encrypted data to be base64-encoded and contains the salt, nonce, and ciphertext.
func DecryptData(encryptedData string) (string, error) {
	// Decode the base64-encoded encrypted data
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	// Extract the salt from the encrypted data
	salt := encryptedBytes[:16]

	// Derive the encryption key using Argon2 key derivation function with the extracted salt
	key := argon2.IDKey([]byte(secryptkey), salt, 1, 64*1024, 4, 32)

	// Create a new AES cipher block using the derived encryption key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM mode instance for decryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce and ciphertext from the encrypted data
	nonceSize := gcm.NonceSize()
	if len(encryptedBytes) < nonceSize+16 {
		return "", ErrorInvalidCipherText
	}
	nonce := encryptedBytes[16 : 16+nonceSize]
	ciphertext := encryptedBytes[16+nonceSize:]

	// Decrypt the ciphertext using the nonce and the derived encryption key
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	// Convert the decrypted plaintext bytes to a string and return it
	return string(plaintext), nil
}
