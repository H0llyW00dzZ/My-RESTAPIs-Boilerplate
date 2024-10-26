// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// Config holds configuration options for encryption and other operations (TODO).
type Config struct {
	AllowVerfy bool
	compress   bool
	isBinary   bool
	modTime    int64
}

// NewDefaultConfig creates a default configuration.
func NewDefaultConfig() *Config {
	return &Config{
		AllowVerfy: false,
		compress:   true,
		isBinary:   true,
		modTime:    crypto.GetUnixTime(),
	}
}

// Option is a function that modifies the Config.
type Option func(*Config)

// WithBinary sets the IsBinary option.
func WithBinary(isBinary bool) Option { return func(c *Config) { c.isBinary = isBinary } }

// WithModTime sets the ModTime option.
func WithModTime(modTime int64) Option { return func(c *Config) { c.modTime = modTime } }

// WithAllowVerify allows keys that cannot be used for encryption to be stored in the keyring for future verification (TODO).
func WithAllowVerify(allowVerfy bool) Option { return func(c *Config) { c.AllowVerfy = allowVerfy } }

// WithCompress sets the option to use compression during encryption.
func WithCompress(compress bool) Option { return func(c *Config) { c.compress = compress } }
