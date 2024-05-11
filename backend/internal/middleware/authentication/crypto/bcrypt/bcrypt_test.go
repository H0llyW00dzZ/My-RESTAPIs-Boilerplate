// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package bcrypt_test

import (
	"testing"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "password123"

	hashedPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}
}

func TestComparePassword(t *testing.T) {
	password := "password123"
	incorrectPassword := "incorrect"

	hashedPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !bcrypt.ComparePassword(password, hashedPassword) {
		t.Error("ComparePassword returned false for correct password")
	}

	if bcrypt.ComparePassword(incorrectPassword, hashedPassword) {
		t.Error("ComparePassword returned true for incorrect password")
	}
}

func TestBcryptService_HashPassword(t *testing.T) {
	bcryptService := bcrypt.New()
	password := "password123"

	hashedPassword, err := bcryptService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}
}

func TestBcryptService_ComparePassword(t *testing.T) {
	bcryptService := bcrypt.New()
	password := "password123"
	incorrectPassword := "incorrect"

	hashedPassword, err := bcryptService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !bcryptService.ComparePassword(password, hashedPassword) {
		t.Error("ComparePassword returned false for correct password")
	}

	if bcryptService.ComparePassword(incorrectPassword, hashedPassword) {
		t.Error("ComparePassword returned true for incorrect password")
	}
}
