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
// even for Kubernetes Ingress, or use it as a Kubernetes Ingress (replace an Nginx, hahaha) for smooth sailing ⛵ ☸
// Also note that this is just a test. Due to the hybrid technique and atomics used in this load balancer mechanism,
// it can be used to improve region (region-based routing) depending on the client as well. For example, if a client is from Indonesia,
// it can forward the request to an Indonesian server or another server. On the other hand, it can also improve
// security mechanisms such as firewalls, authentication mechanisms, etc.
func TestNewProxying(t *testing.T) {
	// Create a test server that will act as the backend server
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert the custom header in the proxied request
		customHeader := r.Header.Get("X-Custom-Header")
		expectedHeader := "ahoy"
		if customHeader != expectedHeader {
			t.Errorf("Expected custom header '%s', got '%s'", expectedHeader, customHeader)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("yo ho ho ⛵ ☸"))
	}))
	defer backendServer.Close()

	// Create a Fiber app
	app := fiber.New()

	// Create a proxying middleware with custom options
	proxyingMiddleware := middleware.NewProxying(
		middleware.WithProxyingServers([]string{backendServer.URL}),
		middleware.WithProxyingTimeout(5*time.Second),
		middleware.WithProxyingModifyRequest(func(c *fiber.Ctx) error {
			c.Request().Header.Set("X-Custom-Header", "ahoy")
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
	expectedBody := "yo ho ho ⛵ ☸"
	if string(body) != expectedBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedBody, string(body))
	}
}
