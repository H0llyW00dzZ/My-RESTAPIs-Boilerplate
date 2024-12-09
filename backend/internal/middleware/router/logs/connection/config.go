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
}

// DefaultConfig provides default settings for the connection logger middleware.
var DefaultConfig = Config{
	Next: nil,
}
