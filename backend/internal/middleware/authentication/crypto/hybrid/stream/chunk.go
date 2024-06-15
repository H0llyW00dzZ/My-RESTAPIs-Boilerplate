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
func (s *Stream) encryptChunk(chunk []byte) ([]byte, []byte, error) {
	// Generate a nonce for AES-CTR.
	aesNonce := make([]byte, aesNonceSize)
	if _, err := rand.Read(aesNonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the chunk using AES-CTR.
	aesStream := cipher.NewCTR(s.aesBlock, aesNonce)
	aesEncryptedChunk := make([]byte, len(chunk))
	aesStream.XORKeyStream(aesEncryptedChunk, chunk)

	// Prepend the AES nonce to the AES-CTR encrypted chunk.
	aesEncryptedChunkWithNonce := append(aesNonce, aesEncryptedChunk...)

	// Generate a nonce for XChaCha20-Poly1305.
	chachaNonce := make([]byte, s.chacha.NonceSize())
	if _, err := rand.Read(chachaNonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the AES-CTR encrypted chunk (including the AES nonce) using XChaCha20-Poly1305.
	chachaEncryptedChunk := s.chacha.Seal(nil, chachaNonce, aesEncryptedChunkWithNonce, nil)

	return chachaNonce, chachaEncryptedChunk, nil
}

// decryptChunk decrypts a single chunk using XChaCha20-Poly1305 and AES-CTR.
func (s *Stream) decryptChunk(chachaNonce, chachaEncryptedChunk []byte) ([]byte, error) {
	// Decrypt the chunk using XChaCha20-Poly1305.
	aesEncryptedChunk, err := s.chacha.Open(nil, chachaNonce, chachaEncryptedChunk, nil)
	if err != nil {
		return nil, err
	}

	// Extract the AES nonce from the beginning of the AES-CTR encrypted chunk.
	aesNonce := aesEncryptedChunk[:aesNonceSize]
	aesEncryptedChunk = aesEncryptedChunk[aesNonceSize:]

	// Decrypt the chunk using AES-CTR.
	aesStream := cipher.NewCTR(s.aesBlock, aesNonce)
	chunk := make([]byte, len(aesEncryptedChunk))
	aesStream.XORKeyStream(chunk, aesEncryptedChunk)

	return chunk, nil
}

// encryptAndWriteChunk encrypts a chunk, calculates the HMAC (if enabled), and writes it to the output stream.
func (s *Stream) encryptAndWriteChunk(chunk []byte, output io.Writer) error {
	chachaNonce, encryptedChunk, err := s.encryptChunk(chunk)
	if err != nil {
		return err
	}

	if s.hmac != nil {
		s.hmac.Reset()
		s.hmac.Write(encryptedChunk)
		hmacDigest := s.hmac.Sum(nil)
		encryptedChunk = append(encryptedChunk, hmacDigest...)
	}

	if err := s.writeChunk(encryptedChunk, chachaNonce, output); err != nil {
		return err
	}

	return nil
}

// readAndDecryptChunk reads an encrypted chunk from the input stream, verifies the HMAC (if enabled), and decrypts it.
func (s *Stream) readAndDecryptChunk(input io.Reader) ([]byte, error) {
	chunkSize, chachaNonce, err := s.readChunkMetadata(input)
	if err != nil {
		return nil, err
	}

	encryptedChunk := make([]byte, chunkSize)
	if _, err := io.ReadFull(input, encryptedChunk); err != nil {
		return nil, err
	}

	var hmacDigest []byte
	// Note This Improve making it extremely difficult to tamper with the encrypted data without being detected.
	if s.hmac != nil {
		hmacDigestSize := s.hmac.Size()
		if len(encryptedChunk) < hmacDigestSize {
			// TODO: This error uncovered in test, since it performs low-level operations on I/O primitives. error handle it's different.
			return nil, errors.New("invalid HMAC digest size")
		}
		hmacDigest = encryptedChunk[len(encryptedChunk)-hmacDigestSize:]
		encryptedChunk = encryptedChunk[:len(encryptedChunk)-hmacDigestSize]
	} else {
		// If HMAC is not enabled, check if the encrypted chunk size matches the expected size
		if len(encryptedChunk) != int(chunkSize) {
			// TODO: This error uncovered in test, since it performs low-level operations on I/O primitives. error handle it's different.
			return nil, errors.New("encrypted chunk size mismatch")
		}
	}

	chunk, err := s.decryptChunk(chachaNonce, encryptedChunk)
	if err != nil {
		return nil, err
	}

	if s.hmac != nil {
		s.hmac.Reset()
		s.hmac.Write(encryptedChunk)
		expectedHMACDigest := s.hmac.Sum(nil)
		if subtle.ConstantTimeCompare(hmacDigest, expectedHMACDigest) != 1 {
			return nil, errors.New("HMAC verification failed")
		}
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
//
// Also note that these TODOs won't break the cipher text because they are outside the encrypted data.
func (s *Stream) writeChunk(encryptedChunk, chachaNonce []byte, output io.Writer) error {
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
func (s *Stream) readChunkMetadata(input io.Reader) (uint16, []byte, error) {
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
