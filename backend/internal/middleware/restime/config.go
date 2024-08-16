// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package restime

import (
	"github.com/gofiber/fiber/v2"
)

// Config is the configuration struct for the ResponseTime middleware.
type Config struct {
	// HeaderName is the name of the header where the response time will be set.
	// Default is "X-Response-Time".
	HeaderName string

	// Next is a function that determines whether the response time measurement should be skipped for a particular request.
	// If Next returns true, the middleware will skip the response time measurement and pass the request to the next middleware/handler.
	// Default is nil, meaning no requests will be skipped.
	Next func(*fiber.Ctx) bool
}

// DefaultConfig is the default configuration for the ResponseTime middleware.
var DefaultConfig = Config{
	HeaderName: "X-Response-Time",
	Next:       nil,
}
