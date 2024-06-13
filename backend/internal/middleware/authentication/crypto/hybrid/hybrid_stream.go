// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package hybrid

import (
	"bytes"
	"encoding/hex"
	"strings"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
)

// StreamService is the interface for the stream-based encryption service.
type StreamService interface {
	// Encrypt encrypts data using a hybrid encryption scheme with streams.
	// It takes the data as input and returns the hex-encoded encrypted data.
	Encrypt(data string) (string, error)

	// Decrypt decrypts data using a hybrid decryption scheme with streams.
	// It takes the hex-encoded encrypted data as input and returns the decrypted data.
	Decrypt(encodedData string) (string, error)
}

// streamService is an implementation of the stream-based encryption StreamService interface.
// TODO: Adopt this high level of security for the stream-based encryption/decryption service, which is suitable for integration with Bubble Tea TUIs,
// by creating a separate repository and focusing on the client-side implementation rather than a web/server implementation.
//
// Reasons:
//
//   - Limited functionality: The current implementation may not be able to leverage advanced cryptography techniques.
//
//   - Cost: Implementing advanced cryptography techniques can be expensive in terms of resources usage. (for example, "argon2" the cost is arround 120MB round-trip tiket)
//
// Also, will add some additional features related to cryptography in the separate repository.
type streamService struct {
	stream *stream.Stream
}

// NewStreamService creates a new instance of the stream-based encryption service.
// It takes the AES and XChaCha20-Poly1305 keys as input.
func NewStreamService(aesKey, chachaKey []byte) (StreamService, error) {
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		return nil, err
	}

	return &streamService{
		stream: s,
	}, nil
}

// Encrypt encrypts data using a hybrid encryption scheme with streams.
// It takes the data as input and returns the hex-encoded encrypted data.
func (s *streamService) Encrypt(data string) (string, error) {
	input := bytes.NewBufferString(data)
	encryptedOutput := &bytes.Buffer{}

	err := s.stream.Encrypt(input, encryptedOutput)
	if err != nil {
		return "", err
	}

	encodedData := hex.EncodeToString(encryptedOutput.Bytes())
	return encodedData, nil
}

// Decrypt decrypts data using a hybrid decryption scheme with streams.
// It takes the hex-encoded encrypted data as input and returns the decrypted data.
func (s *streamService) Decrypt(encodedData string) (string, error) {
	encryptedData, err := hex.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	encryptedInput := bytes.NewBuffer(encryptedData)
	decryptedOutput := &strings.Builder{}

	err = s.stream.Decrypt(encryptedInput, decryptedOutput)
	if err != nil {
		return "", err
	}

	return decryptedOutput.String(), nil
}
