// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"hash"

	"golang.org/x/crypto/chacha20poly1305"
)

// Stream represents a Hybrid stream encryption/decryption object.
type Stream struct {
	aesBlock cipher.Block
	chacha   cipher.AEAD
	hmac     hash.Hash
}

// New creates a new Stream instance with the provided AES and ChaCha20-Poly1305 keys.
// HMAC authentication is disabled by default.
func New(aesKey, chachaKey []byte) (*Stream, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	chacha, err := chacha20poly1305.NewX(chachaKey)
	if err != nil {
		return nil, err
	}

	return &Stream{
		aesBlock: aesBlock,
		chacha:   chacha,
	}, nil
}

// EnableHMAC enables HMAC authentication for the stream using the provided key.
func (s *Stream) EnableHMAC(key []byte) {
	s.hmac = hmac.New(sha256.New, key)
}
