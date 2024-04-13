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
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
	return cors.New(cors.Config{
		AllowOrigins: "https://example.com, https://api.example.com",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400, // 24 hours
	})
}
