// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: This should fix the type assertion bug (it's a language-specific problem).

package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// WithCacheStorage is an option function for cache middleware that sets the cache storage backend.
func WithCacheStorage(storage fiber.Storage) func(*cache.Config) {
	return func(config *cache.Config) {
		config.Storage = storage
	}
}

// WithSessionStorage is an option function for session middleware that sets the session storage backend.
func WithSessionStorage(storage fiber.Storage) func(*session.Config) {
	return func(config *session.Config) {
		config.Storage = storage
	}
}

// WithRateLimiterStorage is an option function for rate limiter middleware that sets the rate limiter storage backend.
func WithRateLimiterStorage(storage fiber.Storage) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.Storage = storage
	}
}

// WithCSRFStorage is an option function for CSRF middleware that sets the CSRF storage backend.
func WithCSRFStorage(storage fiber.Storage) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.Storage = storage
	}
}

// WithIdempotencyStorage is an option function for idempotency middleware that sets the idempotency storage backend.
func WithIdempotencyStorage(storage fiber.Storage) func(*idempotency.Config) {
	return func(config *idempotency.Config) {
		config.Storage = storage
	}
}
