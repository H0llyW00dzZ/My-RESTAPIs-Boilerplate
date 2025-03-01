// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package oauth2

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	providerGoogle = "Google"
)

// Config represents the configuration for the OAuth2 Manager.
type Config struct {
	Provider      string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	SessionConfig session.Config
}

// Manager represents an OAuth2 manager that handles the OAuth2 flow.
// It contains the OAuth2 configuration required for authentication.
type Manager struct {
	config *oauth2.Config
	store  *session.Store
}

// New creates a new instance of the OAuth2 Manager.
// It takes a Config struct as a parameter and returns a pointer to the Manager.
func New(cfg Config) *Manager {
	var config *oauth2.Config

	switch cfg.Provider {
	// TODO: This still needs improvement because Google has many types of OAuth2 (e.g., for desktop, which has been used to implement OAuth2-CLI before, and for web)
	case providerGoogle:
		config = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	default:
		panic(fmt.Sprintf("unsupported provider: %s", cfg.Provider))
	}

	store := session.New(cfg.SessionConfig)

	return &Manager{
		config: config,
		store:  store,
	}
}
