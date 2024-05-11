// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword compares a plaintext password with the stored bcrypt hash.
// It returns true if the password matches the hash, false otherwise.
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
