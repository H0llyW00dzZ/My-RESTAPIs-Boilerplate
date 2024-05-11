// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package bcrypt

// Service is the interface for the bcrypt password hashing service.
// It provides methods for password hashing and comparison.
type Service interface {
	// HashPassword takes a plaintext password and returns the bcrypt hash of the password.
	HashPassword(password string) (string, error)

	// ComparePassword compares a plaintext password with the stored bcrypt hash.
	// It returns true if the password matches the hash, false otherwise.
	ComparePassword(password, hash string) bool
}

// Hash is an implementation of the bcrypt password hashing Service interface.
type Hash struct{}

// New creates a new instance of the bcrypt password hashing service.
func New() Service {
	return &Hash{}
}

// HashPassword takes a plaintext password and returns the bcrypt hash of the password.
func (s *Hash) HashPassword(password string) (string, error) {
	return HashPassword(password)
}

// ComparePassword compares a plaintext password with the stored bcrypt hash.
// It returns true if the password matches the hash, false otherwise.
func (s *Hash) ComparePassword(password, hash string) bool {
	return ComparePassword(password, hash)
}
