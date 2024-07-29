// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package helper

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateAPIKey generates a secure API key using cryptographic techniques.
// It accepts optional parameters:
//
//	length to specify the length of the random byte slice (default is 50)
//	prefix to specify a custom prefix for the API key (default is "sk-")
//
// Note: When using this function for cryptographic purposes (e.g., encryption, decryption, or authenticated signing) by combining logic for enhances security,
// consider setting the length to 32 instead of the default value of 50.
func GenerateAPIKey(options ...any) string {
	// Set the default length to 50 and default prefix to "sk-"
	apiKeyLength := 50
	prefix := "sk-"

	// Parse the optional parameters
	for _, option := range options {
		switch v := option.(type) {
		case int:
			apiKeyLength = v
		case string:
			prefix = v
		default:
			panic("Invalid option type for GenerateAPIKey")
		}
	}

	// Generate a random byte slice of the specified length
	randomBytes := make([]byte, apiKeyLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Note: This is not possible.
		panic(err)
	}

	// Encode the random bytes using base64
	apiKey := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Prepend the prefix to the API key
	return prefix + apiKey
}
