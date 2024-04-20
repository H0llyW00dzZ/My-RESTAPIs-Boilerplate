// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
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

// NewCacheMiddleware creates a new cache middleware with the specified expiration time and cache control flag.
// It retrieves the Redis storage interface from the provided database service and configures the cache middleware accordingly.
func NewCacheMiddleware(db database.Service, expiration time.Duration, cacheControl bool) fiber.Handler {
	// Retrieve the Redis storage interface from the database service.
	cacheMiddlewareService := db.FiberStorage()
	// Create a new cache middleware with the desired configuration.
	return cache.New(cache.Config{
		Expiration:   expiration,
		CacheControl: cacheControl,
		Storage:      cacheMiddlewareService,
	})
}

// NewRateLimiter creates a new rate limiter middleware with the specified maximum number of requests,
// expiration time, and a custom message to log when the rate limit is reached.
// It retrieves the Redis storage interface from the provided database service and configures the rate limiter middleware accordingly.
func NewRateLimiter(db database.Service, max int, expiration time.Duration, limitReachedMessage string) fiber.Handler {
	// Retrieve the Redis storage interface from the database service.
	rateLimiterService := db.FiberStorage()
	// Create a new rate limiter middleware with the desired configuration.
	return limiter.New(limiter.Config{
		Storage:    rateLimiterService,
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c *fiber.Ctx) error {
			log.LogUserActivity(c, limitReachedMessage)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.ErrTooManyRequests.Message,
			})
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
