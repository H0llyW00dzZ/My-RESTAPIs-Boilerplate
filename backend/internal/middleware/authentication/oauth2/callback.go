// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package oauth2

import (
	"fmt"
	"h0llyw00dz-template/backend/pkg/gc"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// HandleCallback handles the callback request from Google after the user has authenticated.
// It retrieves the authorization code from the query parameters, exchanges it for an access token,
// and then uses the access token to retrieve the user's information from the Google API.
func (m *Manager) HandleCallback(c *fiber.Ctx) error {
	ctx := c.Context()
	code := c.Query("code")
	state := c.Query("state")

	// Get the session from the store
	sess, err := m.store.Get(c)
	if err != nil {
		return err
	}

	// Verify the state parameter
	storedState := sess.Get("oauth2_state")
	if state != storedState {
		return c.Status(http.StatusBadRequest).SendString("Invalid state parameter")
	}

	token, err := m.config.Exchange(ctx, code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	client := m.config.Client(ctx, token)
	// TODO: This still needs improvement because Google has many types of OAuth2 (e.g., for desktop, which has been used to implement OAuth2-CLI before, and for web)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer resp.Body.Close()

	// Get a buffer from the pool
	buf := gc.BufferPool.Get()

	defer func() {
		buf.Reset()            // Reset the buffer to prevent data leaks
		gc.BufferPool.Put(buf) // Return the buffer to the pool for reuse
	}()

	// Read the response body into the buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// TODO: This still needs improvement.
	var userInfo map[string]any
	// Use the decoder from the Fiber app configuration
	if err := c.App().Config().JSONDecoder(buf.Bytes(), &userInfo); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Access user information
	//
	// TODO: This still needs improvement.
	email := userInfo["email"].(string)
	name := userInfo["name"].(string)

	// Perform further actions with the user information
	//
	// TODO: Remove this later when it is fully improved.
	fmt.Printf("User logged in: Email: %s, Name: %s\n", email, name)

	// TODO: This still needs improvement.
	return c.SendString("User logged in successfully")
}
