// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorOverflowOccurred is an error variable that represents an overflow condition during the conversion process.
// It is returned by the ToBytes function when the converted value exceeds the maximum value that can be stored in an int.
// This error indicates that the conversion failed due to an overflow and the result cannot be represented accurately.
var ErrorOverflowOccurred = errors.New("convert: overflow occurred during conversion")

// ToBytes converts a string representation of a size with units (KiB, MiB, GiB, TiB) to an integer number of bytes.
// It supports both floating-point and non-floating-point values.
//
// Example usage:
//
//	bytes, err := convert.ToBytes("10KiB")
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

	// Just in case it overflows, however it is not possible to write a test for this
	if bytes < 0 {
		return 0, ErrorOverflowOccurred
	}

	return bytes, nil
}
