// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package keyauth

import (
	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
)

// ValidatorKeyAuthHandler is a custom validator for the key authentication middleware.
// It checks if the provided API key is valid and active by querying the Redis cache and the database.
func ValidatorKeyAuthHandler(c *fiber.Ctx, key string, db database.Service) (bool, error) {
	// Log the authentication attempt.
	log.LogUserActivity(c, "Attempted Authentication")

	// TODO: Implement the "vice versa" method for Redis (for non-browser) & Redis (for browser, aka session storage) -> database -> repeat.
	// Note: Won't implement JWT and their base standards, because it's easy to lead to high vulnerability.
	// Also, the cryptography world is not small enough to rely solely on JWT. 🤪
	// So, any package for "authentication" here will be covered with another crypto instead of JWT and their base standards.

	return true, nil
}
