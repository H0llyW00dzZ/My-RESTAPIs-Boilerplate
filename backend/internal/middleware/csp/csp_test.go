// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package csp_test

import (
	"encoding/json"
	"fmt"
	"h0llyw00dz-template/backend/internal/middleware/csp"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCSPMiddleware(t *testing.T) {
	app := fiber.New()

	// Test case 1: Default configuration
	app.Use(csp.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hell0 W0rldz",
		})
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.0.1")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expectedCSP := "script-src 'nonce-37d7a80604871e579850a658c7add2ae7557d0c6abcc9b31ecddc4424207eba3'"
	if resp.Header.Get("Content-Security-Policy") != expectedCSP {
		t.Errorf("Unexpected Content-Security-Policy header value: %s", resp.Header.Get("Content-Security-Policy"))
	}
	if resp.Header.Get("Content-Security-Policy") == "" {
		t.Error("Content-Security-Policy header is empty")
	}

	// Test case 2: Custom configuration
	app.Use(csp.New(csp.Config{
		RandomnessGenerator: func(clientIP string) string {
			return "custom_randomness"
		},
		ContextKey: "custom_csp_key",
		CSPValueGenerator: func(randomness string, customValues map[string]string) string {
			customValues["default-src"] = "'self'"
			return fmt.Sprintf("script-src 'nonce-%s' %s", randomness, customValues["default-src"])
		},
	}))

	app.Get("/custom", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"custom_csp_key": c.Locals("custom_csp_key"),
		})
	})

	req = httptest.NewRequest("GET", "/custom", nil)
	req.Header.Set("CF-Connecting-IP", "127.1.1.1")
	req.Header.Set("CF-Ray", "12345")
	req.Header.Set("CF-IPCountry", "US")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expectedCSP = "script-src 'nonce-custom_randomness' 'self'"
	if resp.Header.Get("Content-Security-Policy") != expectedCSP {
		t.Errorf("Unexpected Content-Security-Policy header value: %s", resp.Header.Get("Content-Security-Policy"))
	}

	var responseBody map[string]any
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if responseBody["custom_csp_key"] != "custom_randomness" {
		t.Errorf("Unexpected custom context key value: %v", responseBody["custom_csp_key"])
	}

	// Test case 3: Multiple IP addresses
	app.Use(csp.New())

	app.Get("/multiple-ips", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Multiple IPs",
		})
	})

	req = httptest.NewRequest("GET", "/multiple-ips", nil)
	req.Header.Set("X-Real-IP", "192.168.0.1, 192.168.0.2, 192.168.0.3")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expectedCSP = "script-src 'nonce-37d7a80604871e579850a658c7add2ae7557d0c6abcc9b31ecddc4424207eba3'"
	if resp.Header.Get("Content-Security-Policy") != expectedCSP {
		t.Errorf("Unexpected Content-Security-Policy header value: %s", resp.Header.Get("Content-Security-Policy"))
	}
	if resp.Header.Get("Content-Security-Policy") == "" {
		t.Error("Content-Security-Policy header is empty")
	}

	// Test case 4: Invalid IP address provided
	app.Use(csp.New())

	app.Get("/invalid-ip", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Invalid IP provided",
		})
	})

	req = httptest.NewRequest("GET", "/invalid-ip", nil)
	req.Header.Set("X-Real-IP", "invalid-ip-address")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the CSP header is set correctly or remains empty as expected
	if resp.Header.Get("Content-Security-Policy") == "" {
		t.Error("Content-Security-Policy header is empty as expected for invalid IP")
	} else {
		// Note: This test may yield different results on different machines. However, on my laptop,
		// the default local IP for testing is "19e36255972107d42b8cecb77ef5622e842e8a50778a6ed8dd1ce94732daca9e", which corresponds to 0.0.0.0.
		expectedCSP := "script-src 'nonce-19e36255972107d42b8cecb77ef5622e842e8a50778a6ed8dd1ce94732daca9e'"
		if resp.Header.Get("Content-Security-Policy") != expectedCSP {
			t.Errorf("Unexpected Content-Security-Policy header value: %s", resp.Header.Get("Content-Security-Policy"))
		}
	}
}
