// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package bcrypt

import "golang.org/x/crypto/bcrypt"

// hashPassword takes a plaintext password and returns the bcrypt hash of the password.
func (b *Hash) hashPassword(password string) (string, error) {
	// Generate a salt with a cost from struct
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
