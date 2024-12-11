// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package connectionlogger

import (
	"github.com/gofiber/fiber/v2"
)

// Config defines the configuration options for the connection logger middleware.
type Config struct {
	// Next defines a function to skip middleware execution.
	Next func(*fiber.Ctx) bool

	// BufferedChannelCount specifies the size of the buffered channel used
	// to manage connection operations. This determines how many connection
	// operations can be queued before blocking. A larger buffer can help
	// accommodate bursts of requests without blocking, but may consume more memory.
	BufferedChannelCount int
}

// DefaultConfig provides default settings for the connection logger middleware.
var DefaultConfig = Config{
	Next: nil,
	// TODO: There might be another way to get Fiber's concurrency configuration.
	// Currently, it uses the global variable [fiber.DefaultConcurrency], which depends on the concurrency setting (number of goroutines).
	BufferedChannelCount: fiber.DefaultConcurrency,
}
