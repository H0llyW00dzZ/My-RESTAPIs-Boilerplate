// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
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

// encryptAndWriteChunk encrypts a chunk and writes it to the output stream.
func encryptAndWriteChunk(aesBlock cipher.Block, chacha cipher.AEAD, chunk []byte, output io.Writer) error {
	chachaNonce, encryptedChunk, err := encryptChunk(aesBlock, chacha, chunk)
	if err != nil {
		return err
	}

	if err := writeChunk(encryptedChunk, chachaNonce, output); err != nil {
		return err
	}

	return nil
}

// readAndDecryptChunk reads an encrypted chunk from the input stream and decrypts it.
func readAndDecryptChunk(aesBlock cipher.Block, chacha cipher.AEAD, input io.Reader) ([]byte, error) {
	chunkSize, chachaNonce, err := readChunkMetadata(input)
	if err != nil {
		return nil, err
	}

	encryptedChunk := make([]byte, chunkSize)
	if _, err := io.ReadFull(input, encryptedChunk); err != nil {
		return nil, err
	}

	chunk, err := decryptChunk(aesBlock, chacha, chachaNonce, encryptedChunk)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// writeChunk writes the encrypted chunk and its metadata to the output stream.
func writeChunk(encryptedChunk, chachaNonce []byte, output io.Writer) error {
	chunkSizeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(chunkSizeBuf, uint16(len(encryptedChunk)))

	if _, err := output.Write(chunkSizeBuf); err != nil {
		return err
	}
	if _, err := output.Write(chachaNonce); err != nil {
		return err
	}
	if _, err := output.Write(encryptedChunk); err != nil {
		return err
	}

	return nil
}

// readChunkMetadata reads the chunk size and ChaCha20-Poly1305 nonce from the input stream.
func readChunkMetadata(input io.Reader) (uint16, []byte, error) {
	chunkSizeBuf := make([]byte, 2)
	if _, err := io.ReadFull(input, chunkSizeBuf); err != nil {
		if err == io.EOF {
			return 0, nil, err
		}
		return 0, nil, err
	}
	chunkSize := binary.BigEndian.Uint16(chunkSizeBuf)

	chachaNonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err := io.ReadFull(input, chachaNonce); err != nil {
		if err == io.EOF {
			return 0, nil, err
		}
		return 0, nil, err
	}

	return chunkSize, chachaNonce, nil
}
