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
	uppernumcaseCharset = uppercaseCharset + numberCharset
	lowernumcaseCharset = lowercaseCharset + numberCharset
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
	// UpperNumCase generates text using a mix of uppercase letters and numbers.
	UpperNumCase
	// LowerNumCase generates text using a mix of lowercase letters and numbers.
	LowerNumCase
)

var (
	// ErrorsGenerateText is returned when an invalid text case is provided to GenerateText.
	// This error indicates that the specified TextCase does not match any of the predefined cases.
	ErrorsGenerateText = errors.New("crypto/rand: invalid text case")
)

// charsets is a slice that maps each TextCase to its corresponding character set.
// This allows GenerateText to select the appropriate characters based on the specified
// textCase, ensuring flexibility in generating different types of random strings
// (e.g., lowercase, uppercase, mixed, etc.).
//
// Note: A slice is used here instead of a map for better performance and simplicity.
// Since the TextCase values are small, fixed, and sequential integers (starting from 0),
// accessing the slice by index is a constant-time operation (O(1)) and more memory-efficient
// than a map. Additionally, the bounds check ensures safety by preventing out-of-bounds access.
//
// This approach avoids the complexity of using many switch or if-else statements, while
// maintaining high performance and clear organization. It is well-suited for cases where
// the number of TextCase values is small and fixed.
var charsets = []string{
	lowercaseCharset,    // Lowercase
	uppercaseCharset,    // Uppercase
	mixedCharset,        // Mixed
	specialCharset,      // Special
	mixedSpecialCharset, // MixedSpecial
	numberCharset,       // Number
	uppernumcaseCharset, // UpperNumCase
	lowernumcaseCharset, // LowerNumCase
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

	// Ensure the textCase is within bounds
	if int(textCase) < 0 || int(textCase) >= len(charsets) {
		return "", ErrorsGenerateText
	}

	charset := charsets[textCase]
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
