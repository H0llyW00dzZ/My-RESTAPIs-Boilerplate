// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand_test

import (
	"crypto/elliptic"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: If this function is not safe for concurrent use, it would lead to race conditions or produce the same result
// across multiple goroutines. This function is particularly useful in production environments where the system requires
// multiple goroutines (e.g., when 10000 goroutines are needed, it provides 10000 readers, each with its own instance so it always welcome).
func TestFixedSize32Bytes(t *testing.T) {
	r := rand.FixedSize32Bytes()

	// Test reading from the reader
	buf := make([]byte, 32)
	n, err := r.Read(buf)

	// Check the number of bytes read
	if n != 32 {
		t.Errorf("Expected to read 32 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test reading again to ensure it generates new random bytes
	buf2 := make([]byte, 32)
	_, err = r.Read(buf2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that the two reads generate different random bytes
	if string(buf) == string(buf2) {
		t.Error("Expected different random bytes on subsequent reads")
	}
}

// Note: If this function is not safe for concurrent use, it would lead to race conditions or produce the same result
// across multiple goroutines. This function is particularly useful in production environments where the system requires
// multiple goroutines (e.g., when 10000 goroutines are needed, it provides 10000 readers, each with its own instance so it always welcome).
func TestFixedReaderRead(t *testing.T) {
	r := rand.FixedSize32Bytes()

	// Test reading with a buffer smaller than the fixed size
	buf := make([]byte, 16)
	n, err := r.Read(buf)

	// Check the number of bytes read
	if n != 16 {
		t.Errorf("Expected to read 16 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test reading with a buffer larger than the fixed size
	buf = make([]byte, 64)
	n, err = r.Read(buf)

	// Check the number of bytes read
	if n != 32 {
		t.Errorf("Expected to read 32 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test reading with an empty buffer
	buf = make([]byte, 0)
	n, err = r.Read(buf)

	// Check the number of bytes read
	if n != 0 {
		t.Errorf("Expected to read 0 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Note: If this function is not safe for concurrent use, it would lead to race conditions or produce the same result
// across multiple goroutines. This function is particularly useful in production environments where the system requires
// multiple goroutines (e.g., when 10000 goroutines are needed, it provides 10000 readers, each with its own instance so it always welcome).
func TestFixedSizeECDSA(t *testing.T) {
	curves := []elliptic.Curve{
		elliptic.P256(),
		elliptic.P384(),
		elliptic.P521(),
	}

	for _, curve := range curves {
		r := rand.FixedSizeECDSA(curve)

		// Determine the expected size based on the curve
		expectedSize := 32
		if curve.Params().BitSize > 256 {
			expectedSize = 48
		}

		// Test reading from the reader
		buf := make([]byte, expectedSize)
		n, err := r.Read(buf)

		// Check the number of bytes read
		if n != expectedSize {
			t.Errorf("Expected to read %d bytes for curve %s, but read %d bytes", expectedSize, curve.Params().Name, n)
		}

		// Check for any errors
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", curve.Params().Name, err)
		}

		// Test reading again to ensure it generates new random bytes
		buf2 := make([]byte, expectedSize)
		_, err = r.Read(buf2)
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", curve.Params().Name, err)
		}

		// Check that the two reads generate different random bytes
		if string(buf) == string(buf2) {
			t.Errorf("Expected different random bytes on subsequent reads for curve %s", curve.Params().Name)
		}
	}
}

func TestFixedSizeECC(t *testing.T) {
	curves := []struct {
		curve        elliptic.Curve
		expectedSize int
	}{
		{elliptic.P224(), 28},
		{elliptic.P256(), 32},
		{elliptic.P384(), 48},
		{elliptic.P521(), 66},
	}

	for _, c := range curves {
		r := rand.FixedSizeECC(c.curve)

		// Test reading from the reader
		buf := make([]byte, c.expectedSize)
		n, err := r.Read(buf)

		// Check the number of bytes read
		if n != c.expectedSize {
			t.Errorf("Expected to read %d bytes for curve %s, but read %d bytes", c.expectedSize, c.curve.Params().Name, n)
		}

		// Check for any errors
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", c.curve.Params().Name, err)
		}

		// Test reading again to ensure it generates new random bytes
		buf2 := make([]byte, c.expectedSize)
		_, err = r.Read(buf2)
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", c.curve.Params().Name, err)
		}

		// Check that the two reads generate different random bytes
		if string(buf) == string(buf2) {
			t.Errorf("Expected different random bytes on subsequent reads for curve %s", c.curve.Params().Name)
		}
	}
}

// TestFixedSizeRSA tests the FixedSizeRSA function to ensure it generates the correct number of random bytes.
func TestFixedSizeRSA(t *testing.T) {
	modulusSizes := []struct {
		modulusBits  int
		expectedSize int
	}{
		{1024, 128},
		{2048, 256},
		{3072, 384},
		{4096, 512},
	}

	for _, ms := range modulusSizes {
		r := rand.FixedSizeRSA(ms.modulusBits)

		// Test reading from the reader
		buf := make([]byte, ms.expectedSize)
		n, err := r.Read(buf)

		// Check the number of bytes read
		if n != ms.expectedSize {
			t.Errorf("Expected to read %d bytes for modulus %d bits, but read %d bytes", ms.expectedSize, ms.modulusBits, n)
		}

		// Check for any errors
		if err != nil {
			t.Errorf("Unexpected error for modulus %d bits: %v", ms.modulusBits, err)
		}

		// Test reading again to ensure it generates new random bytes
		buf2 := make([]byte, ms.expectedSize)
		_, err = r.Read(buf2)
		if err != nil {
			t.Errorf("Unexpected error for modulus %d bits: %v", ms.modulusBits, err)
		}

		// Check that the two reads generate different random bytes
		if string(buf) == string(buf2) {
			t.Errorf("Expected different random bytes on subsequent reads for modulus %d bits", ms.modulusBits)
		}
	}
}

func TestGenerateFixedUUID(t *testing.T) {
	// Define a regex pattern for a valid UUID v4 and variant (RFC 4122) bits.
	uuidPattern := `^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`
	re := regexp.MustCompile(uuidPattern)

	// Define a regex pattern for a valid UUID v4 and variant (RFC 4122) bits without hyphens.
	uuidPatternWithoutHyphens := `^[a-f0-9]{32}$`
	reWithoutHyphens := regexp.MustCompile(uuidPatternWithoutHyphens)

	tests := []struct {
		name          string
		removeHyphens bool
		pattern       *regexp.Regexp
	}{
		{"WithHyphens", false, re},
		{"WithoutHyphens", true, reWithoutHyphens},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate a UUID and check for errors.
			uuid, err := rand.GenerateFixedUUID(rand.UUIDFormat{RemoveHyphens: tt.removeHyphens})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			t.Logf("Generated UUID: %s", uuid)

			// Verify the format of the UUID.
			if !tt.pattern.MatchString(uuid) {
				t.Errorf("UUID %s does not match expected format", uuid)
			}

			// Generate another UUID and ensure it is different from the first.
			uuid2, err := rand.GenerateFixedUUID(rand.UUIDFormat{RemoveHyphens: tt.removeHyphens})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			t.Logf("Generated UUID: %s", uuid2)

			if uuid == uuid2 {
				t.Error("expected different UUIDs, but got the same")
			}
		})
	}
}

func textCaseToString(tc rand.TextCase) string {
	switch tc {
	case rand.Lowercase:
		return "Lowercase"
	case rand.Uppercase:
		return "Uppercase"
	case rand.Mixed:
		return "Mixed"
	case rand.Special:
		return "Special"
	case rand.MixedSpecial:
		return "MixedSpecial"
	case rand.Number:
		return "Number"
	case rand.UpperNumCase:
		return "UpperNumCase"
	case rand.LowerNumCase:
		return "LowerNumCase"
	case rand.NumSpecial:
		return "NumSpecial"
	case rand.LowercaseSpecial:
		return "LowercaseSpecial"
	case rand.UppercaseSpecial:
		return "UppercaseSpecial"
	default:
		return "Unknown"
	}
}

// Note: Avoid using hardcoded numbers directly as the textCase parameter (e.g., rand.GenerateText(10, 0)).
// Instead, use predefined constants (e.g., rand.GenerateText(10, rand.Lowercase)).
// This approach is considered idiomatic in Go for better readability and maintainability.
// The constants are bound to the TextCase type.
func TestGenerateRandomText(t *testing.T) {
	tests := []struct {
		length   int
		textCase rand.TextCase
	}{
		{10, rand.Lowercase},
		{10, rand.Uppercase},
		{10, rand.Mixed},
		{10, rand.Special},
		{10, rand.MixedSpecial},
		{10, rand.Number},
		{10, rand.UpperNumCase},
		{10, rand.LowerNumCase},
		{10, rand.NumSpecial},
		{10, rand.LowercaseSpecial},
		{10, rand.UppercaseSpecial},
	}

	for _, tt := range tests {
		t.Run(textCaseToString(tt.textCase), func(t *testing.T) {
			result, err := rand.GenerateText(tt.length, tt.textCase)
			assert.NoError(t, err)
			assert.Equal(t, tt.length, len(result))

			// Log the generated text
			t.Logf("Generated text for %s: %s", textCaseToString(tt.textCase), result)

			switch tt.textCase {
			case rand.Lowercase:
				assert.Regexp(t, "^[a-z]+$", result)
			case rand.Uppercase:
				assert.Regexp(t, "^[A-Z]+$", result)
			case rand.Mixed:
				assert.Regexp(t, "^[a-zA-Z0-9]+$", result)
			case rand.Special:
				assert.Regexp(t, "^[!@#$%^&*()\\-_=+\\[\\]{}|;:,.<>?/\\\\]+$", result)
			case rand.MixedSpecial:
				assert.Regexp(t, "^[a-zA-Z0-9!@#$%^&*()\\-_=+\\[\\]{}|;:,.<>?/\\\\]+$", result)
			case rand.Number:
				assert.Regexp(t, "^[0-9]+$", result)
			case rand.UpperNumCase:
				assert.Regexp(t, "^[A-Z0-9]+$", result)
			case rand.LowerNumCase:
				assert.Regexp(t, "^[a-z0-9]+$", result)
			case rand.NumSpecial:
				assert.Regexp(t, "^[0-9!@#$%^&*()\\-_=+\\[\\]{}|;:,.<>?/\\\\]+$", result)
			case rand.LowercaseSpecial:
				assert.Regexp(t, "^[a-z!@#$%^&*()\\-_=+\\[\\]{}|;:,.<>?/\\\\]+$", result)
			case rand.UppercaseSpecial:
				assert.Regexp(t, "^[A-Z!@#$%^&*()\\-_=+\\[\\]{}|;:,.<>?/\\\\]+$", result)
			}
		})
	}
}

func TestGenerateRandomTextInvalidInputs(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		textCase rand.TextCase
		expected string
	}{
		{
			name:     "Negative Length",
			length:   -1,
			textCase: rand.Mixed,
			expected: "length -1 must be greater than 0",
		},
		{
			name:     "Invalid TextCase",
			length:   10,
			textCase: -1, // Assuming -1 is an invalid textCase
			expected: rand.ErrorsGenerateText.Error(),
		},
		{
			name:     "Negative Length and Invalid TextCase",
			length:   -1,
			textCase: -1,
			expected: "length -1 must be greater than 0",
		},
		{
			name:     "Invalid TextCase (999999999999999999)",
			length:   10,
			textCase: 999999999999999999, // Assuming 999999999999999999 is an invalid textCase
			expected: rand.ErrorsGenerateText.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rand.GenerateText(tt.length, tt.textCase)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestChoice(t *testing.T) {
	choices := []string{"apple", "banana", "cherry"}
	results := make(map[string]bool)

	// Run the test multiple times to ensure randomness
	for range 100 {
		choice, err := rand.Choice(choices)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check if the choice is within the expected set
		if choice != "apple" && choice != "banana" && choice != "cherry" {
			t.Errorf("unexpected choice: %v", choice)
		}

		// Record the result
		results[choice] = true
	}

	// Ensure all choices have been selected at least once
	if len(results) != len(choices) {
		t.Errorf("not all choices were selected: %v", results)
	}
}
