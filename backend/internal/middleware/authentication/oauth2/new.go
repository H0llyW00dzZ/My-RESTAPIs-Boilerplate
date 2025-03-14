// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package oauth2

import (
	"fmt"
	"h0llyw00dz-template/backend/internal/database"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// ProviderGoogle represents the identifier for the Google OAuth2 provider.
	// It is used to specify the provider when configuring the OAuth2 Manager.
	ProviderGoogle = "Google"
)

// googleScopes defines the scopes required for Google OAuth2 authentication.
var googleScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

// Config represents the configuration for the OAuth2 Manager.
type Config struct {
	Provider      string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	SessionConfig session.Config
	// Note: The DB field cannot be nil.
	DB database.ServiceAuth
}

// Manager represents an OAuth2 manager that handles the OAuth2 flow.
// It contains the OAuth2 configuration required for authentication.
type Manager struct {
	config *oauth2.Config
	store  *session.Store
	db     database.ServiceAuth
}

// New creates a new instance of the OAuth2 Manager.
// It takes a Config struct as a parameter and returns a pointer to the Manager.
//
// Example Usage:
//
//	// Create an instance of the database service
//	dbService := database.New()
//
//	// Create the OAuth2 configuration
//	cfg := oauth2.Config{
//		Provider:      oauth2.ProviderGoogle,
//		ClientID:      "your-client-id",
//		ClientSecret:  "your-client-secret",
//		RedirectURL:   "your-redirect-url",
//		SessionConfig: sessionConfig,
//		DB:            dbService.AuthUser(),
//	}
//
//	// Create an instance of the OAuth2 manager
//	manager := oauth2.New(cfg)
func New(cfg Config) *Manager {
	var config *oauth2.Config

	switch cfg.Provider {
	// TODO: This still needs improvement because Google has many types of OAuth2 (e.g., for desktop, which has been used to implement OAuth2-CLI before, and for web)
	case ProviderGoogle:
		config = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       googleScopes,
			Endpoint:     google.Endpoint,
		}
	default:
		panic(fmt.Sprintf("unsupported provider: %s", cfg.Provider))
	}

	store := session.New(cfg.SessionConfig)

	return &Manager{
		config: config,
		store:  store,
		// TODO: The DB field will be used to verify the user in the database after the token is exchanged, as mentioned earlier in callback.go.
		db: cfg.DB,
	}
}
