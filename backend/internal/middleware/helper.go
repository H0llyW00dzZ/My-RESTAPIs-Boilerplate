// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"hash/fnv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

// ratelimiterMsg is a custom handler function for the rate limiter middleware.
// It returns a closure that logs a message indicating that a visitor has been rate limited
// and sends an error response with a "Too Many Requests" status code and an appropriate error message.
//
// Parameters:
//
//	customMessage: A custom message to be logged when a visitor is rate limited.
//
// Returns:
//
//	A closure that takes a Fiber context and returns an error.
//	The closure logs the custom message and sends an error response indicating that the rate limit has been reached.
func ratelimiterMsg(customMessage string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.LogUserActivity(c, customMessage)
		// Note: Custom messages for the "Too Many Requests" HTTP response will not be implemented due to lazy 🤪.
		// The default error message provided by Fiber will be used instead.
		return helper.SendErrorResponse(c, fiber.StatusTooManyRequests, fiber.ErrTooManyRequests.Message)
	}
}

// WithKeyGenerator is an option function for NewCacheMiddleware and NewRateLimiter that sets a custom key generator.
// It takes a keyGenerator function that generates a unique key based on the Fiber context and returns a closure
// that configures the key generator for the specified middleware configuration.
//
// The function uses a type switch to determine the type of the middleware configuration and sets the appropriate
// key generator field. It supports cache.Config and limiter.Config types. If an unsupported config type is passed,
// it panics with an appropriate error message.
//
// The use of a type switch in this function is considered a better approach than using multiple if-else statements,
// as it provides a cleaner and more concise way to handle different configuration types, even if there are a large
// number of cases.
//
// Parameters:
//
//	keyGenerator: A function that takes a Fiber context and returns a unique key string.
//
// Returns:
//
//	A closure that takes a middleware configuration and sets the key generator based on the configuration type.
//
// TODO: Implement additional key generators for other middleware configurations besides NewCacheMiddleware and NewRateLimiter
func WithKeyGenerator(keyGenerator func(*fiber.Ctx) string) interface{} {
	return func(config interface{}) {
		// Note: This a better switch-statement, it doesn't matter if there is so many switch (e.g, 1 billion switch case)
		switch cfg := config.(type) {
		case *cache.Config:
			cfg.KeyGenerator = keyGenerator
			// TODO: Implement a custom key generator for any sensitive data such as API keys or OAuth tokens,
			// since the default rate limiter key in Fiber is based on c.IP()
		case *limiter.Config:
			cfg.KeyGenerator = keyGenerator
		default:
			panic(fmt.Sprintf("unsupported config type: %T", config))
		}
	}
}

// WithCacheSkipper is an option function for NewCacheMiddleware that sets a custom cache skipper.
func WithCacheSkipper(cacheSkipper func(*fiber.Ctx) bool) func(*cache.Config) {
	return func(config *cache.Config) {
		config.Next = cacheSkipper
	}
}

// WithValidator is an option function for NewKeyAuthMiddleware that sets a custom validator.
func WithValidator(validator func(*fiber.Ctx, string) (bool, error)) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.Validator = validator
	}
}

// WithErrorHandler is an option function for NewKeyAuthMiddleware that sets a custom error handler.
func WithErrorHandler(errorHandler func(*fiber.Ctx, error) error) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.ErrorHandler = errorHandler
	}
}

// WithSuccessHandler is an option function for NewKeyAuthMiddleware that sets a custom success handler.
func WithSuccessHandler(successHandler func(*fiber.Ctx) error) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.SuccessHandler = successHandler
	}
}

// WithKeyLookup is an option function for NewKeyAuthMiddleware that sets a custom key lookup.
func WithKeyLookup(keyLookup string) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.KeyLookup = keyLookup
	}
}

// WithContextKey is an option function for NewKeyAuthMiddleware that sets a custom context key.
func WithContextKey(contextKey string) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.ContextKey = contextKey
	}
}

// WithMax is an option function for NewRateLimiter that sets the maximum number of requests.
func WithMax(max int) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.Max = max
	}
}

// WithExpiration is an option function for NewRateLimiter that sets the expiration time.
func WithExpiration(expiration time.Duration) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.Expiration = expiration
	}
}

// WithLimitReached is an option function for NewRateLimiter that sets a custom limit reached handler.
func WithLimitReached(limitReached func(*fiber.Ctx) error) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.LimitReached = limitReached
	}
}
