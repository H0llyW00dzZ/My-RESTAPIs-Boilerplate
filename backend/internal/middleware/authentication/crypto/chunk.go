// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/cipher"
	"crypto/rand"
)

const (
	aesNonceSize = 16 // AES-CTR nonce size (16 bytes for a 128-bit nonce)
	chunkSize    = 1024
)

// encryptChunk encrypts a single chunk using AES-CTR and ChaCha20-Poly1305.
func encryptChunk(aesBlock cipher.Block, chacha cipher.AEAD, chunk []byte) ([]byte, []byte, error) {
	// Generate a nonce for AES-CTR.
	aesNonce := make([]byte, aesNonceSize)
	if _, err := rand.Read(aesNonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the chunk using AES-CTR.
	aesStream := cipher.NewCTR(aesBlock, aesNonce)
	aesEncryptedChunk := make([]byte, len(chunk))
	aesStream.XORKeyStream(aesEncryptedChunk, chunk)

	// Prepend the AES nonce to the AES-CTR encrypted chunk.
	aesEncryptedChunkWithNonce := append(aesNonce, aesEncryptedChunk...)

	// Generate a nonce for ChaCha20-Poly1305.
	chachaNonce := make([]byte, chacha.NonceSize())
	if _, err := rand.Read(chachaNonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the AES-CTR encrypted chunk (including the AES nonce) using ChaCha20-Poly1305.
	chachaEncryptedChunk := chacha.Seal(nil, chachaNonce, aesEncryptedChunkWithNonce, nil)

	return chachaNonce, chachaEncryptedChunk, nil
}

// decryptChunk decrypts a single chunk using ChaCha20-Poly1305 and AES-CTR.
func decryptChunk(aesBlock cipher.Block, chacha cipher.AEAD, chachaNonce, chachaEncryptedChunk []byte) ([]byte, error) {
	// Decrypt the chunk using ChaCha20-Poly1305.
	aesEncryptedChunk, err := chacha.Open(nil, chachaNonce, chachaEncryptedChunk, nil)
	if err != nil {
		return nil, err
	}

	// Extract the AES nonce from the beginning of the AES-CTR encrypted chunk.
	aesNonce := aesEncryptedChunk[:aesNonceSize]
	aesEncryptedChunk = aesEncryptedChunk[aesNonceSize:]

	// Decrypt the chunk using AES-CTR.
	aesStream := cipher.NewCTR(aesBlock, aesNonce)
	chunk := make([]byte, len(aesEncryptedChunk))
	aesStream.XORKeyStream(chunk, aesEncryptedChunk)

	return chunk, nil
}
