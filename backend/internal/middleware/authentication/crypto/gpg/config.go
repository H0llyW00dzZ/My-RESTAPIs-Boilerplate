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
	armor      bool
	suffix     string
}

// NewDefaultConfig creates a default configuration.
func NewDefaultConfig() *Config {
	return &Config{
		AllowVerfy: false,
		compress:   true,
		isBinary:   true,
		modTime:    crypto.GetUnixTime(),
		armor:      false,
		suffix:     newGPGModern,
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

// WithArmor sets the option to armor the encrypted message.
//
// Note: Default is false to minimize memory allocation, which is effective for on-the-fly encryption (when set false).
//
// Learn more about how on-the-fly encryption works at: https://www.east-tec.com/kb/safebit/protecting-your-confidential-information/what-does-on-the-fly-encryption-mean/
func WithArmor(armor bool) Option { return func(c *Config) { c.armor = armor } }

// WithCustomSuffix sets a custom suffix for the output filename.
//
// This option is effective when armor is enabled and the input is not a file.
// It allows you to specify a suffix other than the default ".gpg".
// Ensure the suffix is not empty and differs from the default to apply custom behavior
// during filename extraction. The suffix will be used if the output file has an extension
// matching the custom suffix.
func WithCustomSuffix(suffix string) Option { return func(c *Config) { c.suffix = suffix } }
