// Copyright (c) 2025 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const (
	lowercaseCharset    = "abcdefghijklmnopqrstuvwxyz"
	uppercaseCharset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	mixedCharset        = lowercaseCharset + uppercaseCharset + numberCharset
	specialCharset      = "!@#$%^&*()-_=+[]{}|;:,.<>?/\\"
	mixedSpecialCharset = mixedCharset + specialCharset
	numberCharset       = "0123456789"
)

// TextCase defines the type for specifying text case.
type TextCase int

const (
	// Lowercase generates text using only lowercase letters.
	Lowercase TextCase = iota
	// Uppercase generates text using only uppercase letters.
	Uppercase
	// Mixed generates text using a mix of lowercase, uppercase letters, and numbers.
	Mixed
	// Special generates text using only special characters.
	Special
	// MixedSpecial generates text using a mix of lowercase, uppercase letters, numbers, and special characters.
	MixedSpecial
	// Number generates text using only numeric characters.
	Number
)

var (
	// ErrorsGenerateText is returned when an invalid text case is provided to GenerateText.
	// This error indicates that the specified TextCase does not match any of the predefined cases.
	ErrorsGenerateText = errors.New("crypto/rand: invalid text case")
)

// charsets maps each TextCase to its corresponding character set.
// This allows GenerateText to select the appropriate characters
// based on the specified textCase, ensuring flexibility in generating
// different types of random strings (e.g., lowercase, uppercase, mixed, etc.).
//
// Note: There is a trade-off in performance. Previously, using many switch cases
// could lead to complexity issues, especially when adding more character sets.
// However, using a map reduces complexity and improves maintainability, organization, and scalability.
// Although there might be a slight performance decrease during benchmarking,
// the benefits of avoiding excessive if-else or switch statements outweigh it.
var charsets = map[TextCase]string{
	Lowercase:    lowercaseCharset,
	Uppercase:    uppercaseCharset,
	Mixed:        mixedCharset,
	Special:      specialCharset,
	MixedSpecial: mixedSpecialCharset,
	Number:       numberCharset,
}

// GenerateText generates a random text string of the specified length and case.
func GenerateText(length int, textCase TextCase) (string, error) {
	// Note: This is not explicitly enforced because if a length of 1 is predictable,
	// it's generally not a security issue. The reason it's not explicitly set,
	// for example, with a minimum greater than 5, is that this function is used
	// not only for strong random generation (e.g., password generation) but for other purposes as well.
	if length <= 0 {
		return "", fmt.Errorf("crypto/rand: length %d must be greater than 0", length)
	}

	charset, exists := charsets[textCase]
	if !exists {
		return "", ErrorsGenerateText
	}

	text := make([]byte, length)
	charsetLen := int64(len(charset))
	for i := range text {
		// Note: This method is cryptographically secure. The randomness is unpredictable,
		// and no one can predict it. It is safe against classical side-channel attacks.
		// While quantum computing poses challenges to certain cryptographic algorithms,
		// the generation of random numbers itself remains secure with current quantum capabilities.
		index, err := rand.Int(rand.Reader, big.NewInt(charsetLen))
		if err != nil {
			return "", fmt.Errorf("crypto/rand: failed to generate random text: %w", err)
		}
		text[i] = charset[index.Int64()]
	}

	return string(text), nil
}
