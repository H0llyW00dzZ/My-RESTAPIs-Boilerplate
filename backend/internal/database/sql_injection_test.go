// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package database_test

import (
	"h0llyw00dz-template/backend/internal/database"
	"testing"
)

func TestIsValidTableName(t *testing.T) {
	// Add more SQL injection commands for testing if needed.
	// If something is valid that shouldn't be, it indicates a vulnerability SQL injection in the regex pattern.
	tests := []struct {
		name      string
		tableName string
		expected  bool
	}{
		{"ValidTableName", "Users", true},
		{"ValidWithUnderscore", "user_data", true},
		{"InvalidWithSemicolon", "Users; DROP TABLE Accounts;", false},
		{"SQLInjectionSleep", "Users; SLEEP(10); --", false},
		{"InvalidWithSpace", "user data", false},
		{"InvalidWithDash", "user-data", false},
		{"EmptyString", "", false},
		{"SQLInjectionUnion", "Users; UNION SELECT * FROM Admins; --", false},
		{"SQLInjectionComment", "Users; --", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := database.IsValidTableName(tt.tableName); got != tt.expected {
				t.Errorf("IsValidTableName(%q) = %v; want %v", tt.tableName, got, tt.expected)
			}
		})
	}
}
