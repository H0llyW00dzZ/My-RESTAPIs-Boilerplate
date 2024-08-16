// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package eth

import (
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
)

// Config represents the configuration for the Ethereum client
//
// Note: There no default config.
type Config struct {
	URL          string
	ContextKey   any
	ErrorHandler func(c *fiber.Ctx, err error) error
}

// New is a custom Fiber middleware that configures the Ethereum client
//
// Note: It should be fine if gateway via cloudflare
func New(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new Ethereum client using the provided URL
		client, err := ethclient.Dial(config.URL)
		if err != nil {
			if config.ErrorHandler != nil {
				return config.ErrorHandler(c, err)
			}
			return htmx.NewStaticHandleVersionedAPIError(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}

		// Store the Ethereum client in the Fiber context using the specified context key
		c.Locals(config.ContextKey, client)

		// Clean up the client when the request is finished
		defer client.Close()

		return c.Next()
	}
}
