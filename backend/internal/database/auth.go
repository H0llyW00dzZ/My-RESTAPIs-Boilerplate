// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package database

import (
	"database/sql"
	"errors"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

var (
	// ErrInvalidAPIKey is a custom error variable that represents an invalid API key.
	ErrInvalidAPIKey = errors.New("Invalid API key")
	// ErrExpiredAPIKey is a custom error variable that represents an API key expired.
	ErrExpiredAPIKey = errors.New("API Key Expired")
)

// ServiceAuth is an interface that defines methods for user authentication and management.
// It provides a contract for implementing user-related database operations.
type ServiceAuth interface {
	// FiberStorage returns the [fiber.Storage] interface for fiber storage middleware.
	FiberStorage() fiber.Storage
}

// serviceAuth is a concrete implementation of the ServiceAuth interface.
// It encapsulates the database connection and provides methods to interact with user data.
type serviceAuth struct {
	db           *sql.DB
	fiberStorage fiber.Storage
	bcrypt       bcrypt.Service
}

// NewServiceAuth creates a new instance of the ServiceAuth interface.
// It takes a database connection as a parameter and returns a new serviceAuth instance.
func NewServiceAuth(db *sql.DB, fiberStorage fiber.Storage, bcryptService bcrypt.Service) ServiceAuth {
	return &serviceAuth{
		db:           db,
		fiberStorage: fiberStorage,
		bcrypt:       bcryptService,
	}
}

// FiberStorage returns the fiber.Storage instance used for caching.
// This method provides access to the Redis storage instance that is used for caching
// API Keys in the auth middleware.
func (s *serviceAuth) FiberStorage() fiber.Storage {
	return s.fiberStorage
}
