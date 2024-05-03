// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package signature

import (
	"crypto/subtle"
	"encoding/hex"
)

// VerifyHMACSignatureFromFile verifies the HMAC signature for the given file content, secret key, and hex signature.
func VerifyHMACSignatureFromFile(filePath, secretKey, hexSignature string) (bool, error) {
	// Generate the expected HMAC signature from the file content
	expectedSignature, err := GenerateHMACSignatureFromFile(filePath, secretKey)
	if err != nil {
		return false, err
	}

	// Decode the provided hex signature
	providedSignature, err := hex.DecodeString(hexSignature)
	if err != nil {
		return false, err
	}

	// Compare the expected signature with the provided signature using ConstantTimeCompare
	return subtle.ConstantTimeCompare(expectedSignature, providedSignature) == 1, nil
}
