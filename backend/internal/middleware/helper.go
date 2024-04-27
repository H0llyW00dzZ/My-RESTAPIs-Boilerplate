// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"fmt"
	"hash/fnv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/google/uuid"
)

// generateGoogleUUIDFromIP generates a deterministic UUID based on the provided IP address.
func generateGoogleUUIDFromIP(ipAddress string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(ipAddress)).String()
}

// generateSignatureFromIP generates a signature based on the provided IP address.
func generateSignatureFromIP(ipAddress string) string {
	// Generate a UUID based on the IP address
	uuid := uuid.NewSHA1(uuid.NameSpaceURL, []byte(ipAddress))

	// Generate a signature by taking the first 8 characters of the UUID
	signature := uuid.String()[:8]

	return signature
}

// hashForSignature creates a hash from the IP and User-Agent to use in generating a UUID.
func hashForSignature(toHash string) string {
	h := fnv.New64a()
	h.Write([]byte(toHash))
	return fmt.Sprintf("%x", h.Sum64())
}

// WithKeyGenerator is an option function for NewCacheMiddleware that sets a custom key generator.
func WithKeyGenerator(keyGenerator func(*fiber.Ctx) string) func(*cache.Config) {
	return func(config *cache.Config) {
		config.KeyGenerator = keyGenerator
	}
}

// WithCacheSkipper is an option function for NewCacheMiddleware that sets a custom cache skipper.
func WithCacheSkipper(cacheSkipper func(*fiber.Ctx) bool) func(*cache.Config) {
	return func(config *cache.Config) {
		config.Next = cacheSkipper
	}
}
