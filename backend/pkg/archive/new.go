// Copyright (c) 2025 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package archive

import "strings"

// Archiver represents the archiving functionality with customizable options.
type Archiver struct {
	*Config
}

// NewArchiver creates a new Archiver with the specified options.
// If no options are provided, the default values from the Config struct will be used.
func NewArchiver(configs ...Config) *Archiver {
	config := DefaultConfig()
	if len(configs) > 0 {
		config = configs[0]
	}

	if config.TimeFormat == "" {
		config.TimeFormat = "%d" // Default to Unix timestamp format
	}

	// Note: If you're not familiar with the formatting logic, let me explain how it works for the file name and time.
	// The file name format uses two placeholders: %s for the base file name and %s for the timestamp.
	// For example, if the base file name is "example.log" and the timestamp is "20230608123456",
	// the resulting file name will be "example.log_20230608123456.tar.gz".
	if config.FileNameFormat == "" || strings.Count(config.FileNameFormat, "%s") < 2 {
		config.FileNameFormat = defaultFileNameFormat // Default to %s_%s.tar.gz
	}

	return &Archiver{
		Config: &config,
	}
}
