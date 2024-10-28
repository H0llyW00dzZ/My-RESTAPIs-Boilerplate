// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand

import (
	"crypto/rand"
	"fmt"
)

// GenerateFixedUUID creates a new UUID using [crypto/rand] for high randomness.
//
// Note: Unlike most UUID implementations bound by RFC standards,
// this is purely random and not bound to any specific format/resource (e.g., disk (serial), memory, clock, other hardware id).
func GenerateFixedUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Set the version (4) and variant (RFC 4122) bits.
	//
	// Example:
	// Generated UUID: 14215a72-cebd-4b3a-9d98-86cde9c261e0
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}