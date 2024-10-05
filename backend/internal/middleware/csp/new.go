// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package csp

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"unique"

	"github.com/gofiber/fiber/v2"
)

// New creates a new instance of the CSP middleware with the provided configuration.
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		// Test Skipped For This when Next True
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get client IP address
		clientIPs := getClientIP(c, cfg.IPHeader)

		// Check if Cloudflare is detected
		cloudflareRayIDHandle := unique.Make(c.Get(log.CloudflareRayIDHeader))
		if cloudflareRayIDHandle.Value() != "" {
			for i, clientIP := range clientIPs {
				clientIPs[i] = clientIP + " - Cloudflare detected - Ray ID: " + cloudflareRayIDHandle.Value()
			}
		}

		// Get country code from request header
		countryCodeHandle := unique.Make(c.Get(log.CloudflareIPCountryHeader))
		if countryCodeHandle.Value() != "" {
			for i, clientIP := range clientIPs {
				clientIPs[i] = clientIP + ", Country: " + countryCodeHandle.Value()
			}
		}

		// Generate the randomness using the configured generator
		var randomness string
		if len(clientIPs) > 0 {
			randomness = cfg.RandomnessGenerator(clientIPs[0])
		}

		// Store randomness in context
		c.Locals(cfg.ContextKey, randomness)

		// Create a map to store custom values for the CSP header
		customValues := make(map[string]string)

		// Set the CSP header
		//
		// Important: Since this CSP header uses a direct digest (using SHA256) without base64 encoding plus immutable. which is idiomatic way.
		// When using base64 encoding, consider storing the base64 encoded in c.Locals first or somewhere (e.g, database). Avoid fetching the value from
		// the header and then putting it in the render or direct in the render, as the format will be different due to sanitization.
		cspValue := cfg.CSPValueGenerator(randomness, customValues)

		// Set CSP header
		c.Set("Content-Security-Policy", cspValue)

		// Continue to next middleware
		return c.Next()
	}
}
