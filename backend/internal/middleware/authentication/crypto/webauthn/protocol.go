// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package webauthn

import "github.com/go-webauthn/webauthn/webauthn"

var (
	// web holds the instance of the WebAuthn client.
	web *webauthn.WebAuthn
)

// Init initializes the WebAuthn client with the provided configuration.
// It returns an error if the initialization fails.
func Init(config *webauthn.Config) error {
	var err error
	web, err = webauthn.New(config)
	return err
}
