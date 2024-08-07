// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package stream

import (
	"io"
)

// Encrypt reads from the input stream, encrypts the data using AES-CTR and XChaCha20-Poly1305,
// calculates the HMAC if enabled, and writes it to the output stream.
//
// Note: This function requires a builder for the output, such as a string builder, rune builder, or byte builder,
// since it performs low-level operations on I/O primitives. It is designed as the core of cryptographic operations and is compatible with
// the standard library.
func (s *Stream) Encrypt(input io.Reader, output io.Writer) error {
	chunk := make([]byte, ChunkSize)
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
