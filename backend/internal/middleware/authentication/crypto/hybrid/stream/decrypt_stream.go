// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"io"
)

// Decrypt reads from the input stream, decrypts the data using ChaCha20-Poly1305 and AES-CTR,
// and writes it to the output stream.
func (s *Stream) Decrypt(input io.Reader, output io.Writer) error {
	for {
		chunk, err := readAndDecryptChunk(s.aesBlock, s.chacha, input)
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
