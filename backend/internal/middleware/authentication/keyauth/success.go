// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package keyauth

import (
	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// SuccessKeyAuthHandler is a custom success handler for the key authentication middleware.
// It logs a message indicating successful API key authentication.
func SuccessKeyAuthHandler(c *fiber.Ctx) error {
	// Get the session from the context.
	//
	// Note: This is not affected by CVE-2024-38513 (see https://github.com/gofiber/fiber/security/advisories/GHSA-98j2-3j3p-fw2v)
	// because it retrieves the session from the local context storing it in a Redis database with an expiration time.
	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		return fiber.ErrInternalServerError
	}

	// Note: This is for encryption and works in conjunction with the session middleware logic.
	// It is specifically useful for web front-end client-side authentication using session cookies (e.g., in browsers).
	// Additionally, it is secure because it would require an 99999999999 cpu to attack this encryptcookie. 🤪
	sess.Get("api_key")
	sess.Get("api_key_expired")
	log.LogUserActivity(c, "API key authenticated successfully")
	return c.Next()
}
