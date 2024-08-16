// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package restime

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// New creates a new instance of the ResponseTime middleware with the provided configuration.
// If no configuration is provided, the default configuration will be used.
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := DefaultConfig

	// Override default config with provided configuration
	if len(config) > 0 {
		if config[0].HeaderName != "" {
			cfg.HeaderName = config[0].HeaderName
		}
		if config[0].Next != nil {
			cfg.Next = config[0].Next
		}
	}

	return func(c *fiber.Ctx) error {
		// Check if the request should be skipped
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Start the timer
		start := time.Now()

		// Execute the next middleware/handler
		err := c.Next()

		// Calculate the response time
		responseTime := time.Since(start).Milliseconds()

		// Set the response time header only if the request is done
		if err == nil {
			c.Set(cfg.HeaderName, fmt.Sprintf("%dms", responseTime))
		}

		return err
	}
}
