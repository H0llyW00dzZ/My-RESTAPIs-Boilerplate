// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"fmt"
	"h0llyw00dz-template/backend/internal/database"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"hash/fnv"
	"strings"
	"time"

	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/google/uuid"
)

// NewCacheMiddleware creates a new cache middleware with the specified expiration time, cache control flag,
// and an optional custom key generator. It retrieves the Redis storage interface from the provided database
// service and configures the cache middleware accordingly.
//
// If a custom key generator is provided, it will be used to generate cache keys based on the request context.
// Otherwise, the default key generation mechanism of the Fiber cache middleware will be used, which generates
// cache keys based on the request method and path.
//
// Parameters:
//
//	db: The database service instance that provides the Redis storage interface.
//	expiration: The expiration time for cached entries.
//	cacheControl: A boolean flag indicating whether to include cache control headers in the response.
//	keyGenerator: An optional custom key generator function that takes the request context and returns a string key.
//
// Returns:
//
//	A Fiber handler function representing the configured cache middleware.
func NewCacheMiddleware(db database.Service, expiration time.Duration, cacheControl bool, keyGenerator ...func(*fiber.Ctx) string) fiber.Handler {
	// Retrieve the Redis storage interface from the database service.
	cacheMiddlewareService := db.FiberStorage()

	// Create a new cache middleware configuration.
	config := cache.Config{
		Expiration:   expiration,
		CacheControl: cacheControl,
		Storage:      cacheMiddlewareService,
	}

	// Check if a custom key generator is provided.
	if len(keyGenerator) > 0 {
		config.KeyGenerator = keyGenerator[0]
	}

	// Create a new cache middleware with the configured options.
	return cache.New(config)
}

// NewRateLimiter creates a new rate limiter middleware with the specified maximum number of requests,
// expiration time, and a custom message to log when the rate limit is reached.
// It retrieves the Redis storage interface from the provided database service and configures the rate limiter middleware accordingly.
func NewRateLimiter(db database.Service, max int, expiration time.Duration, limitReachedMessage string) fiber.Handler {
	// Retrieve the Redis storage interface from the database service.
	rateLimiterService := db.FiberStorage()
	// Create a new rate limiter middleware with the desired configuration.
	// TODO: Implement a custom key generator for any sensitive data such as API keys or OAuth tokens,
	// since the default rate limiter key in Fiber is based on c.IP()
	return limiter.New(limiter.Config{
		Storage:    rateLimiterService,
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c *fiber.Ctx) error {
			log.LogUserActivity(c, limitReachedMessage)
			return helper.SendErrorResponse(c, fiber.StatusTooManyRequests, fiber.ErrTooManyRequests.Message)
		},
	})
}

// NewCORSMiddleware creates a new CORS middleware with a better configuration.
// It allows specific origins, methods, headers, and credentials, and sets the maximum age for preflight requests.
func NewCORSMiddleware() fiber.Handler {
	// Note: In the Fiber framework v3, this CORS middleware configuration provides better security and low overhead.
	// For example, it allows blocking internal IPs by setting `AllowPrivateNetwork` to false (read more: https://docs.gofiber.io/api/middleware/cors).
	return cors.New(cors.Config{
		// Better Formatting
		AllowOrigins: strings.Join([]string{
			"https://example.com",
			"https://api.example.com",
		}, ","),
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: strings.Join([]string{
			"Content-Type",
			"Authorization",
		}, ","),
		ExposeHeaders: strings.Join([]string{
			"Content-Length",
		}, ","),
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	})
}

// NewETagMiddleware creates a new ETag middleware with the default configuration.
// It generates strong ETags for response caching and validation.
func NewETagMiddleware() fiber.Handler {
	return etag.New(etag.Config{
		Weak: false,
	})
}

// NewFaviconMiddleware creates a new favicon middleware to serve a favicon file.
// It takes the file path of the favicon and the URL path where the favicon will be served.
func NewFaviconMiddleware(filePath, urlPath string) fiber.Handler {
	return favicon.New(favicon.Config{
		File: filePath,
		URL:  urlPath,
	})
}

// NewPprofMiddleware creates a new pprof middleware with a custom path.
// It allows easy access to the pprof profiling tools and logs user activity.
func NewPprofMiddleware(path, pprofMessage string) fiber.Handler {
	// Example Usage: app.Use(NewPprofMiddleware("/pprof", "Accessed pprof profiling tools"))
	return func(c *fiber.Ctx) error {
		log.LogUserActivity(c, pprofMessage)
		return pprof.New(pprof.Config{
			Prefix: path,
		})(c)
	}
}

// NewIPBasedUUIDMiddleware creates a new middleware that generates a deterministic UUID based on the client's IP address
// and attaches it to the Fiber context for reusability.
func NewIPBasedUUIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the client's IP address from the Fiber context
		ipAddress := c.IP()

		// Generate a deterministic UUID based on the IP address
		uuid := generateGoogleUUIDFromIP(ipAddress)

		// Attach the generated UUID to the Fiber context
		c.Locals("ip_based_uuid", uuid)

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// generateGoogleUUIDFromIP generates a deterministic UUID based on the provided IP address.
func generateGoogleUUIDFromIP(ipAddress string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(ipAddress)).String()
}

// NewSignatureMiddleware creates a new middleware that generates a signature based on the client's IP address
// and attaches it to the Fiber context for security purposes.
func NewSignatureMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the client's IP address from the Fiber context
		ipAddress := c.IP()

		// Generate a signature based on the IP address
		signature := generateSignatureFromIP(ipAddress)

		// Attach the generated signature to the Fiber context
		c.Locals("signature", signature)

		// Continue to the next middleware or handler
		return c.Next()
	}
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

// CustomKeyGenerator generates a custom cache key based on the request and logs the visitor activity.
func CustomKeyGenerator(c *fiber.Ctx) string {
	// Get client's IP and User-Agent
	clientIP := c.IP()
	userAgent := c.Get(fiber.HeaderUserAgent)

	// Create a string to hash
	toHash := fmt.Sprintf("%s-%s", clientIP, userAgent)

	// Create a fnv hash and write our string to it
	signature := hashForSignature(toHash)

	// Generate a UUID based on the hash
	signatureUUID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(signature))

	// Log visitor activity with the signature for the frontend
	log.LogVisitorf("Frontend cache generated for visitor activity: IP [%s], User-Agent [%s], Signature [%s], UUID [%s]", clientIP, userAgent, signature, signatureUUID.String())

	// Generate a custom cache key with the hashed signature and UUID
	return fmt.Sprintf("cache_front_end:%s:%s:%s", signature, signatureUUID.String(), c.Path())
}
