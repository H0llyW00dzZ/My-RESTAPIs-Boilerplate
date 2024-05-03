// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"os"
)

// GenerateHMACSignatureFromFile generates an HMAC signature for the given file content and secret key.
//
// Note: This function is suitable for automatically generated files such as backups, code generation, or mirrored files for frontend that are used by goroutine schedulers.
func GenerateHMACSignatureFromFile(filePath, secretKey string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new HMAC instance with the secret key
	mac := hmac.New(sha256.New, []byte(secretKey))

	// Copy the file content to the HMAC instance
	if _, err := io.Copy(mac, file); err != nil {
		return nil, err
	}

	// Calculate the HMAC signature
	signature := mac.Sum(nil)

	return signature, nil
}
