// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package connectionlogger

import (
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler that logs the current number of active connections.
//
// Note: This [connectionlogger] should be placed at the root router first, then put a logger middleware as documented here: https://docs.gofiber.io/api/middleware/logger
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := DefaultConfig

	// Override default config with provided configuration
	if len(config) > 0 {
		if config[0].Next != nil {
			cfg.Next = config[0].Next
		}
	}

	return func(c *fiber.Ctx) error {
		// Check if the request should be skipped
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Increment the active connection count
		//
		// Note: This is safe for concurrent use. However, using a mutex can decrease performance, so it's not recommended (too bad using mutex).
		atomic.AddInt64(&activeConnections, 1)
		defer func() {
			// Decrement the active connection count when the request is done
			atomic.AddInt64(&activeConnections, -1)
		}()

		// Process the request
		return c.Next()
	}
}
