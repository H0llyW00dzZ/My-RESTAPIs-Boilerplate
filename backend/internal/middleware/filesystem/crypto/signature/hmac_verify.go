// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package signature

import (
	"crypto/subtle"
	"encoding/hex"
)

// VerifyHMACSignatureFromFile verifies the HMAC signature for the given file content, secret key, and hex signature.
//
// Note: This function is suitable for automatically generated files such as backups, code generation, or mirrored files for frontend that are used by goroutine schedulers.
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
	// Note: This is generally safe from timing attacks as it uses constant-time comparison.
	return subtle.ConstantTimeCompare(expectedSignature, providedSignature) == 1, nil
}
