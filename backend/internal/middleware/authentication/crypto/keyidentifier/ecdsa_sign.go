// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io"

	"errors"
)

// signUUID signs the given UUID using ECDSA and returns the signature in ASN.1 DER format.
func (k *KeyIdentifier) signUUID(uuid string) ([]byte, error) {
	// Check if the hash function is set in the configuration
	if k.config.Digest == nil {
		return nil, errors.New("crypto/keyidentifier: hash function is not set in the configuration")
	}

	// Digest the UUID using the configured hash function
	h := k.config.Digest()
	h.Write([]byte(uuid))
	digest := h.Sum(nil)

	// Sign the Digest using ECDSA and return the signature in ASN.1 DER format
	signature, err := ecdsa.SignASN1(k.secureRandom(), k.config.PrivateKey, digest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign UUID: %v", err)
	}

	return signature, nil
}

// secureRandom returns a secure random number generator.
//
// If a custom random number generator is provided in the configuration [Config.Rand],
// it will be used. Otherwise, the default [crypto/rand] is used.
//
// The purpose of this function is to allow flexibility in the source of randomness used
// for cryptographic operations. By default, it uses the [crypto/rand] package, which provides
// a cryptographically secure random number generator. However, if a custom random number
// generator is needed for specific requirements (advanced use cases) or testing purposes, it can be set in the
// configuration.
func (k *KeyIdentifier) secureRandom() io.Reader {
	// If no custom random number generator is provided in the configuration,
	if k.config.Rand == nil {
		// use the default [crypto/rand] as the secure random number generator.
		return rand.Reader
	}

	// If a custom random number generator is provided in the configuration,
	// return it as the secure random number generator.
	return k.config.Rand
}
