// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

// GetKeyFunc generates a unique key for each request and returns a function that retrieves the key from the context.
func (k *KeyIdentifier) GetKeyFunc() func(*fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		// Generate a random UUID
		id := utils.UUIDv4()

		// Sign the UUID using ECDSA or HSM
		//
		// Note: When ECDSA or HSM is configured, it becomes a premium UUID that can be used in ASN programming.
		// Also note that this ECDSA or HSM is suitable for workers as well, such as goroutine workers, for example:
		// - Maintaining/securing internal mechanisms (e.g., database, ingress, etc.)
		// On the other hand:
		// - Implementing cryptographic authentication mechanisms for clients instead of using JWT, email, password, username, or other credentials
		// It's not only for TLS/code signing or other mechanisms that only maintain/secure external mechanisms. That's why it's implemented here.
		switch {
		case k.config.PrivateKey != nil && k.config.SignedContextKey != nil:
			signature, err := k.signUUIDWithECDSA(id)
			if err != nil {
				panic(fmt.Errorf("failed to sign UUID: %v", err))
			}
			c.Locals(k.config.SignedContextKey, signature)

			// Test Skipped for HSM
		case k.config.HSM != nil && k.config.SignedContextKey != nil:
			signature, err := k.signUUIDWithHSM(id)
			if err != nil {
				panic(fmt.Errorf("failed to sign UUID: %v", err))
			}
			c.Locals(k.config.SignedContextKey, signature)
		}

		// Set the key in the context
		key := k.config.Prefix + id

		// Return the generated key
		// Note: This won't be affected anyway and is not supported by other Fiber middleware mechanisms that use storage + key generators
		// because the way they are implemented is incorrect (e.g., Custom KeyGenerator does not work properly in the fiber cache middleware).
		return utils.CopyString(key)
	}
}

// GetKey generates a unique key for each request and retrieves it from the context.
func (k *KeyIdentifier) GetKey() string {
	// Generate a random UUID
	//
	// TODO: Do we really need to improve this by using a cryptographic technique similar to how Bitcoin generates addresses?
	id := utils.UUIDv4()

	// Set the key in the context
	key := k.config.Prefix + id

	// Return the generated key
	// Note: This won't be affected anyway and is not supported by other Fiber middleware mechanisms that use storage + key generators
	// because the way they are implemented is incorrect (e.g., Custom KeyGenerator does not work properly in the fiber cache middleware).
	return utils.CopyString(key)
}

// GenerateCacheKey generates a cache key based on the request method, URL path, and query parameters.
//
// It takes the following parameter:
//   - c: The Fiber context.
//
// It returns the generated cache key as a string.
//
// This function generates a cache key by concatenating the request method, URL path, and sorted query parameters.
// It then computes the SHA-256 hash of the concatenated string and returns the hexadecimal representation of the hash.
//
// Example usage:
//
//	cacheKey := k.GenerateCacheKey(c)
//
//	// Create a new KeyIdentifier instance
//	cacheKeyGen := keyidentifier.New(keyidentifier.Config{
//		Prefix: "frontend:",
//	})
//
//	// Use the cache middleware with the custom key generator
//	app.Use(cache.New(cache.Config{
//		KeyGenerator: cacheKeyGen.GenerateCacheKey,
//		Storage: mystorage,
//
//	}))
//
// Note: This is now suitable and secure to use with the Fiber cache middleware because it computes the SHA-256 hash of the key instead of using c.Patch() (Default Fiber).
// For example, "frontend:44658f661a1a27cf94e51bf48947525e5dfcfb6f95050b52800300f2554b7f99_GET_body",
// where 44658f661a1a27cf94e51bf48947525e5dfcfb6f95050b52800300f2554b7f99_GET_body is the actual key to get the value.
// Previously, it was not secure because the key directly used c.Path(), which could leak sensitive information to the public, for example, in Redis/Valkey logs, commander panels, cloud.
// Also note that this only works with the Fiber cache middleware and can enhance speed performance for HTTP/3 in load balancers.
func (k *KeyIdentifier) GenerateCacheKey(c *fiber.Ctx) string {
	// Get the request method
	method := c.Method()

	// Get the URL path
	path := c.Path()

	// Get the sorted query parameters
	queryParams := getSortedQueryParams(c.Request().URI().QueryArgs())

	// Concatenate the method, path, and query parameters
	key := fmt.Sprintf("%s:%s?%s", method, path, queryParams)

	// Compute the SHA-256 hash of the key
	//
	// Note: not possible use k.Config.Digest because this required no error
	digest := sha256.Sum256([]byte(key))

	// Convert the hash to a hexadecimal string
	cacheKey := hex.EncodeToString(digest[:])

	// No need to copy; this is already an immutable, built-in, secure cryptographic hash.
	return k.config.Prefix + cacheKey
}

// getSortedQueryParams returns the sorted query parameters as a string.
//
// It takes the following parameter:
//   - queryParams: The query parameters as a *fasthttp.Args.
//
// It returns the sorted query parameters as a string.
//
// This function sorts the query parameter keys, concatenates the key-value pairs, and returns the resulting string.
// The purpose of sorting the query parameters is to ensure that the order of the parameters does not affect the generated cache key.
func getSortedQueryParams(queryParams *fasthttp.Args) string {
	var params []string

	// Iterate over the query parameters
	queryParams.VisitAll(func(key, value []byte) {
		params = append(params, fmt.Sprintf("%s=%s", string(key), string(value)))
	})

	// Sort the key-value pairs
	slices.Sort(params)

	// Join the sorted key-value pairs with "&"
	return strings.Join(params, "&")
}
