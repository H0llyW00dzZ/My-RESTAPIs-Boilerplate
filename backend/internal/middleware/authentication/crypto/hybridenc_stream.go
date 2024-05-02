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

// HybridEncryptStream reads from the input stream, encrypts the data using AES-CTR and ChaCha20-Poly1305, and writes it to the output stream.
func HybridEncryptStream(input io.Reader, output io.Writer, aesKey, chachaKey []byte) error {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}

	chacha, err := chacha20poly1305.New(chachaKey)
	if err != nil {
		return err
	}

	chunk := make([]byte, chunkSize)
	for {
		n, err := input.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		chachaNonce, encryptedChunk, err := encryptChunk(aesBlock, chacha, chunk[:n])
		if err != nil {
			return err
		}

		// Write the size of the encrypted chunk to the output stream.
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

		if err == io.EOF {
			break
		}
	}

	return nil
}
