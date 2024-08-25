// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CustomNextContentType is a helper function that creates a custom Next function for the fiber middleware.
//
// The returned Next function checks the content type of the response and determines whether to skip
// the middleware based on the provided content types. If the response content type starts with any of the
// specified content types, the Next function returns true, indicating that the middleware should be skipped.
//
// Example usage:
//
//	// Create a custom skipper function to skip the middleware for HTML responses
//	htmlSkipper := CustomNextContentType(fiber.MIMETextHTMLCharsetUTF8)
//
// Parameters:
//   - contentTypes: Variadic string parameters representing the content types to skip the middleware for.
//
// Returns:
//   - A function that takes a [fiber.Ctx] as input and returns a boolean value indicating whether to
//     skip the middleware for the given request based on the response content type.
//
// Note: This function is suitable for cache middleware.
func CustomNextContentType(contentTypes ...string) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		for _, contentType := range contentTypes {
			if strings.HasPrefix(string(c.Response().Header.ContentType()), contentType) {
				return true
			}
		}
		return false
	}
}

// CustomNextPathAvailable is a helper function that creates a custom Next function for the fiber middleware.
//
// The returned Next function checks if the requested path is available in the provided map of paths.
// If the path is found in the map, the Next function returns false, indicating that the middleware should not be skipped.
// Otherwise, it returns true, indicating that the middleware should be skipped.
//
// Example usage:
//
//	// Create a custom skipper function to skip the middleware for available paths
//	pathSkipper := CustomNextPathAvailable(map[string]bool{
//	    "v2/users": true,
//	    "v2/products": true,
//	})
//
// Parameters:
//   - paths: A map of string keys representing the available paths, and bool values indicating their availability.
//
// Returns:
//   - A function that takes a [*fiber.Ctx] as input and returns a boolean value indicating whether to
//     skip the middleware for the given request based on the availability of the requested path.
func CustomNextPathAvailable(paths map[string]bool) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		_, ok := paths[c.Path()]
		return !ok
	}
}

// CustomNextStack is a helper function that creates a custom Next stack function for the fiber middleware
// by combining multiple custom Next functions based on their keys in the provided map.
//
// The returned Next function checks each custom Next function in the map and determines whether
// to skip the middleware based on their combined result. If any of the custom Next functions
// returns true, indicating that the middleware should be skipped, the combined Next function
// returns true. Otherwise, it returns false.
//
// Example usage:
//
//	// Create custom Next functions
//	pathSkipper := CustomNextPathAvailable(map[string]bool{
//	    "/api/v1/users": true,
//	    "/api/v1/products": true,
//	})
//	htmlSkipper := CustomNextContentType(fiber.MIMETextHTMLCharsetUTF8)
//
//	// Combine the custom Next stack functions
//	customNextStack := CustomNextStack(map[string]func(*fiber.Ctx) bool{
//	    "pathSkipper":  pathSkipper,
//	    "htmlSkipper":  htmlSkipper,
//	})
//
// Parameters:
//   - nextFuncs: A map of string keys representing the names of the custom Next stack functions,
//     and func(*fiber.Ctx) bool values representing the custom Next stack functions themselves.
//
// Returns:
//   - A function that takes a [fiber.Ctx] as input and returns a boolean value indicating whether to
//     skip the middleware for the given request based on the combined result of the custom Next stack functions.
func CustomNextStack(nextFuncs map[string]func(*fiber.Ctx) bool) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		for _, nextFunc := range nextFuncs {
			if nextFunc(c) {
				return true
			}
		}
		return false
	}
}

// CustomNextStatusCode is a helper function that creates a custom Next function for the fiber middleware.
//
// The returned Next function checks the HTTP status code of the response and determines whether to skip
// the middleware based on the provided status codes. If the response status code matches any of the
// specified status codes, the Next function returns true, indicating that the middleware should be skipped.
//
// Example usage:
//
//	// Create a custom skipper function to skip the middleware for 404 and 500 status codes
//	redirectSkipper := CustomNextStatusCode(fiber.StatusMovedPermanently)
//
// Parameters:
//   - statusCodes: Variadic int parameters representing the HTTP status codes to skip the middleware for.
//
// Returns:
//   - A function that takes a [fiber.Ctx] as input and returns a boolean value indicating whether to
//     skip the middleware for the given request based on the response status code.
//
// Note: This function is suitable for redirect or error handling middleware or middleware that should be skipped for certain status codes.
func CustomNextStatusCode(statusCodes ...int) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		status := c.Response().StatusCode()
		for _, code := range statusCodes {
			if status == code {
				return true
			}
		}
		return false
	}
}
