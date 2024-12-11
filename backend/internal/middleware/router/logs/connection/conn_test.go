// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package connectionlogger_test

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	connectionlogger "h0llyw00dz-template/backend/internal/middleware/router/logs/connection"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func TestConnectionLoggerMiddleware(t *testing.T) {
	app := fiber.New()

	// Buffer to capture log output
	var buf bytes.Buffer

	// Initialize middleware with the app
	middleware := connectionlogger.New()
	httplog := logger.New(logger.Config{
		Output: &buf,
		CustomTags: map[string]logger.LogFunc{
			"testLog": connectionlogger.GetActiveConnections,
		},
		Format: "${testLog}",
	})

	// Add middleware to the app
	app.Use(middleware, httplog)

	// Define a simple handler
	app.Get("/", func(c *fiber.Ctx) error {
		// Simulate some processing time, let's say keep-alive concurrently
		time.Sleep(100 * time.Millisecond)
		return c.SendString("Hello, World!")
	})

	// Test the middleware
	t.Run("Check active connections with concurrency", func(t *testing.T) {
		// Note: ingress-nginx might become a bottleneck with 10K ~ 1 million requests at same-time.
		// However, this middleware can handle effectively.
		concurrentRequests := 1000
		start := make(chan struct{})
		var wg sync.WaitGroup

		// Reset the buffer
		buf.Reset()

		// Launch multiple requests concurrently
		for i := 0; i < concurrentRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start // Wait for the start signal
				req := httptest.NewRequest("GET", "/", nil)
				app.Test(req, -1)
			}()
		}

		// Start all requests at the same time
		close(start)

		// Wait for all requests to complete
		wg.Wait()

		// Capture the logger output
		logOutput := buf.String()

		// Check for expected log output
		if !strings.Contains(logOutput, "1000 Active Connections") {
			t.Errorf("Expected log output to contain '3 Active Connections', got '%s'", logOutput)
		}
	})
}
