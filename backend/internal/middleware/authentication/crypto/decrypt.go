// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"golang.org/x/crypto/chacha20poly1305"
)

// decrypt decrypts the given ciphertext using a cascade of ciphers.
// Note: This method is a reliable and secure approach for decryption.
// An alternative method is to combine the techniques described in RFC 5652
// (https://www.rfc-editor.org/rfc/rfc5652.html) and RFC 5652 Section 6.3
// (https://datatracker.ietf.org/doc/html/rfc5652#section-6.3).
// However, that approach might have a higher risk of corrupting the ciphertext
// if not implemented carefully.
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	// Second decryption: ChaCha20-Poly1305
	chachaCipher := func(data []byte) ([]byte, error) {
		aead, err := chacha20poly1305.New(key)
		if err != nil {
			return nil, err
		}

		nonceSize := aead.NonceSize()
		if len(data) < nonceSize {
			return nil, ErrorInvalidCipherText
		}

		nonce := data[:nonceSize]
		ciphertext := data[nonceSize:]

		plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, err
		}

		return plaintext, nil
	}

	plaintext, err := decryptWithCipher(ciphertext, chachaCipher)
	if err != nil {
		return nil, err
	}

	// First decryption: AES
	aesCipher := func(data []byte) ([]byte, error) {
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}

		nonceSize := gcm.NonceSize()
		if len(data) < nonceSize {
			return nil, ErrorInvalidCipherText
		}

		nonce := data[:nonceSize]
		ciphertext := data[nonceSize:]

		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, err
		}

		return plaintext, nil
	}

	plaintext, err = decryptWithCipher(plaintext, aesCipher)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// DecryptData decrypts the given encrypted data using AES decryption with the same derived encryption key used during encryption.
// It expects the encrypted data and signature to be base64-encoded.
// It verifies the signature before decrypting the data.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the decryption key.
func DecryptData(encryptedData, signature string, useArgon2 bool, secryptKey, signKey string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return "", err
	}

	if !verifySignature(encryptedBytes, signatureBytes, signKey) {
		return "", ErrorInvalidSignature
	}

	salt := encryptedBytes[:16]
	ciphertext := encryptedBytes[16:]

	key := deriveKey(salt, useArgon2, secryptKey)
	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
