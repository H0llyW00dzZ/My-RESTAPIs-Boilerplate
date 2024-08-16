// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package eth provides a middleware for configuring an Ethereum client in a Fiber application.
//
// The middleware creates an Ethereum client using the provided configuration and stores it in the Fiber context
// for subsequent use in route handlers. It also handles errors that may occur during the client creation process.
//
// Configuration:
//
// The [Config] struct allows configuring the Ethereum client:
//
//	type Config struct {
//		URL         string
//		ContextKey  any
//		ErrorHandler func(c *fiber.Ctx, err error) error
//		Next         func(*fiber.Ctx) bool
//	}
//
// Configuration:
//   - URL: The URL of the Ethereum network to connect to (e.g., "https://eth.btz.pm").
//   - ContextKey: The key used to store the Ethereum client in the Fiber context. This key must be specified when creating the Config struct.
//   - ErrorHandler: A custom error handler function to handle errors that occur during client creation.
//     If not provided, it defaults to using htmx.NewStaticHandleVersionedAPIError.
//   - Next: A function that determines whether to skip the middleware and proceed to the next middleware or route handler.
//     It takes a Fiber context as input and returns a boolean value.
//     If true, the middleware is skipped, and the next middleware or route handler is executed.
//     If false, the middleware continues its execution.
//
// Retrieving the Ethereum Client:
//
// In route handlers, the Ethereum client can be retrieved from the Fiber context using the specified ContextKey:
//
//	client := c.Locals(config.ContextKey).(*ethclient.Client)
//
// Make sure to type-assert the retrieved value to [*ethclient.Client].
//
// Error Handling:
//
// If an error occurs during the Ethereum client creation, the middleware will call the specified ErrorHandler function.
// A custom error handler can be provided in the [Config] to handle errors according to specific requirements.
// If no custom error handler is provided, it defaults to using [htmx.NewStaticHandleVersionedAPIError].
//
// Cleaning Up:
//
// The middleware automatically closes the Ethereum client when the request is finished using defer [client.Close].
// This ensures proper cleanup of resources.
//
// Note: Make sure to import the necessary dependencies and update the import paths based on the project structure.
package eth
