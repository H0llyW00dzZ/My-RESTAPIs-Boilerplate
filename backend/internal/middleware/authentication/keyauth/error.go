// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyauth

import (
	"errors"
	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/pkg/restapis/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

// ErrorKeyAuthHandler is a custom error handler for the key authentication middleware.
// It handles different types of authentication errors and sends appropriate error responses.
func ErrorKeyAuthHandler(c *fiber.Ctx, err error) error {
	// Header KeyAuth + Session Middleware Logic Request ID for enhancement
	//
	// Note: This is a different request ID. Even though the header is the same ("x-request-id"),
	// it can be connected to the key auth + session middleware logic.
	c.Locals(keyAuthRequestID)
	switch {
	case errors.Is(err, keyauth.ErrMissingOrMalformedAPIKey):
		// Log the authentication attempt.
		log.LogUserActivity(c, "Attempted Authentication")
		log.LogUserActivity(c, "Missing or Malformed API Key")
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Authentication required")
	case errors.Is(err, database.ErrInvalidAPIKey):
		log.LogUserActivity(c, "Invalid API key")
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid API key")
	case errors.Is(err, database.ErrExpiredAPIKey):
		log.LogUserActivity(c, "API Key Expired")
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "API Key Expired")
	default:
		log.LogErrorf("Unexpected error during API key validation: %v", err)
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}
}
