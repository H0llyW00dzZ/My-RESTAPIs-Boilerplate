// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
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

// CustomNextHeader is a helper function that creates a custom Next function for the fiber middleware.
//
// The returned Next function checks the presence of a specific header in the request and determines
// whether to skip the middleware based on the provided header keys. If any of the specified header keys
// are found in the request headers, the Next function returns true, indicating that the middleware should be skipped.
//
// Example usage:
//
//	// Create a custom skipper function to skip the middleware if the "X-Cache" or "X-Proxy" headers are present
//	cacheSkipper := CustomNextHeader("X-Cache", "X-Proxy")
//
// Parameters:
//   - headerKeys: Variadic string parameters representing the header keys to check for in the request headers.
//
// Returns:
//   - A function that takes a [fiber.Ctx] as input and returns a boolean value indicating whether to
//     skip the middleware for the given request based on the presence of any of the specified header keys.
func CustomNextHeader(headerKeys ...string) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		for _, headerKey := range headerKeys {
			// TODO: Remove utils.CopyString ? as it is only used for references and not for storing values
			if utils.CopyString(c.Get(headerKey)) != "" {
				return true
			}
		}
		return false
	}
}

// CustomNextHostName is a helper function that creates a custom Next function for the fiber middleware.
//
// The returned Next function checks the hostname of the request and determines whether to skip
// the middleware based on the provided hostnames. If the request hostname matches any of the
// specified hostnames, the Next function returns true, indicating that the middleware should be skipped.
//
// Example usage:
//
//	// Create a custom skipper function to skip the middleware for "example.com" and "www.example.com"
//	hostSkipper := CustomNextHostName("example.com", "www.example.com")
//
// Parameters:
//   - hostnames: Variadic string parameters representing the hostnames to skip the middleware for.
//
// Returns:
//   - A function that takes a [fiber.Ctx] as input and returns a boolean value indicating whether to
//     skip the middleware for the given request based on the request hostname.
//
// Note: This function is suitable for middleware that should be skipped for specific hostnames or domains.
func CustomNextHostName(hostnames ...string) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		// TODO: Remove utils.CopyString ? as it is only used for references and not for storing values
		currentHostname := utils.CopyString(c.Hostname())
		for _, hostname := range hostnames {
			if currentHostname == hostname {
				return true
			}
		}
		return false
	}
}

// CustomKeyGenerator generates a custom cache key based on the request and logs the visitor activity.
func CustomKeyGenerator(c *fiber.Ctx) string {
	// Get client's IP and User-Agent
	clientIP := c.IP()
	userAgent := c.Get(fiber.HeaderUserAgent)

	// Create a string to hash
	toHash := fmt.Sprintf("%s-%s", clientIP, userAgent)

	// Create a fnv hash and write our string to it
	signature := hashForSignature(toHash)

	// Generate a UUID based on the hash
	signatureUUID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(signature))

	// Log visitor activity with the signature for the frontend
	// Note: Rename "generated" to "initiated" because This cache is used for fiber storage operations (e.g., get, set, delete, reset, close).
	log.LogVisitorf("Frontend cache initiated for visitor activity: IP [%s], User-Agent [%s], Signature [%s], UUID [%s]", clientIP, userAgent, signature, signatureUUID.String())

	// Generate a custom cache key with the hashed signature and UUID for fiber storage operations (e.g., get, set, delete, reset, close).
	return fmt.Sprintf(utils.CopyString("cache_front_end:%s:%s:%s"), signature, signatureUUID.String(), c.Path())
}

// CustomCacheSkipper is a function that determines whether to skip caching for a given request path.
// It returns true if the request path starts with any of the specified prefixes.
func CustomCacheSkipper(prefixes ...string) func(*fiber.Ctx) bool {
	// Note: This safely integrates with the context (e.g., fiber ctx or std library ctx).
	return func(c *fiber.Ctx) bool {
		for _, prefix := range prefixes {
			if strings.HasPrefix(c.Path(), prefix) {
				return true
			}
		}
		return false
	}
}
