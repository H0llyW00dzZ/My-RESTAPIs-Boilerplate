// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware_test

import (
	"h0llyw00dz-template/backend/internal/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TODO: Implement another load balancer, as the current implementation of this Balancer should be sufficient
// even for Kubernetes Ingress, or use it as a Kubernetes Ingress (replace a nginx hahaha) for smooth sailing ⛵ ☸
func TestNewProxying(t *testing.T) {
	// Create a test server that will act as the backend server
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert the custom header in the proxied request
		customHeader := r.Header.Get("X-Custom-Header")
		expectedHeader := "custom-value"
		if customHeader != expectedHeader {
			t.Errorf("Expected custom header '%s', got '%s'", expectedHeader, customHeader)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Backend server response"))
	}))
	defer backendServer.Close()

	// Create a Fiber app
	app := fiber.New()

	// Create a proxying middleware with custom options
	proxyingMiddleware := middleware.NewProxying(
		middleware.WithProxyingServers([]string{backendServer.URL}),
		middleware.WithProxyingTimeout(5*time.Second),
		middleware.WithProxyingModifyRequest(func(c *fiber.Ctx) error {
			c.Request().Header.Set("X-Custom-Header", "custom-value")
			return nil
		}),
	)

	// Register the proxying/standalone load balancer middleware
	app.Use(proxyingMiddleware)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Assert the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	// Assert the response body
	expectedBody := "Backend server response"
	if string(body) != expectedBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedBody, string(body))
	}
}
