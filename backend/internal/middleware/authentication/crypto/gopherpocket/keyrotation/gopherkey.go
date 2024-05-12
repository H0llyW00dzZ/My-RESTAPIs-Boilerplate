// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package keyrotation

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	log "h0llyw00dz-template/backend/internal/logger"
	"sync"
	"time"
)

const (
	// KeySize specifies the size of the encryption keys in bytes.
	// It is set to 32 bytes (256 bits) for both AES and ChaCha20-Poly1305 keys.
	KeySize = 32
)

var (
	// ErrInvalidKeySize is returned when the provided initial key sizes are invalid.
	// It indicates that the initial AES and ChaCha20-Poly1305 keys must be 32 bytes each.
	ErrInvalidKeySize = errors.New("invalid initial key size")

	// ErrKeyRotationFailed is returned when an error occurs during key rotation.
	// It signifies that the key rotation process encountered an error and failed.
	ErrKeyRotationFailed = errors.New("key rotation failed")
)

// KeyManager manages the rotation of encryption keys.
type KeyManager struct {
	currentAESKey    []byte
	currentChaChaKey []byte
	rotationInterval time.Duration
	mutex            sync.RWMutex
}

// NewKeyManager creates a new instance of KeyManager with the provided initial keys and rotation interval.
// It returns an error if the initial key sizes are invalid.
func NewKeyManager(initialAESKey, initialChaChaKey []byte, rotationInterval time.Duration) (*KeyManager, error) {
	if len(initialAESKey) != KeySize || len(initialChaChaKey) != KeySize {
		return nil, ErrInvalidKeySize
	}
	return &KeyManager{
		currentAESKey:    initialAESKey,
		currentChaChaKey: initialChaChaKey,
		rotationInterval: rotationInterval,
	}, nil
}

// GetCurrentKeys returns the current AES and ChaCha20-Poly1305 keys.
// It uses a read lock to ensure thread-safe access to the keys and performs constant-time comparison.
func (km *KeyManager) GetCurrentKeys() ([]byte, []byte) {
	km.mutex.RLock()
	defer km.mutex.RUnlock()

	aesKey := make([]byte, KeySize)
	chaChaKey := make([]byte, KeySize)

	subtle.ConstantTimeCopy(1, aesKey, km.currentAESKey)
	subtle.ConstantTimeCopy(1, chaChaKey, km.currentChaChaKey)

	return aesKey, chaChaKey
}

// RotateKeys generates new random keys and replaces the current keys with the new ones.
// It uses a write lock to ensure exclusive access during key rotation.
func (km *KeyManager) RotateKeys() error {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	newAESKey := make([]byte, KeySize)
	newChaChaKey := make([]byte, KeySize)

	_, err := rand.Read(newAESKey)
	if err != nil {
		return err
	}

	_, err = rand.Read(newChaChaKey)
	if err != nil {
		return err
	}

	km.currentAESKey = newAESKey
	km.currentChaChaKey = newChaChaKey
	return nil
}

// StartKeyRotation starts a background goroutine that periodically rotates the keys based on the specified rotation interval.
// It logs any errors that occur during key rotation.
func (km *KeyManager) StartKeyRotation() {
	ticker := time.NewTicker(km.rotationInterval)
	defer ticker.Stop()

	for range ticker.C {
		err := km.RotateKeys()
		if err != nil {
			log.LogErrorf("Key rotation failed: %v", err)
		}
	}
}
