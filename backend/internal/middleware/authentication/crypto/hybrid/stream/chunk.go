// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// AES-CTR nonce size (16 bytes for a 128-bit nonce)
	//
	// TODO: Do we really need to increase this since the current size is still secure?
	aesNonceSize = 16
	chunkSize    = 1024
)

// encryptChunk encrypts a single chunk using AES-CTR and XChaCha20-Poly1305.
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

	// Generate a nonce for XChaCha20-Poly1305.
	chachaNonce := make([]byte, chacha.NonceSize())
	if _, err := rand.Read(chachaNonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the AES-CTR encrypted chunk (including the AES nonce) using XChaCha20-Poly1305.
	chachaEncryptedChunk := chacha.Seal(nil, chachaNonce, aesEncryptedChunkWithNonce, nil)

	return chachaNonce, chachaEncryptedChunk, nil
}

// decryptChunk decrypts a single chunk using XChaCha20-Poly1305 and AES-CTR.
func decryptChunk(aesBlock cipher.Block, chacha cipher.AEAD, chachaNonce, chachaEncryptedChunk []byte) ([]byte, error) {
	// Decrypt the chunk using XChaCha20-Poly1305.
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

// encryptAndWriteChunk encrypts a chunk, calculates the HMAC (if enabled), and writes it to the output stream.
func encryptAndWriteChunk(aesBlock cipher.Block, chacha cipher.AEAD, hmac hash.Hash, chunk []byte, output io.Writer) error {
	chachaNonce, encryptedChunk, err := encryptChunk(aesBlock, chacha, chunk)
	if err != nil {
		return err
	}

	if hmac != nil {
		hmac.Reset()
		hmac.Write(encryptedChunk)
		hmacDigest := hmac.Sum(nil)
		encryptedChunk = append(encryptedChunk, hmacDigest...)
	}

	if err := writeChunk(encryptedChunk, chachaNonce, output); err != nil {
		return err
	}

	return nil
}

// readAndDecryptChunk reads an encrypted chunk from the input stream, verifies the HMAC (if enabled), and decrypts it.
func readAndDecryptChunk(aesBlock cipher.Block, chacha cipher.AEAD, hmac hash.Hash, input io.Reader) ([]byte, error) {
	chunkSize, chachaNonce, err := readChunkMetadata(input)
	if err != nil {
		return nil, err
	}

	encryptedChunk := make([]byte, chunkSize)
	if _, err := io.ReadFull(input, encryptedChunk); err != nil {
		return nil, err
	}

	if hmac != nil {
		hmacDigestSize := hmac.Size()
		if len(encryptedChunk) < hmacDigestSize {
			return nil, errors.New("invalid HMAC digest size")
		}
		hmacDigest := encryptedChunk[len(encryptedChunk)-hmacDigestSize:]
		encryptedChunk = encryptedChunk[:len(encryptedChunk)-hmacDigestSize]

		hmac.Reset()
		hmac.Write(encryptedChunk)
		expectedHMACDigest := hmac.Sum(nil)
		if subtle.ConstantTimeCompare(hmacDigest, expectedHMACDigest) != 1 {
			return nil, errors.New("HMAC verification failed")
		}
	}

	chunk, err := decryptChunk(aesBlock, chacha, chachaNonce, encryptedChunk)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// writeChunk writes the encrypted chunk and its metadata to the output stream.
//
// TODO: Add an identifier to the chunk metadata. Since this is an object chunk, it is possible to include an identifier.
// For example, when inspecting the binary data, it might look like:
//
//	+---------------------------------------+
//	|                METADATA               |
//	+-------------+-------------------------+
//	|             |        Contestant       |
//	| Competition +-------+--------+--------+
//	|             |  John | Andrea | Robert |
//	+-------------+-------+--------+--------+
//	| Swimming    |  1:30 |   2:05 |   1:15 |
//	+-------------+-------+--------+--------+
//	| Running     | 15:30 |  14:10 |  15:45 |
//	+-------------+-------+--------+--------+
//
// Another example:
//
//	     ID   UUID   TEXT
//		---- ------ ------
//		 1    Text   Text
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

// readChunkMetadata reads the chunk size and XChaCha20-Poly1305 nonce from the input stream.
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
