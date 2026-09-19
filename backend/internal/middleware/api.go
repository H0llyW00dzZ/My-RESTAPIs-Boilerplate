// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// APIRoute represents a single API route, containing the path, HTTP method,
// handler function, and optional rate limiters.
type APIRoute struct {
	Path                      string
	Method                    string
	Handler                   fiber.Handler
	RateLimiter               []fiber.Handler
	KeyAuth                   fiber.Handler
	RequestID                 fiber.Handler
	EncryptedCookieMiddleware fiber.Handler
}

// APIGroup represents a group of API routes under a common prefix.
// It also allows for a group-wide rate limiter.
type APIGroup struct {
	Prefix                    string
	Routes                    []APIRoute
	RateLimiter               []fiber.Handler
	KeyAuth                   fiber.Handler
	RequestID                 fiber.Handler
	EncryptedCookieMiddleware fiber.Handler
}
