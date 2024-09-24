// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package webauthn

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// BeginLogin initiates the WebAuthn login process.
// It takes a [webauthn.User] and optional [webauthn.LoginOption] as input.
// It returns a [protocol.CredentialAssertion], [webauthn.SessionData], and an error.
// The returned session data must be stored securely and passed to the [FinishLogin] function.
func BeginLogin(user webauthn.User, opts ...webauthn.LoginOption) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	return web.BeginLogin(user, opts...)
}

// FinishLogin completes the WebAuthn login process.
// It takes a [webauthn.User], the session data obtained from [BeginLogin], and the parsed credential assertion response.
// It returns a [webauthn.Credential] and an error.
func FinishLogin(user webauthn.User, sessionData *webauthn.SessionData, response *protocol.ParsedCredentialAssertionData) (*webauthn.Credential, error) {
	return web.ValidateLogin(user, *sessionData, response)
}
