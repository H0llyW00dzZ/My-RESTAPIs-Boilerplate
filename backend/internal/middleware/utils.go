// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
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
