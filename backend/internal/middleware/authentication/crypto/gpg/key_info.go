// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

// KeyInfo holds metadata about a public key.
type KeyInfo struct {
	UserIDs           []string
	PrimaryKey        string
	KeyID             uint64
	HexKeyID          string
	CanEncrypt        bool
	CanVerify         bool
	IsExpired         bool
	IsRevoked         bool
	Fingerprint       string
	DigestFingerprint []string
}

// GetKeyInfos returns a slice of KeyInfo structs containing metadata
// about all the public keys managed by the Encryptor.
func (e *Encryptor) GetKeyInfos() []KeyInfo { return e.keyInfos }
