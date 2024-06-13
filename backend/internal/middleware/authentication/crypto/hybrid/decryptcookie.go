// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package hybrid

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/chacha20poly1305"
)

// DecryptCookie decrypts a cookie value using a hybrid decryption scheme.
func DecryptCookie(encodedCookie, key string, encoding string) (string, error) {
	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", ErrorInvalidKey
	}

	var decodedCookie []byte
	switch encoding {
	case "hex":
		decodedCookie, err = hex.DecodeString(encodedCookie)
	default:
		decodedCookie, err = base64.RawURLEncoding.DecodeString(encodedCookie)
	}
	if err != nil {
		return "", ErrorInvalidCookie
	}

	decryptFn := func(encryptedCookie []byte) ([]byte, error) {
		// Extract the nonces and encrypted cookie value
		if len(encryptedCookie) < 12+chacha20poly1305.NonceSizeX {
			return nil, ErrorInvalidCookie
		}
		aesNonce := encryptedCookie[:12]
		chachaNonce := encryptedCookie[12 : 12+chacha20poly1305.NonceSizeX]
		ciphertext := encryptedCookie[12+chacha20poly1305.NonceSizeX:]

		// Decrypt the cookie value using ChaCha20-Poly1305
		aesCiphertext, err := decryptChaCha20Poly1305(ciphertext, chachaNonce, keyDecoded)
		if err != nil {
			return nil, err
		}

		// Decrypt the AES-GCM ciphertext
		plaintext, err := decryptAESGCM(aesCiphertext, aesNonce, keyDecoded)
		if err != nil {
			return nil, err
		}

		return plaintext, nil
	}

	plaintext, err := decryptFn(decodedCookie)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// decryptAESGCM decrypts the ciphertext using AES-GCM and returns the plaintext.
func decryptAESGCM(ciphertext, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// decryptChaCha20Poly1305 decrypts the ciphertext using XChaCha20-Poly1305 and returns the plaintext.
func decryptChaCha20Poly1305(ciphertext, nonce, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
