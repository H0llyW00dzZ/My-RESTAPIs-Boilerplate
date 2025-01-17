// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"errors"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// Encryptor handles encryption using a OpenPGP/GPG public key.
type Encryptor struct {
	publicKey []string
	keyInfos  []KeyInfo
	config    *Config
	*helper
}

// helper provides utility functions that can be embedded in other structs.
// This allows for shared code functionality across different components without duplicating code.
// By embedding helper, structs can access common methods easily.
type helper struct {
	armor  bool
	suffix string
}

var (
	// ErrorCantEncrypt is returned when the provided public key cannot be used for encryption.
	ErrorCantEncrypt = errors.New("Crypto: GPG/OpenPGP the provided key cannot be used for encryption")
)

// NewEncryptor creates a new Encryptor instance with multiple public keys.
//
// Note: Ensure that the provided public key can be used for encryption.
// This function handles multiple keys within an armored key block.
// Filtering keys from a complex, multi-key armored block can be challenging.
//
// TODO: Implement similar logic for Verify/Sign mechanisms that can be used for authentication over the network (GPG Modern) ?
func NewEncryptor(publicKeys []string, opts ...Option) (*Encryptor, error) {
	// Apply user-provided options to override defaults
	config := NewDefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	var validKeys []string
	var keyInfos []KeyInfo
	// Track unique keys by fingerprint
	uniqueKeys := make(map[string]bool)

	for _, pubKey := range publicKeys {
		// Validate the public key
		key, err := crypto.NewKeyFromArmored(pubKey)
		if err != nil {
			continue // Skip invalid keys
		}

		// Extract key information
		keyInfo := extractKeyInfo(key)

		// Check if the key is already added using its fingerprint
		if uniqueKeys[keyInfo.Fingerprint] {
			continue // Skip duplicate keys
		}

		// Check if the key can be used for encryption or Check if the key can be used for future verification
		// prevent duplicate keys
		if key.CanEncrypt() || (config.AllowVerfy && key.CanVerify()) {
			validKeys = append(validKeys, pubKey)
			keyInfos = append(keyInfos, keyInfo)
			uniqueKeys[keyInfo.Fingerprint] = true
		}
	}

	if len(validKeys) == 0 {
		return nil, ErrorCantEncrypt
	}

	return &Encryptor{
		publicKey: validKeys,
		keyInfos:  keyInfos,
		config:    config,
		helper: &helper{
			armor:  config.armor,
			suffix: config.suffix,
		},
	}, nil
}

// extractKeyInfo extracts metadata from a given crypto.Key and returns it as a KeyInfo struct.
// This function gathers essential details about the key, such as its ID, capabilities, and fingerprints.
// If there were many keys used for end-to-end encryption, it recommended to store them in a database.
// This was primarily used for automated backups that sent data over the network due to its efficiency.
func extractKeyInfo(key *crypto.Key) KeyInfo {
	entity := key.GetEntity()
	var userIDs []string
	for _, uid := range entity.Identities {
		userIDs = append(userIDs, uid.UserId.Id)
	}

	return KeyInfo{
		UserIDs:           userIDs,
		PrimaryKey:        entity.PrimaryKey.KeyIdString(),
		CreationDate:      entity.PrimaryKey.CreationTime,
		KeyID:             key.GetKeyID(),
		HexKeyID:          key.GetHexKeyID(),
		CanEncrypt:        key.CanEncrypt(),
		CanVerify:         key.CanVerify(),
		IsExpired:         key.IsExpired(),
		IsRevoked:         key.IsRevoked(),
		Fingerprint:       key.GetFingerprint(),
		DigestFingerprint: key.GetSHA256Fingerprints(),
	}
}

const (
	newGPGModern = ".gpg"
)
