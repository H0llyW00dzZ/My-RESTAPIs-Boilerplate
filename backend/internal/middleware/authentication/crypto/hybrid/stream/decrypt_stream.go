// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package stream

import (
	"io"
)

// Decrypt reads from the input stream, decrypts the data using XChaCha20-Poly1305 and AES-CTR,
// verifies the HMAC if enabled, and writes it to the output stream.
//
// Note: This function requires a builder for the output, such as a string builder, rune builder, or byte builder,
// since it performs low-level operations on I/O primitives. It is designed as the core of cryptographic operations and is compatible with
// the standard library.
func (s *Stream) Decrypt(input io.Reader, output io.Writer) error {
	for {
		chunk, err := s.readAndDecryptChunk(input)
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
