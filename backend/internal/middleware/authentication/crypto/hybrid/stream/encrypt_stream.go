// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"io"
)

// Encrypt reads from the input stream, encrypts the data using AES-CTR and XChaCha20-Poly1305,
// calculates the HMAC if enabled, and writes it to the output stream.
//
// Note: This function requires a builder for the output, such as a string builder, rune builder, or byte builder,
// since it performs low-level operations on I/O primitives.
func (s *Stream) Encrypt(input io.Reader, output io.Writer) error {
	chunk := make([]byte, chunkSize)
	for {
		n, err := input.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			if err := s.encryptAndWriteChunk(chunk[:n], output); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
