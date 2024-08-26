// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package csp

import (
	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
)

// New creates a new instance of the CSP middleware with the provided configuration.
func New(config ...Config) fiber.Handler {
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		// Test Skipped For This when Next True
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		clientIP := getClientIP(c, cfg.IPHeader)

		cloudflareRayID := c.Get(log.CloudflareRayIDHeader)
		if cloudflareRayID != "" {
			clientIP += " - Cloudflare detected - Ray ID: " + cloudflareRayID
		}
		countryCode := c.Get(log.CloudflareIPCountryHeader)
		if countryCode != "" {
			clientIP += ", Country: " + countryCode
		}

		// Generate the randomness using the configured generator
		randomness := cfg.RandomnessGenerator(clientIP)

		// Set the randomness in the context using the configured context key
		c.Locals(cfg.ContextKey, randomness)

		// Create a map to store custom values for the CSP header
		customValues := make(map[string]string)

		// Set the CSP header
		//
		// Important: Since this CSP header uses a direct digest (using SHA256) without base64 encoding plus immutable. which is idiomatic way.
		// When using base64 encoding, consider storing the base64 encoded in c.Locals first or somewhere (e.g, database). Avoid fetching the value from
		// the header and then putting it in the render or direct in the render, as the format will be different due to sanitization.
		cspValue := cfg.CSPValueGenerator(randomness, customValues)
		c.Set("Content-Security-Policy", cspValue)

		// Continue to the next middleware/route handler
		return c.Next()
	}
}
