// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// AES-CTR nonce size (16 bytes for a 128-bit nonce)
	//
	// TODO: Do we really need to increase this since the current size is still secure?
	aesNonceSize = 16
	minChunkBuf  = 2
)

const (
	// ChunkSize is the size of each chunk to be encrypted and decrypted.
	// This value should be large enough to improve performance but small enough to
	// avoid excessive memory consumption.
	ChunkSize = 1024
)

// encryptChunk encrypts a single chunk using AES-CTR and XChaCha20-Poly1305.
func (s *Stream) encryptChunk(chunk []byte) ([]byte, []byte, error) {
	// Generate a nonce for AES-CTR.
	aesNonceCapacity := s.AESNonceCapacity(len(chunk))
	aesNonce := make([]byte, aesNonceSize, aesNonceCapacity)
	if _, err := rand.Read(aesNonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the chunk using AES-CTR.
	aesStream := s.cipher(aesNonce)
	aesEncryptedChunk := make([]byte, len(chunk), aesNonceCapacity)
	aesStream.XORKeyStream(aesEncryptedChunk, chunk)

	// Prepend the AES nonce to the AES-CTR encrypted chunk.
	aesEncryptedChunkWithNonce := append(aesNonce, aesEncryptedChunk...)

	// Generate a nonce and anti-tamper for XChaCha20-Poly1305.
	//
	// Note: This is a technique following the example from the Go documentation
	// (see https://pkg.go.dev/golang.org/x/crypto@v0.24.0/chacha20poly1305#example-NewX).
	// By using a larger capacity for the nonce slice, the output of the nonce along with the
	// cryptographic randomness will always be unique, instead of using a fixed size from s.chacha.NonceSize().
	// This approach is suitable for XChaCha20-Poly1305.
	antiTamper := s.ChachaNonceCapacity(len(aesEncryptedChunkWithNonce))
	anotherChachaNonce := make([]byte, s.chacha.NonceSize(), antiTamper)
	if _, err := rand.Read(anotherChachaNonce); err != nil {
		return nil, nil, err
	}

	// Generate a nonce for XChaCha20-Poly1305.
	chachaNonce := make([]byte, s.chacha.NonceSize(), antiTamper)
	if _, err := rand.Read(chachaNonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the AES-CTR encrypted chunk (including the AES nonce) using XChaCha20-Poly1305.
	//
	// TODO: Consider including the HMAC sum of the AES-CTR encrypted chunk in the "additionalData" parameter.
	//       However, it is not strictly necessary at the moment since XChaCha20-Poly1305 is capable of handling
	//       up to 250GB of data, basically depending on the available memory (RAM) for most use-cases.
	chachaEncryptedChunk := s.chacha.Seal(anotherChachaNonce, chachaNonce, aesEncryptedChunkWithNonce, nil)

	// The resulting ciphertext will always be unique now.
	return chachaNonce, chachaEncryptedChunk, nil
}

// decryptChunk decrypts a single chunk using XChaCha20-Poly1305 and AES-CTR.
func (s *Stream) decryptChunk(chachaNonce, chachaEncryptedChunk []byte) ([]byte, error) {
	// Decrypt the chunk using XChaCha20-Poly1305.
	chachaEncrypted := chachaEncryptedChunk[s.chacha.NonceSize():]
	aesEncryptedChunk, err := s.chacha.Open(nil, chachaNonce, chachaEncrypted, nil)
	if err != nil {
		return nil, err
	}

	// Extract the AES nonce from the beginning of the AES-CTR encrypted chunk.
	aesNonce := aesEncryptedChunk[:aesNonceSize]
	aesEncryptedChunk = aesEncryptedChunk[aesNonceSize:]

	// Decrypt the chunk using AES-CTR.
	aesStream := s.cipher(aesNonce)
	chunk := make([]byte, len(aesEncryptedChunk))
	aesStream.XORKeyStream(chunk, aesEncryptedChunk)

	return chunk, nil
}

// encryptAndWriteChunk encrypts a chunk, calculates the HMAC (if enabled), and writes it to the output stream.
func (s *Stream) encryptAndWriteChunk(chunk []byte, output io.Writer) error {
	chachaNonce, encryptedChunk, err := s.encryptChunk(chunk)
	if err != nil {
		return err
	}

	// Note: The design is intentionally different from the [Digest] function exposed in the public API.
	//       The internal HMAC sum and appending process is separate from the [Digest] function
	//       to maintain a clear distinction between the internal and external functionality.
	//
	//       Also note that when HMAC authentication is enabled, the HMAC is bound to the encrypted chunk.
	//       Tampering with the HMAC (e.g., modifying it and then decrypting without HMAC authentication) will cause the decryption process to fail.
	//       The integrity verification will fail during decryption if the internal HMAC sum has been tampered with.
	//       It is crucial to maintain the integrity of the HMAC and not attempt any modifications.
	if s.hmac != nil {
		s.hmac.Reset()
		s.hmac.Write(encryptedChunk)
		hmacDigest := s.hmac.Sum(nil)
		encryptedChunk = append(encryptedChunk, hmacDigest...)
	}

	if err := s.writeChunk(encryptedChunk, chachaNonce, output); err != nil {
		return err
	}

	return nil
}

// readAndDecryptChunk reads an encrypted chunk from the input stream, verifies the HMAC (if enabled), and decrypts it.
func (s *Stream) readAndDecryptChunk(input io.Reader) ([]byte, error) {
	chunkSize, chachaNonce, err := s.readChunkMetadata(input)
	if err != nil {
		return nil, err
	}

	encryptedChunk, err := s.readEncryptedChunk(input, chunkSize)
	if err != nil {
		return nil, err
	}

	hmacDigest, encryptedChunk, err := s.extractHMACDigest(encryptedChunk, chunkSize)
	if err != nil {
		return nil, err
	}

	chunk, err := s.decryptChunk(chachaNonce, encryptedChunk)
	if err != nil {
		return nil, err
	}

	if s.hmac != nil {
		if err := s.verifyHMAC(encryptedChunk, hmacDigest); err != nil {
			return nil, err
		}
	}

	return chunk, nil
}

// writeChunk writes the encrypted chunk and its metadata to the output stream.
//
// TODO: Add an identifier to the chunk metadata. Since this is an object chunk, it is possible to include an identifier.
// For example, when inspecting the binary data, it might look like:
//
//	+---------------------------------------+
//	|                METADATA               |
//	+-------------+-------------------------+
//	|             |        Contestant       |
//	| Competition +-------+--------+--------+
//	|             |  John | Andrea | Robert |
//	+-------------+-------+--------+--------+
//	| Swimming    |  1:30 |   2:05 |   1:15 |
//	+-------------+-------+--------+--------+
//	| Running     | 15:30 |  14:10 |  15:45 |
//	+-------------+-------+--------+--------+
//
// Another example:
//
//	     ID   UUID   TEXT
//		---- ------ ------
//		 1    Text   Text
//
// Also note that these TODOs won't break the cipher text because they are outside the encrypted data.
func (s *Stream) writeChunk(encryptedChunk, chachaNonce []byte, output io.Writer) error {
	chunkSizeBuf := make([]byte, minChunkBuf)
	binary.BigEndian.PutUint16(chunkSizeBuf, uint16(len(encryptedChunk)))
	// Note: This example demonstrates how to use an identifier in SSL/TLS for a private network that can be used in a private server (not locally).
	// The focus is not solely on SSL/TLS (due boring of SSL/TLS), as this is the core of ciphertext cryptography.
	//
	//		identifier := []byte{0x00, 0x01337} // Example CipherSuite Identifier for H0LLYW00DZZ_AESCTR_XCHACHA20POLY1305_HMAC_SHA256
	//		if _, err := output.Write(identifier); err != nil {
	//	    	return err
	//		}

	if _, err := output.Write(chunkSizeBuf); err != nil {
		return err
	}
	if _, err := output.Write(chachaNonce); err != nil {
		return err
	}
	if _, err := output.Write(encryptedChunk); err != nil {
		return err
	}

	return nil
}

// readChunkMetadata reads the chunk size and XChaCha20-Poly1305 nonce from the input stream.
func (s *Stream) readChunkMetadata(input io.Reader) (uint16, []byte, error) {
	// Note: This example demonstrates how to use an identifier in SSL/TLS for a private network that can be used in a private server (not locally).
	// The focus is not solely on SSL/TLS (due boring of SSL/TLS), as this is the core of ciphertext cryptography.
	//
	//		identifierBuf := make([]byte, 2)
	//		if _, err := io.ReadAtLeast(input, identifierBuf, 2); err != nil {
	//	    	return 0, nil, err
	//		}
	//		identifier := binary.BigEndian.Uint16(identifierBuf)
	//		// Parse the identifier and set the corresponding algorithms and parameters
	//		switch identifier {
	//		case 0x01337:
	//	    	// H0LLYW00DZZ_AESCTR_XCHACHA20POLY1305_HMAC_SHA256
	//	    	// Set the corresponding algorithms and parameters
	//		default:
	//	    	return 0, nil, errors.New("Unknown identifier")
	//		}

	chunkSizeBuf := make([]byte, minChunkBuf)
	if _, err := io.ReadAtLeast(input, chunkSizeBuf, minChunkBuf); err != nil {
		if err == io.ErrUnexpectedEOF {
			return 0, nil, errors.New("XChaCha20-Poly1305: Unexpected Chunk Buffer Size")
		}
		return 0, nil, err
	}
	chunkSize := binary.BigEndian.Uint16(chunkSizeBuf)

	chachaNonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err := io.ReadAtLeast(input, chachaNonce, chacha20poly1305.NonceSizeX); err != nil {
		if err == io.ErrUnexpectedEOF {
			return 0, nil, errors.New("XChaCha20-Poly1305: Unexpected NonceSizeX")
		}
		return 0, nil, err
	}

	return chunkSize, chachaNonce, nil
}

// readEncryptedChunk reads the encrypted chunk from the input stream.
func (s *Stream) readEncryptedChunk(input io.Reader, chunkSize uint16) ([]byte, error) {
	encryptedChunk := make([]byte, chunkSize)
	if _, err := io.ReadFull(input, encryptedChunk); err != nil {
		if err == io.ErrUnexpectedEOF {
			if len(encryptedChunk) > 0 && s.hmac != nil {
				return nil, errors.New("XChaCha20-Poly1305: invalid HMAC digest size") // Middle Error Location in I/O primitives
			}
			return nil, errors.New("XChaCha20-Poly1305: encrypted chunk size mismatch") // Middle Error Location in I/O primitives
		}
		return nil, err
	}
	return encryptedChunk, nil
}

// extractHMACDigest extracts the HMAC digest from the encrypted chunk if HMAC is enabled.
func (s *Stream) extractHMACDigest(encryptedChunk []byte, chunkSize uint16) ([]byte, []byte, error) {
	// Note: This improves security by making it extremely difficult to tamper with the encrypted data without being detected.
	var hmacDigest []byte
	if s.hmac != nil {
		hmacDigestSize := s.hmac.Size()
		if len(encryptedChunk) < hmacDigestSize {
			// TODO: Use math to handle this error differently since it was uncovered in a test and performs low-level operations on I/O primitives.
			// Note: This uncovered probably in TLS, but it's not possible for TLS to encounter this as TLS already uses HMAC.
			return nil, nil, errors.New("Hybrid Scheme: invalid HMAC digest size") // Deep/Unknown Error Location in I/O primitives, probably in TLS
		}
		hmacDigest = encryptedChunk[len(encryptedChunk)-hmacDigestSize:]
		encryptedChunk = encryptedChunk[:len(encryptedChunk)-hmacDigestSize]
	} else {
		// If HMAC is not enabled, check if the encrypted chunk size matches the expected size
		if len(encryptedChunk) != int(chunkSize) {
			// TODO: Use math to handle this error differently since it was uncovered in a test and performs low-level operations on I/O primitives.
			// Note: This uncovered probably in TLS, but it's not possible for TLS to encounter this as TLS already uses HMAC.
			return nil, nil, errors.New("Hybrid Scheme: encrypted chunk size mismatch") // Deep/Unknown Error Location in I/O primitives, probably in TLS
		}
	}
	return hmacDigest, encryptedChunk, nil
}

// verifyHMAC verifies the HMAC of the encrypted chunk.
func (s *Stream) verifyHMAC(encryptedChunk, hmacDigest []byte) error {
	s.hmac.Reset()
	s.hmac.Write(encryptedChunk)
	expectedHMACDigest := s.hmac.Sum(nil)
	if subtle.ConstantTimeCompare(hmacDigest, expectedHMACDigest) != 1 {
		return errors.New("XChaCha20-Poly1305: HMAC verification failed")
	}
	return nil
}
