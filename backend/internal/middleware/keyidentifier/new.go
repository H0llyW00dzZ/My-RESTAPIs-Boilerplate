// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

// Config represents the configuration options for the key identifier.
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
