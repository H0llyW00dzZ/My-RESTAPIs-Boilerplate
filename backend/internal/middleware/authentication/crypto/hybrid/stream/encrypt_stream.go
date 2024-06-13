// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"io"
)

// Encrypt reads from the input stream, encrypts the data using AES-CTR and ChaCha20-Poly1305,
// calculates the HMAC if enabled, and writes it to the output stream.
func (s *Stream) Encrypt(input io.Reader, output io.Writer) error {
	chunk := make([]byte, chunkSize)
	for {
		n, err := input.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			if err := encryptAndWriteChunk(s.aesBlock, s.chacha, s.hmac, chunk[:n], output); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
