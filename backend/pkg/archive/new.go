// Copyright (c) 2025 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package archive

// Archiver represents the archiving functionality with customizable options.
type Archiver struct {
	fileNameFormat string // Format string for the archive filename
	timeFormat     string // Format string for the timestamp
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

	return &Archiver{
		fileNameFormat: config.FileNameFormat,
		timeFormat:     config.TimeFormat,
	}
}
