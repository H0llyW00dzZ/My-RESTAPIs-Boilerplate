// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

// Config represents the configuration options for the key identifier.
//
// Note: The Prefix here is not actually a key, it's a group-key. For example, "session_id_authorized:<uuid>",
// where <uuid> is the actual key to get the value. This is because memory storage is unstructured, unlike
// relational databases that use queries and tables.
type Config struct {
	Prefix string
}

// ConfigDefault is the default configuration for the key identifier.
var ConfigDefault = Config{
	Prefix: "session_id_authorized:",
}

// KeyIdentifier represents the key identifier.
type KeyIdentifier struct {
	config Config
}

// New creates a new instance of the key identifier with the given configuration.
func New(config ...Config) *KeyIdentifier {
	cfg := ConfigDefault

	if len(config) > 0 {
		cfg = config[0]
	}

	return &KeyIdentifier{
		config: cfg,
	}
}
