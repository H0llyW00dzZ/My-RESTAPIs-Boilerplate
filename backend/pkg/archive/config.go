// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package archive

import (
	"time"
)

// Config represents the configuration for the archiving process.
type Config struct {
	DocFile        string        // Path to the document file
	ArchiveDir     string        // Directory to store the archived files
	MaxSize        int64         // Maximum size of the document file before archiving
	CheckInterval  time.Duration // Time interval for checking the document file size
	FileNameFormat string        // Format string for the archive filename
	TimeFormat     string        // Format string for the timestamp
}

const (
	defaultMaxSize        = 10 * 1024 * 1024 // 10 MiB
	defaultInterval       = 5 * time.Minute  // Check every 5 minutes
	defaultFileNameFormat = "%s_%s"          // Default archive filename format
	defaultTimeFormat     = "%d"             // Default to Unix timestamp format
)

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		MaxSize:        defaultMaxSize,
		CheckInterval:  defaultInterval,
		FileNameFormat: defaultFileNameFormat,
	}
}
