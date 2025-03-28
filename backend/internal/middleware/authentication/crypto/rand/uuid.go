// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand

import (
	"fmt"
)

// UUIDFormat represents the format options for generating UUIDs.
type UUIDFormat struct {
	RemoveHyphens bool
}

// GenerateFixedUUID creates a new UUID using [crypto/rand] for high randomness.
//
// Note: Unlike most UUID implementations bound by RFC standards,
// this is purely random and not bound to any specific format/resource (e.g., disk (serial), memory, clock, other hardware id).
//
// The format of the generated UUID can be customized using the UUIDFormat struct.
// If no format options are provided, the default format with hyphens will be used.
func GenerateFixedUUID(format ...UUIDFormat) (string, error) {
	// Set default format options if none are provided
	opts := UUIDFormat{
		RemoveHyphens: false,
	}
	if len(format) > 0 {
		opts = format[0]
	}

	// Use FixedSizeReader to ensure a consistent 16-byte read.
	reader := FixedSizeReader(16)
	uuid := make([]byte, 16)
	if _, err := reader.Read(uuid); err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Set the version (4) and variant (RFC 4122) bits.
	//
	// Example:
	// Generated UUID: 14215a72-cebd-4b3a-9d98-86cde9c261e0
	//
	// Note: This is similar to Google's UUID, but it does not use a pool (mutex); it directly generates random numbers.
	// It is safe to call from multiple goroutines simultaneously.
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	if opts.RemoveHyphens {
		// Construct the UUID string without hyphens directly
		// Note: This much faster.
		return fmt.Sprintf("%08x%04x%04x%04x%012x",
			uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
	}

	// Construct the UUID string with hyphens
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
