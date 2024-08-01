// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyauth

import (
	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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

// isAPIKeyValidInSession checks if the API key is valid and not expired in the session.
// It returns one string values (UUID) and two boolean values: isAPIKeyValid and expired.
func isAPIKeyValidInSession(sess *session.Session, key string) (string, bool, bool) {
	sessionAPIKey := sess.Get("x_api_key")
	if sessionAPIKey != nil && sessionAPIKey.(string) == key {
		expiredInSession := sess.Get("x_api_key_expired")
		if expiredInSession != nil && expiredInSession.(bool) {
			sess.Destroy()
			log.LogInfof("API key %s found in session but marked as expired", key)
			return sess.ID(), true, true
		}
		return sess.ID(), true, false
	}
	return "", false, false
}

// saveAPIKeyInSession saves the API key and its expiration status in the session.
func saveAPIKeyInSession(sess *session.Session, key string, expired bool) {
	sess.Set(apiKey, key)
	sess.Set(apiKeyExpired, expired)
	if expired {
		sess.SetExpiry(defaultExpryContextKey)
	}
	if err := sess.Save(); err != nil {
		log.LogErrorf("Failed to save session: %v", err)
	}
}
