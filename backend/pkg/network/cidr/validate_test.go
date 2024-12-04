// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package cidr_test

import (
	"h0llyw00dz-template/backend/pkg/network/cidr"
	"os"
	"testing"
)

func TestValidateAndParseIPs(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		envVarValue string
		defaultIPS  string
		expected    []string
		expectError bool
	}{
		{
			name:        "Valid single IP",
			envVarValue: "192.168.1.1",
			defaultIPS:  "0.0.0.0/0",
			expected:    []string{"192.168.1.1"},
			expectError: false,
		},
		{
			name:        "Valid CIDR",
			envVarValue: "192.168.1.0/24",
			defaultIPS:  "0.0.0.0/0",
			expected:    []string{"192.168.1.0/24"},
			expectError: false,
		},
		{
			name:        "Multiple valid entries",
			envVarValue: "192.168.1.1, 10.0.0.0/8",
			defaultIPS:  "0.0.0.0/0",
			expected:    []string{"192.168.1.1", "10.0.0.0/8"},
			expectError: false,
		},
		{
			name:        "Invalid IP",
			envVarValue: "999.999.999.999",
			defaultIPS:  "0.0.0.0/0",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid CIDR",
			envVarValue: "999.999.999.999/999",
			defaultIPS:  "0.0.0.0/0",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Multiple invalid entries",
			envVarValue: "999.999.999.999, 999.999.999.999/999",
			defaultIPS:  "0.0.0.0/0",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Multiple invalid entries, expected one valid",
			envVarValue: "999.999.999.999, 999.999.999.999/999, 192.168.1.1, 100.100.100.100, 100.100.100.100/100",
			defaultIPS:  "0.0.0.0/0",
			expected:    []string{"192.168.1.1"},
			expectError: true,
		},
		{
			name:        "Default value used",
			envVarValue: "0.0.0.0/0",
			defaultIPS:  "0.0.0.0/0",
			expected:    []string{"0.0.0.0/0"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable
			os.Setenv("TEST_ENV_VAR", tt.envVarValue)
			defer os.Unsetenv("TEST_ENV_VAR")

			// Call the function
			result, err := cidr.ValidateAndParseIPs("TEST_ENV_VAR", tt.defaultIPS)

			// Check for errors
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, but got: %v", err)
				}
				// Check if the result matches the expected output
				if !equal(result, tt.expected) {
					t.Errorf("expected %v, but got %v", tt.expected, result)
				}
			}
		})
	}
}

// Helper function to compare slices
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
