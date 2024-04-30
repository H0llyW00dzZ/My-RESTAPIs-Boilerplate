// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package keyauth

import (
	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
)

// SuccessKeyAuthHandler is a custom success handler for the key authentication middleware.
// It logs a message indicating successful API key authentication.
func SuccessKeyAuthHandler(c *fiber.Ctx) error {
	log.LogUserActivity(c, "API key authenticated successfully")
	return c.Next()
}
