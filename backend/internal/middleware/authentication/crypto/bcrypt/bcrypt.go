// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

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
type Hash struct {
	cost int
}

var (
	// ErrInvalidCost is returned when the provided cost is outside the allowed range.
	ErrInvalidCost = errors.New("bcrypt: invalid cost")
)

// New creates a new instance of the bcrypt password hashing service.
// If the cost is not provided or is outside the allowed range (MinCost to MaxCost),
// it will return an error.
func New(cost ...int) (Service, error) {
	h := &Hash{}
	if len(cost) > 0 {
		h.cost = cost[0]
	}
	if h.cost < bcrypt.MinCost {
		h.cost = bcrypt.DefaultCost
	}
	if h.cost > bcrypt.MaxCost {
		return nil, ErrInvalidCost
	}
	return h, nil
}

// HashPassword takes a plaintext password and returns the bcrypt hash of the password.
func (b *Hash) HashPassword(password string) (string, error) {
	return b.hashPassword(password)
}

// ComparePassword compares a plaintext password with the stored bcrypt hash.
// It returns true if the password matches the hash, false otherwise.
func (b *Hash) ComparePassword(password, hash string) bool {
	return b.comparePassword(password, hash)
}
