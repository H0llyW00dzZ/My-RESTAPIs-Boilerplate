// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package proxytrust

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// New creates a new middleware handler for checking trusted proxies.
// It accepts an optional Config parameter to customize its behavior.
//
// If no configuration is provided, it uses the [DefaultConfig].
//
// The middleware checks if the incoming request is from a trusted proxy using
// the IsProxyTrusted method. If the proxy is not trusted, it returns a
// StatusGatewayTimeout error. Otherwise, it proceeds to the next handler.
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := DefaultConfig

	// Override default config with provided configuration
	if len(config) > 0 {
		if config[0].Next != nil {
			cfg.Next = config[0].Next
		}
		if config[0].StatusCode != DefaultConfig.StatusCode {
			cfg.StatusCode = config[0].StatusCode
		}
	}

	return func(c *fiber.Ctx) error {
		// Check if the trusted proxy check is enabled
		if !c.App().Config().EnableTrustedProxyCheck {
			// Skip if it is not enabled
			log.Warn("[Router] [ProxyTrust]: EnableTrustedProxyCheck is not enabled, skipping...")
			return c.Next()
		}

		// Check if the request should be skipped
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Validate if the proxy is trusted
		if !c.IsProxyTrusted() {
			// Note: Returning a new error is a better approach instead of returning directly,
			// as it allows the error to be handled by the caller somewhere else in the codebase,
			// especially when the codebase grows larger.
			// Additionally, this is preferable to returning a status of 'forbidden'.
			return fiber.NewError(cfg.StatusCode)
		}

		// Proceed to the next middleware/handler
		return c.Next()
	}
}
