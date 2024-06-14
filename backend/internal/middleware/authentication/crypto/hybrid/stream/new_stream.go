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

// New creates a new Stream instance with the provided AES and XChaCha20-Poly1305 keys.
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
// The HMAC is computed using SHA-256 and follows the Encrypt-then-MAC (EtM) scheme approach
// as specified in (RFC 7366).
//
// When HMAC authentication is enabled, the encryption process is modified as follows:
//
//  1. The plaintext is encrypted using AES-CTR with a randomly generated nonce.
//  2. The AES-CTR encrypted data is then encrypted using XChaCha20-Poly1305 with another randomly generated nonce.
//  3. The HMAC is computed over the XChaCha20-Poly1305 encrypted data using the HMAC key.
//  4. The resulting HMAC tag is appended to the ciphertext.
//
// The decryption process is modified as follows:
//
//  1. The HMAC tag is extracted from the end of the ciphertext.
//  2. The HMAC is computed over the remaining ciphertext using the HMAC key.
//  3. The computed HMAC is compared with the extracted HMAC tag. If they don't match, an error is returned.
//  4. If the HMAC verification succeeds, the ciphertext is decrypted using XChaCha20-Poly1305.
//  5. The resulting data is then decrypted using AES-CTR.
//
// Note: The Digest method can be used to calculate the HMAC digest of the encrypted data for future verification.
// Since the HMAC is computed over both the AES-CTR and XChaCha20-Poly1305 encrypted data, the Digest method is optional
// and can be used to store the HMAC sum separately for additional verification purposes.
func (s *Stream) EnableHMAC(key []byte) {
	s.hmac = hmac.New(sha256.New, key)
}
