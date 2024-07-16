// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

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
