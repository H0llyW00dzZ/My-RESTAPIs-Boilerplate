// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package csp

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
)

// Config represents the configuration options for the CSP middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// RandomnessGenerator is a function that generates a random string for the CSP nonce.
	// It takes a string parameter as input, which can be used as a seed or additional data for randomness generation.
	// The function should return the generated random string.
	//
	// Optional. Default: A default randomness generator function provided by the CSP middleware.
	RandomnessGenerator func(string) string

	// ContextKey is the key used to store and retrieve the CSP data from the request context.
	// It should be a unique and descriptive key to avoid conflicts with other middleware or application data.
	//
	// Optional. Default: "csp_random"
	ContextKey string
}

// DefaultConfig returns the default configuration for the CSP middleware.
func DefaultConfig() Config {
	return Config{
		RandomnessGenerator: defaultRandomnessGenerator,
		ContextKey:          "csp_random",
		Next:                nil,
	}
}

// defaultRandomnessGenerator generates randomness using SHA256 of the client IP.
func defaultRandomnessGenerator(clientIP string) string {
	hash := sha256.Sum256([]byte(clientIP))
	return hex.EncodeToString(hash[:])
}
