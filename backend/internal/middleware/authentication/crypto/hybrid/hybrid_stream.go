// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package hybrid

import (
	"bytes"
	"encoding/hex"
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
type streamService struct {
	aesKey    []byte
	chachaKey []byte
}

// NewStreamService creates a new instance of the stream-based encryption service.
// It takes the AES and ChaCha20-Poly1305 keys as input.
func NewStreamService(aesKey, chachaKey []byte) StreamService {
	return &streamService{
		aesKey:    aesKey,
		chachaKey: chachaKey,
	}
}

// Encrypt encrypts data using a hybrid encryption scheme with streams.
// It takes the data as input and returns the hex-encoded encrypted data.
func (s *streamService) Encrypt(data string) (string, error) {
	input := bytes.NewBufferString(data)
	encryptedOutput := &bytes.Buffer{}

	err := stream.EncryptStream(input, encryptedOutput, s.aesKey, s.chachaKey)
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
	decryptedOutput := &bytes.Buffer{}

	err = stream.DecryptStream(encryptedInput, decryptedOutput, s.aesKey, s.chachaKey)
	if err != nil {
		return "", err
	}

	return decryptedOutput.String(), nil
}
