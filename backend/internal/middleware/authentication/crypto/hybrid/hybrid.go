// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package hybrid

import "errors"

var (
	// ErrorInvalidCookie is a custom error variable that represents an error
	// which occurs when the cookie format is invalid or malformed.
	ErrorInvalidCookie = errors.New("invalid cookie format")
	// ErrorInvalidKey is a custom error variable that represents an error
	// which occurs when the provided key is invalid or cannot be decoded.
	ErrorInvalidKey = errors.New("invalid key")
)

// Service is the interface for the hybrid encryption service.
// It provides methods for cookie encryption and decryption using a hybrid encryption scheme.
type Service interface {
	// EncryptCookie encrypts a cookie value using a hybrid encryption scheme.
	// It takes the cookie value as input and returns the base64-encoded encrypted cookie.
	EncryptCookie(value string) (string, error)

	// DecryptCookie decrypts a cookie value using a hybrid decryption scheme.
	// It takes the base64-encoded encrypted cookie as input and returns the decrypted cookie value.
	DecryptCookie(encodedCookie string) (string, error)
}

// cryptoService is an implementation of the hybrid encryption Service interface.
type cryptoService struct {
	key string
}

// New creates a new instance of the hybrid encryption service.
// It takes the encryption key as input.
//
// TODO: Support Multiple Encoding for the encrypt value not a key (e.g., md5 which is suitable or other).
// Also, note that this encryption is strong, unlike JWT that can still lead to high vulnerability 💀.
func New(key string) Service {
	return &cryptoService{
		key: key,
	}
}

// EncryptCookie encrypts a cookie value using a hybrid encryption scheme.
// It takes the cookie value as input and returns the base64-encoded encrypted cookie.
func (s *cryptoService) EncryptCookie(value string) (string, error) {
	return EncryptCookie(value, s.key)
}

// DecryptCookie decrypts a cookie value using a hybrid decryption scheme.
// It takes the base64-encoded encrypted cookie as input and returns the decrypted cookie value.
func (s *cryptoService) DecryptCookie(encodedCookie string) (string, error) {
	return DecryptCookie(encodedCookie, s.key)
}
