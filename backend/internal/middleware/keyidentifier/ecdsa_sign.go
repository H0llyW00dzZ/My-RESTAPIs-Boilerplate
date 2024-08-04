// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
)

// signUUID signs the given UUID using ECDSA and returns the signature in ASN.1 DER format.
func (k *KeyIdentifier) signUUID(uuid string) ([]byte, error) {
	// Check if the private key is set in the configuration
	if k.config.PrivateKey == nil {
		return nil, fmt.Errorf("private key is not set in the configuration")
	}

	// Check if the hash function is set in the configuration
	if k.config.Digest == nil {
		return nil, fmt.Errorf("hash function is not set in the configuration")
	}

	// Digest the UUID using the configured hash function
	h := k.config.Digest()
	h.Write([]byte(uuid))
	digest := h.Sum(nil)

	// Sign the Digest using ECDSA and return the signature in ASN.1 DER format
	signature, err := ecdsa.SignASN1(rand.Reader, k.config.PrivateKey, digest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign UUID: %v", err)
	}

	return signature, nil
}
