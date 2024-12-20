// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"fmt"
	"strconv"
	"strings"
)

// extractParts extracts the numeric part and the unit part from the input string.
// It returns the numeric part and the unit part as separate strings.
func extractParts(size string) (string, string) {
	// Extract the numeric part and the unit part
	numericPart := strings.TrimRight(size, "kmgtKMGTibIB")
	unitPart := strings.TrimLeft(size, "0123456789.")
	return numericPart, unitPart
}

// parseNumericPart parses the numeric part of the input string into a float64.
// It returns the parsed float64 value and an error if the parsing fails.
func parseNumericPart(numericPart string) (float64, error) {
	// Convert the numeric part to a float64
	return strconv.ParseFloat(numericPart, 64)
}

// convertToBytes converts the numeric value to bytes based on the unit part.
// It returns the converted value in bytes and an error if the unit is invalid or unsupported.
func convertToBytes(num float64, unitPart string) (int, error) {
	// Define the conversion factors for each unit
	//
	// Note: This is based on how computers work (e.g., typically handle units of memory and storage) and should be correct.
	factors := map[string]int{
		"":    1,
		"k":   1024,
		"m":   1024 * 1024,
		"g":   1024 * 1024 * 1024,
		"t":   1024 * 1024 * 1024 * 1024,
		"kib": 1024,
		"mib": 1024 * 1024,
		"gib": 1024 * 1024 * 1024,
		"tib": 1024 * 1024 * 1024 * 1024,
		"kb":  1000,
		"mb":  1000 * 1000,
		"gb":  1000 * 1000 * 1000,
		"tb":  1000 * 1000 * 1000 * 1000,
	}

	// Convert the unit part to lowercase
	unitPart = strings.ToLower(unitPart)

	// Special case for "B" unit
	if unitPart == "b" {
		return int(num), nil
	}

	// Check if the unit part is empty
	if unitPart == "" {
		return 0, fmt.Errorf("invalid size: %f", num)
	}

	// Get the conversion factor based on the unit
	factor, ok := factors[unitPart]
	if !ok {
		return 0, fmt.Errorf("invalid unit: %s", unitPart)
	}

	// Calculate the size in bytes
	bytes := int(num * float64(factor))

	return bytes, nil
}
