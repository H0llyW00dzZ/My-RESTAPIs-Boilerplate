// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

import (
	"crypto"
	"crypto/ecdsa"
	"hash"
	"io"
)

// Config represents the configuration options for the key identifier.
//
// Note: The Prefix here is not actually a key, it's a group-key. For example, "session_id_authorized:<uuid>",
// where <uuid> is the actual key to get the value. This is because memory storage is unstructured, unlike
// relational databases that use queries and tables.
// Also note that when you see logs in the redis/valkey or redis/valkey commander panel, "session_id_authorized:" will be categorized as a group,
// and <uuid> will be the key to get the value.
// Then To create multiple group-keys, similar to a binary tree (see https://en.wikipedia.org/wiki/Binary_tree), simply add another prefix tag. For example, "authorization:session_id_authorized:<uuid>".
type Config struct {
	Prefix           string
	PrivateKey       *ecdsa.PrivateKey
	Digest           func() hash.Hash
	SignedContextKey any
	Rand             io.Reader
	HSM              crypto.Signer
}

// ConfigDefault is the default configuration for the key identifier.
var ConfigDefault = Config{
	Prefix:           "session_id_authorized:",
	PrivateKey:       nil,
	Digest:           nil,
	SignedContextKey: nil,
	Rand:             nil,
	HSM:              nil,
}

// KeyIdentifier represents the key identifier.
type KeyIdentifier struct {
	config Config
}

// New creates a new instance of the key identifier with the given configuration.
//
// Note: It is recommended not to implement this function as a global variable across the codebase,
// as it may result in the same UUID being generated for all instances, rather than unique UUIDs.
//
// Example Usage:
//
//	func generateUUID() {
//		// Create a new key identifier with custom configuration
//		uuid := keyidentifier.New(keyidentifier.Config{
//			Prefix: "custom_prefix:",
//			// Set other configuration options as needed
//		})
//		// Use the key identifier
//		// ...
//	}
func New(config ...Config) *KeyIdentifier {
	cfg := ConfigDefault

	if len(config) > 0 {
		cfg = config[0]
	}

	return &KeyIdentifier{
		config: cfg,
	}
}
