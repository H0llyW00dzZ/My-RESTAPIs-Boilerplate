// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

// Service is the interface for the crypto service.
// It provides methods for encryption, decryption, and ciphertext verification.
type Service interface {
	// Encrypt encrypts the given data using AES encryption with a derived encryption key
	// and signs the ciphertext.
	// It returns the base64-encoded ciphertext and signature.
	Encrypt(data string) (string, string, error)

	// Decrypt decrypts the given encrypted data using AES decryption with the same derived
	// encryption key used during encryption.
	// It expects the encrypted data and signature to be base64-encoded.
	// It returns the decrypted data as a string.
	Decrypt(encryptedData, signature string) (string, error)

	// VerifyCiphertext verifies the integrity of the ciphertext without decrypting it.
	// It checks if the ciphertext has a valid structure and matches the expected format.
	// It expects the encrypted data and signature to be base64-encoded.
	// It returns true if the ciphertext is valid, false otherwise.
	VerifyCiphertext(encryptedData, signature string) bool
}

// cryptoService is an implementation of the crypto Service interface.
type cryptoService struct {
	useArgon2  bool
	secryptKey string
	signKey    string
}

// New creates a new instance of the crypto service.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the encryption key.
func New(useArgon2 bool, secryptKey, signKey string) Service {
	// TODO: Implement additional crypto algorithms based on my learning and knowledge.
	// Consider implementing the "labyrinth algorithm" (hahaha) or other cryptographic techniques.
	return &cryptoService{
		useArgon2:  useArgon2,
		secryptKey: secryptKey,
		signKey:    signKey,
	}
}

// Encrypt encrypts the given data using AES encryption with a derived encryption key
// and signs the ciphertext.
// It returns the base64-encoded ciphertext and signature.
func (s *cryptoService) Encrypt(data string) (string, string, error) {
	return EncryptData(data, s.useArgon2, s.secryptKey, s.signKey)
}

// Decrypt decrypts the given encrypted data using AES decryption with the same derived
// encryption key used during encryption.
// It expects the encrypted data and signature to be base64-encoded.
// It returns the decrypted data as a string.
func (s *cryptoService) Decrypt(encryptedData, signature string) (string, error) {
	return DecryptData(encryptedData, signature, s.useArgon2, s.secryptKey, s.signKey)
}

// VerifyCiphertext verifies the integrity of the ciphertext without decrypting it.
// It checks if the ciphertext has a valid structure and matches the expected format.
// It expects the encrypted data and signature to be base64-encoded.
// It returns true if the ciphertext is valid, false otherwise.
func (s *cryptoService) VerifyCiphertext(encryptedData, signature string) bool {
	return VerifyCiphertext(encryptedData, signature, s.signKey)
}
