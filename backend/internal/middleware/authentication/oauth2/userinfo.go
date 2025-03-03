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
	"golang.org/x/oauth2/google"
)

// getUserInfo retrieves the user information from the OAuth2 provider's API.
// It takes an OAuth2 client and a Fiber context as parameters and returns a map of user information.
func (m *Manager) getUserInfo(c *fiber.Ctx, client *http.Client) (map[string]any, error) {
	var userInfoURL string

	// TODO: This still needs improvement because Google has many types of OAuth2 (e.g., for desktop, which has been used to implement OAuth2-CLI before, and for web)
	switch m.config.Endpoint {
	case google.Endpoint:
		userInfoURL = googleUserInfoURL
	default:
		return nil, fmt.Errorf("unsupported provider endpoint")
	}

	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// TODO: This still needs improvement.
	var userInfo map[string]any
	// Use the decoder from the Fiber app configuration
	if err := c.App().Config().JSONDecoder(buf.Bytes(), &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
