// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package proxytrust_test

import (
	"h0llyw00dz-template/backend/internal/middleware/router/proxytrust"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestMiddleware checks the behavior of the proxy middleware.
func TestMiddleware(t *testing.T) {
	// Define test cases for different trusted proxies
	tests := []struct {
		name           string
		trustedProxies []string
		expectedStatus int
	}{
		{
			name:           "Test Fiber App with 0.0.0.0 as trusted proxy",
			trustedProxies: []string{"0.0.0.0"},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Test Fiber App with 10.0.0.2 as untrusted proxy",
			trustedProxies: []string{"10.0.0.2"},
			expectedStatus: fiber.StatusGatewayTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Fiber app with specific trusted proxies
			app := fiber.New(fiber.Config{
				AppName:                 tt.name,
				EnableIPValidation:      true,
				TrustedProxies:          tt.trustedProxies,
				EnableTrustedProxyCheck: true,
			})

			// Use the proxy middleware with default config
			app.Use(proxytrust.New())

			// Define a simple route
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("Hello, World!")
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Test failed with error: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, but got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// TestMiddlewareWithNext checks the behavior of the proxy middleware with the Next config.
func TestMiddlewareWithNext(t *testing.T) {
	app := fiber.New(fiber.Config{
		TrustedProxies:          []string{"10.0.0.2"},
		EnableTrustedProxyCheck: true,
	})

	// Use the proxy middleware with a custom Next function
	app.Use(proxytrust.New(proxytrust.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip the middleware for this test case
			return true
		},
	}))

	// Define a simple route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.1") // Untrusted IP

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Test failed with error: %v", err)
	}

	// Expect 200 OK because the Next function skips the middleware
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, but got %d", fiber.StatusOK, resp.StatusCode)
	}
}
