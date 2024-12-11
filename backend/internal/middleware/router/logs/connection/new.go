// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package connectionlogger

import (
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler that logs the current number of active connections.
//
// Note: This [connectionlogger] should be placed at the root router first, then put a logger middleware as documented here: https://docs.gofiber.io/api/middleware/logger
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := DefaultConfig

	// Override default config with provided configuration
	if len(config) > 0 {
		if config[0].Next != nil {
			cfg.Next = config[0].Next
		}
		if config[0].BufferedChannelCount != DefaultConfig.BufferedChannelCount {
			cfg.BufferedChannelCount = config[0].BufferedChannelCount
		}
	}

	// Initialize the channel and start the goroutine once
	//
	// Note: This implementation works well on AMD EPYC™ processors. Performance on other processors may vary.
	initTrackActiveConnections.Do(func() {
		connChan = make(chan bool, cfg.BufferedChannelCount)
		go func() {
			for increment := range connChan {
				if increment {
					atomic.AddInt64(&activeConnections, 1)
				} else {
					atomic.AddInt64(&activeConnections, -1)
				}
			}
		}()
	})

	return func(c *fiber.Ctx) error {
		// Check if the request should be skipped
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Increment the active connection count
		//
		// Note: This is safe for concurrent use. However, using a mutex can decrease performance, so it's not recommended (too bad using mutex).
		//
		// Additionally, if issues arise in a Kubernetes environment, they might be due to ingress configurations (e.g., some ingress-nginx configuration causing slowness)
		// leading to inefficiencies or resource constraints. Consider using the Vertical Pod Autoscaler (VPA) if necessary.
		// Also Ensure CoreDNS is adequately scaled, which may require the Horizontal Pod Autoscaler (HPA) for optimal performance (e.g., reduce latency).
		connChan <- true
		defer func() {
			// Decrement the active connection count when the request is done
			connChan <- false
		}()

		// Process the request
		return c.Next()
	}
}
