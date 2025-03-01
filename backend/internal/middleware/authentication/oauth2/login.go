// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package oauth2

import (
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// HandleLogin handles the login request.
// It generates the authorization URL and redirects the user to the Google login page.
func (m *Manager) HandleLogin(c *fiber.Ctx) error {
	// Note: This safe against CSRF Attacks 🤪
	state, err := rand.GenerateFixedUUID()
	if err != nil {
		return err
	}

	// Get the session from the store
	sess, err := m.store.Get(c)
	if err != nil {
		return err
	}

	// Store the state in the session
	sess.Set("oauth2_state", state)
	if err := sess.Save(); err != nil {
		return err
	}

	authURL := m.config.AuthCodeURL(state)
	return c.Redirect(authURL, http.StatusTemporaryRedirect)
}
