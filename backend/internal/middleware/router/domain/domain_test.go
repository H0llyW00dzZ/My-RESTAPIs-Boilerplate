// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package domain_test

import (
	"h0llyw00dz-template/backend/internal/middleware/router/domain"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDomainRouter(t *testing.T) {
	// Create a new Fiber instance for the main domain
	mainApp := fiber.New()
	mainApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Main Domain")
	})

	// Create a new Fiber instance for the API subdomain
	apiApp := fiber.New()
	apiApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Domain")
	})

	// Configure the domain router
	config := domain.Config{
		Hosts: map[string]*fiber.App{
			// Note: It is optional to include "www.example.com" in Hosts for production,
			// because MainDomain will link "www." to example.com (mainApp).
			"example.com":     mainApp,
			"api.example.com": apiApp,
		},
		MainDomain: "example.com",
	}

	// Create a new Fiber instance and apply the domain router middleware
	app := fiber.New()
	app.Use(domain.New(config))

	// Test cases
	tests := []struct {
		hostname   string
		expected   string
		statusCode int
	}{
		{"example.com", "Main Domain", fiber.StatusOK},
		// Note: This already supports non-case-sensitive hostnames (e.g., wWw.example.com)
		{"www.example.com", "Main Domain", fiber.StatusOK},
		{"WWW.example.com", "Main Domain", fiber.StatusOK},
		{"wWW.example.com", "Main Domain", fiber.StatusOK},
		{"WwW.example.com", "Main Domain", fiber.StatusOK},
		{"WWw.example.com", "Main Domain", fiber.StatusOK},
		{"api.example.com", "API Domain", fiber.StatusOK},
		{"unknown.example.com", "Service Unavailable", fiber.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = tt.hostname

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(body), tt.expected)
		})
	}
}
