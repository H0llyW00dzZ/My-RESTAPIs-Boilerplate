// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package helper

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"strconv"
)

// ParseNumericalValue is a helper function to parse numerical values from strings.
// It takes a string value, base, and bit size as input and returns the parsed value as uint64.
// If parsing fails, it returns 0.
func ParseNumericalValue(value string, base, bitSize int) uint64 {
	parsed, err := strconv.ParseUint(value, base, bitSize)
	if err != nil {
		log.LogErrorf("Failed to parse numerical value: %v", err)
		return 0
	}
	return parsed
}

// ParseInt64Value is a helper function to parse int64 values from strings.
// It takes a string value, base, and bit size as input and returns the parsed value as int64.
// If parsing fails, it returns 0.
func ParseInt64Value(value string, base, bitSize int) int64 {
	parsed, err := strconv.ParseInt(value, base, bitSize)
	if err != nil {
		log.LogErrorf("Failed to parse int64 value: %v", err)
		return 0
	}
	return parsed
}
