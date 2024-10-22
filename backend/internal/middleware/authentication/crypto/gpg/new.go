// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// Encryptor handles encryption using a OpenPGP/GPG public key.
type Encryptor struct {
	publicKey string
}

// NewEncryptor creates a new Encryptor instance.
func NewEncryptor(publicKey string) (*Encryptor, error) {
	// Validate the public key
	_, err := crypto.NewKeyFromArmored(publicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	return &Encryptor{publicKey: publicKey}, nil
}
