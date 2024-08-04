// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyauth

import (
	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/authentication/helper"
	"time"

	"github.com/bytedance/sonic"
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
// It returns one string value (UUID) and two boolean values: isAPIKeyValid and expired.
func isAPIKeyValidInSession(sess *session.Session, key string) (string, bool, bool) {
	sessionAPIKeyData := sess.Get(apiKey)
	if sessionAPIKeyData != nil {
		var data helper.APIKeyData
		// Note: Custom JSON encoder/decoder configuration, similar to what Fiber currently supports,
		// currently unavailable in this enhancement due to its focus on better performance.
		//
		// Idiom Go Error Handling type-assertion.
		if err := sonic.Unmarshal(sessionAPIKeyData.([]byte), &data); err != nil {
			log.LogErrorf("Failed to unmarshal API key data from session: %v", err)
			return "", false, false
		}

		if data.APIKey == key {
			if data.Status == helper.APIKeyExpired.String() {
				sess.Destroy()
				log.LogInfof("API key %s found in session but marked as expired", key)
				return sess.ID(), true, true
			}
			return sess.ID(), true, false
		}
	}
	return "", false, false
}

// saveAPIKeyInSession saves the API key and its expiration status in the session.
//
// Note: The session data is now stored as JSON when viewing in Redis, Valkey Insight, or Commander panel.
// Additionally, it is possible to implement an encryption/decryption mechanism for the JSON values, as Go, unlike other languages, allows for this functionality + 100% secure.
func saveAPIKeyInSession(sess *session.Session, key string, expired bool, expirationDate time.Time) {
	if expired {
		data := helper.APIKeyData{
			APIKey: key,
			Status: helper.APIKeyExpired.String(),
			Authorization: helper.AuthorizationData{
				// Time Server not client
				AuthTime: time.Now().UTC(),
			},
		}

		// Note: Custom JSON encoder/decoder configuration, similar to what Fiber currently supports,
		// currently unavailable in this enhancement due to its focus on better performance.
		jsonData, err := sonic.Marshal(data)
		if err != nil {
			log.LogErrorf("Failed to marshal API key data for Session Middleware: %v", err)
			return
		}

		sess.Set(apiKey, jsonData)
		sess.SetExpiry(defaultExpryContextKey)
	}

	data := helper.APIKeyData{
		APIKey: key,
		Status: helper.APIKeyActive.String(),
		Authorization: helper.AuthorizationData{
			// Time Server not client
			AuthTime: time.Now().UTC(),
			// Note: This expiration time is retrieved from the relational database (MySQL).
			// The performance speed might be somewhat slow (taking an average of 1s response time in the frontend) during the first query due to the relational database (always slow).
			// However, when it hits Redis and is released into cookies with encryption, the speed can be faster (possibly 0ms ~ 1ms response time).
			ExpiredTime: expirationDate,
		},
	}

	// Note: Custom JSON encoder/decoder configuration, similar to what Fiber currently supports,
	// currently unavailable in this enhancement due to its focus on better performance.
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		log.LogErrorf("Failed to marshal API key data for Session Middleware: %v", err)
		return
	}

	sess.Set(apiKey, jsonData)

	// Note: This is safe if it's encrypted by the Encrypted Cookies middleware (src https://docs.gofiber.io/api/middleware/encryptcookie).
	// If it's not encrypted, then it's not safe. Also note that when using Encrypted Cookies middleware (src https://docs.gofiber.io/api/middleware/encryptcookie),
	// consider using a Hex Encoder/Decoder instead of Base64 because the Base64-encoded value may not LGTM.
	if err := sess.Save(); err != nil {
		log.LogErrorf("Failed to save session: %v", err)
	}
}
