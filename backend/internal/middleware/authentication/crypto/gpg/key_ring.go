// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// createKeyRing creates a KeyRing from the public key.
func (e *Encryptor) createKeyRing() (*crypto.KeyRing, error) {
	var keyRing *crypto.KeyRing

	for _, pubKey := range e.publicKey {
		key, err := crypto.NewKeyFromArmored(pubKey)
		if err != nil {
			return nil, fmt.Errorf("invalid public key: %w", err)
		}

		if keyRing == nil {
			keyRing, err = crypto.NewKeyRing(key)
		} else {
			err = keyRing.AddKey(key)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to add key to key ring: %w", err)
		}
	}

	return keyRing, nil
}
