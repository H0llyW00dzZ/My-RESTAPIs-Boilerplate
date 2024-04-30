// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	// Import the godotenv package for loading environment variables from a .env file
	// The "_" blank identifier is used to import the package for its side effects (auto-loading .env file)
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/chacha20poly1305"
)

// encrypt encrypts the given data using a cascade of ciphers.
func encrypt(data []byte, key []byte) ([]byte, error) {
	// First encryption: AES
	aesCipher := func(data []byte) ([]byte, error) {
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}

		ciphertext := gcm.Seal(nil, nonce, data, nil)
		ciphertext = append(nonce, ciphertext...)
		return ciphertext, nil
	}

	ciphertext, err := encryptWithCipher(data, aesCipher)
	if err != nil {
		return nil, err
	}

	// Second encryption: ChaCha20-Poly1305
	chachaCipher := func(data []byte) ([]byte, error) {
		aead, err := chacha20poly1305.New(key)
		if err != nil {
			return nil, err
		}

		nonce := make([]byte, aead.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}

		ciphertext := aead.Seal(nil, nonce, data, nil)
		ciphertext = append(nonce, ciphertext...)
		return ciphertext, nil
	}

	ciphertext, err = encryptWithCipher(ciphertext, chachaCipher)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// EncryptData encrypts the given data using AES encryption with a derived encryption key and signs the ciphertext.
// It returns the base64-encoded ciphertext and signature.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the encryption key.
func EncryptData(data string, useArgon2 bool) (string, string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	key := deriveKey(salt, useArgon2)
	ciphertext, err := encrypt([]byte(data), key)
	if err != nil {
		return "", "", err
	}

	encryptedData := append(salt, ciphertext...)
	signature := signData(encryptedData)

	encodedData := base64.StdEncoding.EncodeToString(encryptedData)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	return encodedData, encodedSignature, nil
}
