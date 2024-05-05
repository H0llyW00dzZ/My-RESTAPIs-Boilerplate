// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/aes"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

// HybridDecryptStream reads from the input stream, decrypts the data using ChaCha20-Poly1305 and AES-CTR, and writes it to the output stream.
func HybridDecryptStream(input io.Reader, output io.Writer, aesKey, chachaKey []byte) error {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}

	chacha, err := chacha20poly1305.NewX(chachaKey)
	if err != nil {
		return err
	}

	for {
		chunk, err := readAndDecryptChunk(aesBlock, chacha, input)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if _, err := output.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}
