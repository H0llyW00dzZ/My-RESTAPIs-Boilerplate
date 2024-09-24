// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package webauthn provides a simple and convenient interface for implementing WebAuthn
// authentication in a Go application. It wraps the github.com/go-webauthn/webauthn library
// and provides a higher-level API for common WebAuthn operations.
//
// The package allows initializing a WebAuthn client, beginning and finishing registration, and
// beginning and finishing login. It handles the session data management and provides the necessary
// data structures for the WebAuthn protocol.
//
// To use this package, follow these steps:
//
// 1. Initialize the WebAuthn client by calling the Init function with a valid configuration.
//
//  2. Begin the registration process by calling the BeginRegistration function, passing the user
//     and any additional registration options. This function returns a CredentialCreation instance,
//     session data, and an error if any. Store the session data securely for the next step.
//
//  3. Complete the registration process by calling the FinishRegistration function, passing the user,
//     the session data obtained from the previous step, and the parsed credential creation response.
//     This function returns a Credential instance and an error if any.
//
//  4. Begin the login process by calling the BeginLogin function, passing the user and any additional
//     login options. This function returns a CredentialAssertion instance, session data, and an error
//     if any. Store the session data securely for the next step.
//
//  5. Complete the login process by calling the FinishLogin function, passing the user, the session
//     data obtained from the previous step, and the parsed credential assertion response. This function
//     returns a Credential instance and an error if any.
//
// Note: The session data returned by the BeginRegistration and BeginLogin functions must be stored
// securely and passed to the corresponding FinishRegistration and FinishLogin functions to complete
// the WebAuthn flow.
//
// For more information about WebAuthn and its usage, refer to the GitHub repository at
// github.com/go-webauthn/webauthn and the WebAuthn specification at https://www.w3.org/TR/webauthn-2/.
//
// Also note that WebAuthn can be risky if the data registered on the device gets compromised,
// potentially leading to malware or remote access trojan (RAT) attacks.
package webauthn
