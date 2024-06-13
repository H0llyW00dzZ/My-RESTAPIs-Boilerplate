// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"crypto/hmac"
	"crypto/sha256"
	"io"
)

// Digest calculates the HMAC digest of the encrypted data (if HMAC is enabled).
//
// Note: When HMAC is enabled, don't forget to save/log the HMAC sum into your pocket or whatever it is for future use.
func (s *Stream) Digest(input io.Reader) ([]byte, error) {
	if s.hmac == nil {
		return nil, nil
	}

	// TODO: Allow the use of other hash functions. However,
	// it is not really needed right now since SHA-256 is sufficient on my CPU.
	hmacHash := hmac.New(sha256.New, s.hmac.Sum(nil))

	if _, err := io.Copy(hmacHash, input); err != nil {
		return nil, err
	}

	return hmacHash.Sum(nil), nil
}
