// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package helper

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"testing"
)

func TestParseNumericalValue(t *testing.T) {
	log.InitializeLogger("Gopher Testing", "unix")
	tests := []struct {
		name     string
		value    string
		base     int
		bitSize  int
		expected uint64
	}{
		{
			name:     "Valid uint64",
			value:    "1234567890",
			base:     10,
			bitSize:  64,
			expected: 1234567890,
		},
		{
			name:     "Invalid uint64",
			value:    "invalid",
			base:     10,
			bitSize:  64,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseNumericalValue(tt.value, tt.base, tt.bitSize)
			if result != tt.expected {
				t.Errorf("ParseNumericalValue(%s, %d, %d) = %d, want %d", tt.value, tt.base, tt.bitSize, result, tt.expected)
			}
		})
	}
}

func TestParseInt64Value(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		base     int
		bitSize  int
		expected int64
	}{
		{
			name:     "Valid int64",
			value:    "-1234567890",
			base:     10,
			bitSize:  64,
			expected: -1234567890,
		},
		{
			name:     "Invalid int64",
			value:    "invalid",
			base:     10,
			bitSize:  64,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseInt64Value(tt.value, tt.base, tt.bitSize)
			if result != tt.expected {
				t.Errorf("ParseInt64Value(%s, %d, %d) = %d, want %d", tt.value, tt.base, tt.bitSize, result, tt.expected)
			}
		})
	}
}
