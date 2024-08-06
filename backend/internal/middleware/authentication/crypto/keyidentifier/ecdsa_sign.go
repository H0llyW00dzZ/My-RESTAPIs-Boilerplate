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

// signUUIDWithECDSA signs the given UUID using ECDSA and returns the signature in ASN.1 DER format.
func (k *KeyIdentifier) signUUIDWithECDSA(uuid string) ([]byte, error) {
	// Check if the hash function is set in the configuration
	if k.config.Digest == nil {
		return nil, errors.New("crypto/keyidentifier: hash function is not set in the configuration")
	}

	digest := k.digest([]byte(uuid))

	// Sign the Digest using ECDSA and return the signature in ASN.1 DER format
	signature, err := ecdsa.SignASN1(k.secureRandom(), k.config.PrivateKey, digest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign UUID with ECDSA: %v", err)
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

// signUUIDWithHSM signs the given UUID using the HSM and returns the signature.
//
// Note: Testing this function is skipped due to the challenging of using an HSM in testing mode. However, it is recommended to use "signUUIDWithECDSA" instead.
// Regarding the private key of ECDSA, it can be easily maintained while using "signUUIDWithECDSA", especially in a Kubernetes environment that already has an HSM connector.
func (k *KeyIdentifier) signUUIDWithHSM(uuid string) ([]byte, error) {
	// Check if the hash function is set in the configuration
	if k.config.Digest == nil {
		return nil, errors.New("crypto/keyidentifier: hash function is not set in the configuration")
	}

	digest := k.digest([]byte(uuid))

	// Sign the digest using the HSM and return the signature
	signature, err := k.config.HSM.Sign(k.secureRandom(), digest, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to sign UUID with HSM: %v", err)
	}

	return signature, nil
}

// digest computes the hash of the given UUID using the configured hash function.
//
// It takes the following parameter:
//   - uuid: The UUID as a byte slice.
//
// It returns the computed hash of the UUID as a byte slice.
//
// This function is used internally by the [signUUIDWithECDSA] and [signUUIDWithHSM] functions
// to compute the hash of the UUID before signing it. The hash function used for the digest
// is determined by the "Digest" field in the [Config] struct.
//
// Example usage:
//
//	digest := k.digest([]byte(uuid))
//
// Note: The "Digest" field in the [Config] struct must be set to a valid hash function.
// If the [Digest] field is not set, an error will be returned by the calling function.
func (k *KeyIdentifier) digest(uuid []byte) []byte {
	// Digest the UUID using the configured hash function
	h := k.config.Digest()
	h.Write(uuid)
	return h.Sum(nil)
}
