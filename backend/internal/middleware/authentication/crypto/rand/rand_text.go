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
	mixedCharset        = lowercaseCharset + uppercaseCharset + "0123456789"
	specialCharset      = "!@#$%^&*()-_=+[]{}|;:,.<>?/\\"
	mixedSpecialCharset = mixedCharset + specialCharset
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
)

var (
	// ErrorsGenerateText is returned when an invalid text case is provided to GenerateText.
	// This error indicates that the specified TextCase does not match any of the predefined cases.
	ErrorsGenerateText = errors.New("crypto/rand: invalid text case")
)

// GenerateText generates a random text string of the specified length and case.
func GenerateText(length int, textCase TextCase) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("crypto/rand: length %d must be greater than 0", length)
	}

	// Note: This implementation is optimized for performance.
	// Another method could use a map for better organization, like this:
	//
	// 	var charsets = map[TextCase]string{
	// 		Lowercase:    lowercaseCharset,
	// 		Uppercase:    uppercaseCharset,
	// 		Mixed:        mixedCharset,
	// 		Special:      specialCharset,
	// 		MixedSpecial: mixedSpecialCharset,
	// 	}
	//
	// 	// GenerateText generates a random text string of the specified length and case.
	// 	func GenerateText(length int, textCase TextCase) (string, error) {
	// 		if length <= 0 {
	// 			return "", fmt.Errorf("crypto/rand: length %d must be greater than 0", length)
	// 		}
	//
	// 		charset, exists := charsets[textCase]
	// 		if !exists {
	// 			return "", ErrorsGenerateText
	// 		}
	//
	// 		charsetLen := int64(len(charset))
	// 		text := make([]byte, length)
	// 		for i := range text {
	// 			index, err := rand.Int(rand.Reader, big.NewInt(charsetLen))
	// 			if err != nil {
	// 				return "", fmt.Errorf("crypto/rand: failed to generate random text: %w", err)
	// 			}
	// 			text[i] = charset[index.Int64()]
	// 		}
	//
	// 		return string(text), nil
	// 	}
	//
	// However, using a map might slightly decrease performance due to additional lookups.
	var charset string
	switch textCase {
	case Lowercase:
		charset = lowercaseCharset
	case Uppercase:
		charset = uppercaseCharset
	case Mixed:
		charset = mixedCharset
	case Special:
		charset = specialCharset
	case MixedSpecial:
		charset = mixedSpecialCharset
	default:
		return "", ErrorsGenerateText
	}

	text := make([]byte, length)
	for i := range text {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("crypto/rand: failed to generate random text: %w", err)
		}
		text[i] = charset[index.Int64()]
	}

	return string(text), nil
}
