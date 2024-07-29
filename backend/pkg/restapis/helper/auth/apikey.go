// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package auth

import (
	"h0llyw00dz-template/backend/pkg/restapis/helper"

	"github.com/gofiber/fiber/v2"
)

// AuthenticateAPIKey retrieves the authenticated API key from the request context.
//
// Example Usage:
//
//	// Authenticate the API key
//	apiKey, err := auth.AuthenticateAPIKey(c)
//	if err != nil {
//		return err
//	}
//
// Note: This enhances security by combining Fiber's Key Auth and Session Middleware Logic (See backend/internal/middleware/authentication/keyauth).
func AuthenticateAPIKey(c *fiber.Ctx) (string, error) {
	apiKey, ok := c.Locals(apikey).(string)
	if !ok {
		return "", helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid API key")
	}
	return apiKey, nil
}
