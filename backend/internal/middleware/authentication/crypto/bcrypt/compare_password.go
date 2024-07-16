// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword compares a plaintext password with the stored bcrypt hash.
// It returns true if the password matches the hash, false otherwise.
func (b *Hash) comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
