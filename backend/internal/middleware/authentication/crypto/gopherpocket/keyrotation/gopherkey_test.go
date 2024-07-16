// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyrotation_test

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"testing"
	"time"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/gopherpocket/keyrotation"
)

func TestNewKeyManager(t *testing.T) {
	log.InitializeLogger("Gopher Testing", "unix")
	validAESKey := make([]byte, keyrotation.KeySize)
	validChaChaKey := make([]byte, keyrotation.KeySize)
	invalidKey := make([]byte, keyrotation.KeySize-1)

	testCases := []struct {
		name             string
		initialAESKey    []byte
		initialChaChaKey []byte
		rotationInterval time.Duration
		expectedError    error
	}{
		{
			name:             "Valid keys",
			initialAESKey:    validAESKey,
			initialChaChaKey: validChaChaKey,
			rotationInterval: time.Hour,
			expectedError:    nil,
		},
		{
			name:             "Invalid AES key size",
			initialAESKey:    invalidKey,
			initialChaChaKey: validChaChaKey,
			rotationInterval: time.Hour,
			expectedError:    keyrotation.ErrInvalidKeySize,
		},
		{
			name:             "Invalid ChaCha key size",
			initialAESKey:    validAESKey,
			initialChaChaKey: invalidKey,
			rotationInterval: time.Hour,
			expectedError:    keyrotation.ErrInvalidKeySize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := keyrotation.NewKeyManager(tc.initialAESKey, tc.initialChaChaKey, tc.rotationInterval)
			if err != tc.expectedError {
				t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
			}
		})
	}
}

func TestKeyRotation(t *testing.T) {
	initialAESKey := make([]byte, keyrotation.KeySize)
	initialChaChaKey := make([]byte, keyrotation.KeySize)
	rotationInterval := 100 * time.Millisecond

	km, err := keyrotation.NewKeyManager(initialAESKey, initialChaChaKey, rotationInterval)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}

	go km.StartKeyRotation()

	// Wait for a few rotation intervals
	time.Sleep(3 * rotationInterval)

	// Check that the keys have been rotated
	currentAESKey, currentChaChaKey := km.GetCurrentKeys()
	if string(currentAESKey) == string(initialAESKey) {
		t.Error("AES key has not been rotated")
	}
	if string(currentChaChaKey) == string(initialChaChaKey) {
		t.Error("ChaCha key has not been rotated")
	}
}
