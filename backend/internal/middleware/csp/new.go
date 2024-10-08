// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package csp

import (
	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
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

		// Initialize a slice to hold client IPs with additional information
		var clientIPsWithInfo []string

		// Check if Cloudflare is detected
		// this already unique
		cloudflareRayIDHandle := utils.CopyString(c.Get(log.CloudflareRayIDHeader))
		// Get country code from request header
		// this already unique
		countryCodeHandle := utils.CopyString(c.Get(log.CloudflareIPCountryHeader))

		for _, ip := range clientIPs {
			if cloudflareRayIDHandle != "" {
				ip += " - Cloudflare detected - Ray ID: " + cloudflareRayIDHandle
			}

			if countryCodeHandle != "" {
				ip += ", Country: " + countryCodeHandle
			}

			clientIPsWithInfo = append(clientIPsWithInfo, ip)
		}

		clientIPs = clientIPsWithInfo

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
		c.Set(fiber.HeaderContentSecurityPolicy, cspValue)

		// Continue to next middleware
		return c.Next()
	}
}
