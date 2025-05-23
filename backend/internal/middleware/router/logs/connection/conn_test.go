// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package connectionlogger_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	connectionlogger "h0llyw00dz-template/backend/internal/middleware/router/logs/connection"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	// This is the max concurrent requests on my laptop without overclocking. It depends on the processor.
	//
	// Note: If the test fails, it indicates that the processor can't handle up to 5555 concurrent requests, as it depends on the processor.
	// For testing 500K or 1 million concurrent requests, consider increasing the BufferedChannelCount and adjusting the Fiber configuration for concurrency.
	// If it reaches 500K or 1 million, it indicates better processor performance.
	concurrentRequests = 5555
)

// Note: This test might fail on some single-core machines because 5555 goroutines are considered a high number.
// It may also fail on faster processors, such as AMD processors (e.g., RYZEN or EPYC) with a high number of cores (both single and multi-core),
// because they are too fast to catch the goroutines. In such cases, the time.Sleep duration must be set to lower than 100ms for faster processors.
func TestConnectionLoggerMiddleware(t *testing.T) {
	app := fiber.New()

	// Buffer to capture log output.
	//
	// Use [strings.Builder] to reduce memory usage.
	// Allow the processor to handle concurrency efficiently in the wilds.
	var buf strings.Builder

	// Initialize middleware with the app
	middleware := connectionlogger.New()
	httplog := logger.New(logger.Config{
		Output: &buf,
		CustomTags: map[string]logger.LogFunc{
			"testLog": connectionlogger.GetActiveConnections,
		},
		Format: "${testLog}\n",
	})

	// Add middleware to the app
	app.Use(middleware, httplog)

	// Define a simple handler
	app.Get("/", func(c *fiber.Ctx) error {
		// Simulate some processing time, let's say keep-alive concurrently
		//
		// Note: This depends on the processor. If there are more than 1K concurrent requests,
		// it might need to increase the time.Sleep duration or adjust the Fiber configuration.
		time.Sleep(100 * time.Millisecond)
		return c.SendString("Hello, World!")

	})

	// Test the middleware
	t.Run("Check active connections with concurrency", func(t *testing.T) {
		// Note: ingress-nginx might become a bottleneck with 10K ~ 1 million requests at same-time.
		// However, this middleware can handle effectively.
		//
		// Additionally, this implementation isn't perfect yet. There are currently no storage mechanisms,
		// and the BufferedChannelCount is not bound to Fiber's context concurrency.
		// If storage mechanisms are implemented and BufferedChannelCount is aligned with Fiber's context concurrency,
		// it could perfectly manage high concurrency.
		start := make(chan struct{})
		var wg sync.WaitGroup

		// Reset the buffer
		buf.Reset()

		// Launch multiple requests concurrently
		for range concurrentRequests {
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
		if !strings.Contains(logOutput, fmt.Sprintf("%d Active Connections", concurrentRequests)) {
			t.Errorf("Expected log output to contain '%d Active Connections', got '%s'", concurrentRequests, logOutput)
		}
	})
}

// Note: This test might fail on some single-core machines because 5555 goroutines are considered a high number.
// It may also fail on faster processors, such as AMD processors (e.g., RYZEN or EPYC) with a high number of cores (both single and multi-core),
// because they are too fast to catch the goroutines. In such cases, the time.Sleep duration must be set to lower than 100ms for faster processors.
func TestConnectionLoggerMiddlewareImmutable(t *testing.T) {
	app := fiber.New(
		fiber.Config{
			Immutable: true,
		},
	)

	// Buffer to capture log output.
	//
	// Use [strings.Builder] to reduce memory usage.
	// Allow the processor to handle concurrency efficiently in the wilds.
	var buf strings.Builder

	// Initialize middleware with the app
	middleware := connectionlogger.New()
	httplog := logger.New(logger.Config{
		Output: &buf,
		CustomTags: map[string]logger.LogFunc{
			"testLog": connectionlogger.GetActiveConnections,
		},
		Format: "${testLog}\n",
	})

	// Add middleware to the app
	app.Use(middleware, httplog)

	// Define a simple handler
	app.Get("/", func(c *fiber.Ctx) error {
		// Simulate some processing time, let's say keep-alive concurrently
		//
		// Note: This depends on the processor. If there are more than 1K concurrent requests,
		// it might need to increase the time.Sleep duration or adjust the Fiber configuration.
		time.Sleep(100 * time.Millisecond)
		return c.SendString("Hello, World!")

	})

	// Test the middleware
	t.Run("[Immutable] Check active connections with concurrency", func(t *testing.T) {
		// Note: ingress-nginx might become a bottleneck with 10K ~ 1 million requests at same-time.
		// However, this middleware can handle effectively.
		//
		// Additionally, this implementation isn't perfect yet. There are currently no storage mechanisms,
		// and the BufferedChannelCount is not bound to Fiber's context concurrency.
		// If storage mechanisms are implemented and BufferedChannelCount is aligned with Fiber's context concurrency,
		// it could perfectly manage high concurrency.
		start := make(chan struct{})
		var wg sync.WaitGroup

		// Reset the buffer
		buf.Reset()

		// Launch multiple requests concurrently
		for range concurrentRequests {
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
		if !strings.Contains(logOutput, fmt.Sprintf("%d Active Connections", concurrentRequests)) {
			t.Errorf("Expected log output to contain '%d Active Connections', got '%s'", concurrentRequests, logOutput)
		}
	})
}
