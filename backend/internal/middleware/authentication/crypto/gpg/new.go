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
func NewEncryptor(publicKeys []string) (*Encryptor, error) {
	var validKeys []string

	for _, pubKey := range publicKeys {
		// Validate the public key
		key, err := crypto.NewKeyFromArmored(pubKey)
		if err != nil {
			continue // Skip invalid keys
		}

		// Check if the key can be used for encryption
		if key.CanEncrypt() {
			validKeys = append(validKeys, pubKey)
		}
	}

	if len(validKeys) == 0 {
		return nil, ErrorCantEncrypt
	}

	return &Encryptor{publicKey: validKeys}, nil
}
