// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package proxytrust

import (
	"github.com/gofiber/fiber/v2"
)

// Config defines the configuration options for the proxy middleware.
type Config struct {
	// Next is a function that determines whether the middleware should skip
	// processing for a particular request. If Next returns true, the middleware
	// will skip its logic and pass the request to the next handler.
	// Default is nil, meaning no requests will be skipped.
	Next func(*fiber.Ctx) bool
}

// DefaultConfig is the default configuration for the proxy middleware.
var DefaultConfig = Config{
	Next: nil,
}
