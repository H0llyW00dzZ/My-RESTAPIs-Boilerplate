// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/aes"
	"encoding/binary"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

// HybridDecryptStream reads from the input stream, decrypts the data using ChaCha20-Poly1305 and AES-CTR, and writes it to the output stream.
func HybridDecryptStream(input io.Reader, output io.Writer, aesKey, chachaKey []byte) error {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}

	chacha, err := chacha20poly1305.New(chachaKey)
	if err != nil {
		return err
	}

	for {
		// Read the size of the encrypted chunk from the input stream.
		chunkSizeBuf := make([]byte, 2)
		if _, err := io.ReadFull(input, chunkSizeBuf); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		chunkSize := binary.BigEndian.Uint16(chunkSizeBuf)

		chachaNonce := make([]byte, chacha.NonceSize())
		if _, err := io.ReadFull(input, chachaNonce); err != nil {
			return err
		}

		encryptedChunk := make([]byte, chunkSize)
		if _, err := io.ReadFull(input, encryptedChunk); err != nil {
			return err
		}

		chunk, err := decryptChunk(aesBlock, chacha, chachaNonce, encryptedChunk)
		if err != nil {
			return err
		}

		if _, err := output.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}
