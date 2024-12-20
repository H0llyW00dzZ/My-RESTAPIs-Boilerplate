// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert_test

import (
	"h0llyw00dz-template/backend/pkg/convert"
	"testing"
)

func TestToBytes(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"1024B", 1024, false},
		{"1024b", 1024, false},
		{"10KiB", 10240, false},
		{"5.2MiB", 5452595, false},
		{"3GiB", 3221225472, false},
		{"1.5TiB", 1649267441664, false},
		{"2.7GB", 2700000000, false},
		{"500MB", 500000000, false},
		{"1000B", 1000, false},
		{"1.5KB", 1500, false},
		{"2.5MB", 2500000, false},
		{"3.5GB", 3500000000, false},
		{"4.5TB", 4500000000000, false},
		{"invalid", 0, true},
		{"10XYZ", 0, true},
		{"", 0, true},
		{"1.5", 0, true},
		{"KB", 0, true},
		{"10InvalidUnit", 0, true},
		{"1,0ABC", 0, true},
		{"10,KiB", 0, true},
		{"10KiBMiBGiBTiB", 0, true},
		{"10kibmibgibtib", 0, true},
	}

	for _, tc := range testCases {
		result, err := convert.ToBytes(tc.input)
		if tc.hasError {
			if err == nil {
				t.Errorf("Expected an error for input: %s", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input: %s, error: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("Unexpected result for input: %s, got: %d, expected: %d", tc.input, result, tc.expected)
			}
		}
	}
}
