// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package hybrid

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

// EncryptCookie encrypts a cookie value using a hybrid encryption scheme.
func EncryptCookie(value, key string) (string, error) {
	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", ErrorInvalidKey
	}

	encryptFn := func(plaintext []byte) ([]byte, error) {
		// Encrypt the plaintext using AES-GCM
		ciphertext, aesNonce, err := encryptAESGCM(plaintext, keyDecoded)
		if err != nil {
			return nil, err
		}

		// Encrypt the AES-GCM ciphertext using ChaCha20-Poly1305
		encryptedCookie, chachaNonce, err := encryptChaCha20Poly1305(ciphertext, keyDecoded)
		if err != nil {
			return nil, err
		}

		// Combine the nonces and encrypted cookie value
		// Note: this strong, required 99999999999999 cpu to brute force it.
		noncesAndCiphertext := append(aesNonce, chachaNonce...)
		noncesAndCiphertext = append(noncesAndCiphertext, encryptedCookie...)
		return noncesAndCiphertext, nil
	}

	encryptedCookie, err := encryptFn([]byte(value))
	if err != nil {
		return "", err
	}

	encodedCookie := base64.RawURLEncoding.EncodeToString(encryptedCookie)
	return encodedCookie, nil
}

// encryptAESGCM encrypts the plaintext using AES-GCM and returns the ciphertext and nonce.
func encryptAESGCM(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// encryptChaCha20Poly1305 encrypts the plaintext using ChaCha20-Poly1305 and returns the ciphertext and nonce.
func encryptChaCha20Poly1305(plaintext, key []byte) ([]byte, []byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}
