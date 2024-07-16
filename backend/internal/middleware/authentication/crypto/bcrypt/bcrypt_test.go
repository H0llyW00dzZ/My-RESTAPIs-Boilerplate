// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package bcrypt_test

import (
	"testing"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/bcrypt"
)

func TestBcryptService_HashPassword(t *testing.T) {
	bcryptService, err := bcrypt.New()
	if err != nil {
		t.Fatalf("Failed to create bcrypt service: %v", err)
	}
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
	bcryptService, err := bcrypt.New()
	if err != nil {
		t.Fatalf("Failed to create bcrypt service: %v", err)
	}
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

func TestBcryptService_HashPasswordWithCustomCost(t *testing.T) {
	bcryptService, err := bcrypt.New(12)
	if err != nil {
		t.Fatalf("Failed to create bcrypt service: %v", err)
	}
	password := "password123"

	hashedPassword, err := bcryptService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}
}

func TestBcryptService_ComparePasswordWithCustomCost(t *testing.T) {
	bcryptService, err := bcrypt.New(12)
	if err != nil {
		t.Fatalf("Failed to create bcrypt service: %v", err)
	}
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

func TestBcryptService_InvalidCost(t *testing.T) {
	_, err := bcrypt.New(32)
	if err == nil {
		t.Error("Expected an error for invalid cost, but got nil")
	}
	if err != bcrypt.ErrInvalidCost {
		t.Errorf("Expected ErrInvalidCost, but got: %v", err)
	}
}
