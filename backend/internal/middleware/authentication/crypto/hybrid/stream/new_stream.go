// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"hash"

	"golang.org/x/crypto/chacha20poly1305"
)

// Stream represents a Hybrid stream encryption/decryption object.
//
// The Stream struct combines AES-CTR and XChaCha20-Poly1305 encryption algorithms
// to provide a secure and efficient way to encrypt and decrypt data streams.
// It also supports optional HMAC authentication for added integrity and authenticity.
//
// Online Tools for Cipher Analysis and Identification:
//
//   - Boxentriq Cipher Identifier: https://www.boxentriq.com/code-breaking/cipher-identifier
//     This online tool helps identify the type of cipher used based on the ciphertext.
//     It supports various classical and modern ciphers.
//
//   - Hex-Works: https://hex-works.com/
//     Hex-Works provides a set of online tools for working with hexadecimal data, including encryption,
//     decryption, and analysis. It supports AES, DES, RC4, and other ciphers.
//
// Note: The security of the encrypted data depends on the secure generation and management of the encryption keys.
// Make sure to use strong, randomly generated keys and keep them confidential.
//
// Also note that if the results from the above tools cannot identify the cipher used, it is considered
// a strong indication that your data and privacy are secure. If the tools fail to identify the cipher,
// it suggests that the encryption scheme is robust and resistant to common analysis techniques,
// providing a high level of confidentiality and security for your data.
type Stream struct {
	aesBlock       cipher.Block
	chacha         cipher.AEAD
	hmac           hash.Hash
	cipher         func([]byte) cipher.Stream
	customizeNonce *CustomizeCapacityNonce
}

const (
	// additionalCapacityPercentage represents the percentage of additional capacity
	// to be added to the nonce capacity when it exceeds the minimum required size.
	// This constant is used in both the AESNonceCapacity and ChachaNonceCapacity methods
	// to calculate the nonce capacity for AES-CTR and XChaCha20-Poly1305 respectively.
	// The default value is set to 0.05, which means an additional 5% capacity will be added.
	additionalCapacityPercentage = 0.05 // use 5% capacity
)

// CustomizeCapacityNonce allows customizing the nonce capacity for AES-CTR and XChaCha20-Poly1305.
// The default value for both AESNonceCapacity and ChachaNonceCapacity is 0.05 (5% additional capacity).
//
// Note: Customizing the nonce capacity for AES-CTR and XChaCha20-Poly1305 won't increase the bit size of the nonce.
// It just enhances the ciphertext to make it always unique. The default is set to 5%, which provides a balance between
// security and performance. If the capacity is set to a high value (e.g., 100%), it can lead to performance issues
// when encrypting or decrypting large amounts of data (e.g., 100GB data then it need 200GB).
type CustomizeCapacityNonce struct {
	AESNonceCapacity    float64
	ChachaNonceCapacity float64
}

// New creates a new Stream instance with the provided AES and XChaCha20-Poly1305 keys.
// HMAC authentication is disabled by default.
func New(aesKey, chachaKey []byte) (*Stream, error) {
	if len(aesKey) != 16 && len(aesKey) != 24 && len(aesKey) != 32 {
		return nil, errors.New("Hybrid Scheme: Invalid AES-CTR key size")
	}

	if len(chachaKey) != 32 {
		return nil, errors.New("Hybrid Scheme: Invalid XChaCha20-Poly1305 key size")
	}

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
		cipher: func(nonce []byte) cipher.Stream {
			return cipher.NewCTR(aesBlock, nonce)
		},
		customizeNonce: &CustomizeCapacityNonce{
			AESNonceCapacity:    additionalCapacityPercentage,
			ChachaNonceCapacity: additionalCapacityPercentage,
		},
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

// AESNonceCapacity calculates the nonce capacity for AES-CTR based on the length of the encrypted data.
// It takes the length of the encrypted data as input and returns the calculated nonce capacity.
func (s *Stream) AESNonceCapacity(encryptedLen int) int {
	return s.calculateAESNonceCapacity(encryptedLen)
}

// ChachaNonceCapacity calculates the nonce capacity for XChaCha20-Poly1305 based on the length of the encrypted data.
// It takes the length of the encrypted data as input and returns the calculated nonce capacity.
func (s *Stream) ChachaNonceCapacity(encryptedLen int) int {
	return s.calculateChachaNonceCapacity(s.chacha.NonceSize(), encryptedLen+s.chacha.Overhead())
}

// CustomizeNonceCapacity allows customizing the nonce capacity for AES-CTR and XChaCha20-Poly1305.
// The default value for both AESNonceCapacity and ChachaNonceCapacity is 0.05 (5% additional capacity).
//
// Example usage:
//
//	stream, err := stream.New(aesKey, chachaKey)
//	if err != nil {
//	    // Handle error
//	}
//
//	// Customize nonce capacity
//	stream.CustomizeNonceCapacity(0.1, 0.08)
//
// In this example, the nonce capacity for AES-CTR is set to 0.1 (10% additional capacity)
// and the nonce capacity for XChaCha20-Poly1305 is set to 0.08 (8% additional capacity).
func (s *Stream) CustomizeNonceCapacity(aesCapacity, chachaCapacity float64) {
	s.customizeNonce.AESNonceCapacity = aesCapacity
	s.customizeNonce.ChachaNonceCapacity = chachaCapacity
}
