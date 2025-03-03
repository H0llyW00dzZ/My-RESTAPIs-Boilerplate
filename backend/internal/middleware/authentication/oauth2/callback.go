// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package oauth2

import (
	"fmt"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// HandleCallback handles the callback request from Google after the user has authenticated.
// It retrieves the authorization code from the query parameters, exchanges it for an access token,
// and then uses the access token to retrieve the user's information from the Google API.
//
// TODO: This still needs improvement and must be combined with Fiber's rate limiter to protect against bots bruteforce attacks.
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
	//
	// Note: This already uses Fiber's session middleware mechanism, so the session is not using OAuth2's built-in mechanisms such as token, JWT, JWS, etc.
	// This approach is better oauth2 custom and considered safer.
	storedState := sess.Get("oauth2_state")

	// Just in case it keeps getting an "Invalid state parameter" error, it's because storedState is of type interface{}/any.
	// Explicitly converting it to a string using .(string) is a better approach because the "state" is a string,
	// and it's effective in protecting against CSRF attacks 🤪.
	if state != storedState.(string) {
		// Useful protect against bots bruteforce attacks.
		sess.Destroy()
		return helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid state parameter")
	}

	token, err := m.config.Exchange(ctx, code)
	if err != nil {
		sess.Destroy()
		return helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	client := m.config.Client(ctx, token)
	// TODO: This still needs improvement because Google has many types of OAuth2 (e.g., for desktop, which has been used to implement OAuth2-CLI before, and for web)
	userInfo, err := m.getUserInfo(c, client)
	if err != nil {
		sess.Destroy()
		return helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
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
