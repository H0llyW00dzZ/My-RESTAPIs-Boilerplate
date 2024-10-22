// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// Encryptor handles encryption using a OpenPGP/GPG public key.
type Encryptor struct {
	publicKey string
}

var (
	// ErrorCantEncrypt is returned when the provided public key cannot be used for encryption.
	ErrorCantEncrypt = errors.New("Crypto: GPG/OpenPGP the provided key cannot be used for encryption")
)

// NewEncryptor creates a new Encryptor instance.
func NewEncryptor(publicKey string) (*Encryptor, error) {
	// Validate the public key
	key, err := crypto.NewKeyFromArmored(publicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	// Check if the key can be used for encryption
	if !key.CanEncrypt() {
		return nil, ErrorCantEncrypt
	}

	return &Encryptor{publicKey: publicKey}, nil
}
