// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package csp_test

import (
	"encoding/json"
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
			"csp_random": c.Locals("csp_random"),
		})
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.0.1")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Header.Get("Content-Security-Policy") != "script-src 'nonce-37d7a80604871e579850a658c7add2ae7557d0c6abcc9b31ecddc4424207eba3'" {
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
	if resp.Header.Get("Content-Security-Policy") != "script-src 'nonce-custom_randomness'" {
		t.Errorf("Unexpected Content-Security-Policy header value: %s", resp.Header.Get("Content-Security-Policy"))
	}

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if responseBody["custom_csp_key"] != "custom_randomness" {
		t.Errorf("Unexpected custom context key value: %v", responseBody["custom_csp_key"])
	}
}
