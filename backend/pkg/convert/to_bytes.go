// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"fmt"
	"strings"
)

// ToBytes converts a string representation of a size with units (KiB, MiB, GiB, TiB) to an integer number of bytes.
// It supports both floating-point and non-floating-point values.
//
// Example usage:
//
//	bytes, err := ToBytes("10KiB")
//	if err != nil {
//	    // Handle the error you poggers
//	}
//	fmt.Printf("10KiB = %d bytes\n", bytes)
//
// The function returns an error if the input string has an invalid format or an unsupported unit.
//
// TODO: Switch to []string instead of string to allow bulk operations?
func ToBytes(size string) (int, error) {
	size = strings.TrimSpace(size)

	numericPart, unitPart := extractParts(size)

	if numericPart == "" {
		return 0, fmt.Errorf("invalid size: %s", size)
	}

	num, err := parseNumericPart(numericPart)
	if err != nil {
		return 0, fmt.Errorf("invalid size: %s", size)
	}

	bytes, err := convertToBytes(num, unitPart)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}
