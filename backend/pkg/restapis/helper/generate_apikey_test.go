// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package helper_test

import (
	"strings"
	"testing"

	"h0llyw00dz-template/backend/pkg/restapis/helper"
)

func TestGenerateAPIKey_DefaultLengthAndPrefix(t *testing.T) {
	apiKey := helper.GenerateAPIKey()
	if len(apiKey) != 70 {
		t.Errorf("Expected API key length to be 70, but got %d", len(apiKey))
	}
	if !strings.HasPrefix(apiKey, "sk-") {
		t.Errorf("Expected API key to have prefix 'sk-', but got %s", apiKey)
	}
}

func TestGenerateAPIKey_CustomLength(t *testing.T) {
	apiKey := helper.GenerateAPIKey(64)
	if len(apiKey) != 89 {
		t.Errorf("Expected API key length to be 89, but got %d", len(apiKey))
	}
	if !strings.HasPrefix(apiKey, "sk-") {
		t.Errorf("Expected API key to have prefix 'sk-', but got %s", apiKey)
	}
}

func TestGenerateAPIKey_CustomPrefix(t *testing.T) {
	apiKey := helper.GenerateAPIKey("custom-")
	if len(apiKey) != 74 {
		t.Errorf("Expected API key length to be 74, but got %d", len(apiKey))
	}
	if !strings.HasPrefix(apiKey, "custom-") {
		t.Errorf("Expected API key to have prefix 'custom-', but got %s", apiKey)
	}
}

func TestGenerateAPIKey_CustomLengthAndPrefix(t *testing.T) {
	apiKey := helper.GenerateAPIKey(32, "api-")
	if len(apiKey) != 47 {
		t.Errorf("Expected API key length to be 47, but got %d", len(apiKey))
	}
	if !strings.HasPrefix(apiKey, "api-") {
		t.Errorf("Expected API key to have prefix 'api-', but got %s", apiKey)
	}
}

func TestGenerateAPIKey_InvalidOptionType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "Invalid option type for GenerateAPIKey" {
				t.Errorf("Expected panic message 'Invalid option type for GenerateAPIKey', but got '%v'", r)
			}
		} else {
			t.Errorf("Expected GenerateAPIKey to panic on invalid option type, but it didn't")
		}
	}()
	helper.GenerateAPIKey(true)
}
