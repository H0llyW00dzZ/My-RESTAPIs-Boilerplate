// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/aes"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

// HybridEncryptStream reads from the input stream, encrypts the data using AES-CTR and ChaCha20-Poly1305, and writes it to the output stream.
func HybridEncryptStream(input io.Reader, output io.Writer, aesKey, chachaKey []byte) error {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}

	chacha, err := chacha20poly1305.NewX(chachaKey)
	if err != nil {
		return err
	}

	chunk := make([]byte, chunkSize)
	for {
		n, err := input.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			if err := encryptAndWriteChunk(aesBlock, chacha, chunk[:n], output); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
