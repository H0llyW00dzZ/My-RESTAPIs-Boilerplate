// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package webauthn

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// BeginRegistration initiates the WebAuthn registration process.
// It takes a webauthn.User and optional webauthn.RegistrationOption as input.
// It returns a [protocol.CredentialCreation], [webauthn.SessionData], and an error.
// The returned session data must be stored securely and passed to the [FinishRegistration] function.
func BeginRegistration(user webauthn.User, opts ...webauthn.RegistrationOption) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	return web.BeginRegistration(user, opts...)
}

// FinishRegistration completes the WebAuthn registration process.
// It takes a [webauthn.User], the session data obtained from [BeginRegistration], and the parsed credential creation response.
// It returns a [webauthn.Credential] and an error.
func FinishRegistration(user webauthn.User, sessionData *webauthn.SessionData, response *protocol.ParsedCredentialCreationData) (*webauthn.Credential, error) {
	return web.CreateCredential(user, *sessionData, response)
}
