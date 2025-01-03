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
	docFile       string        // Path to the document file
	archiveDir    string        // Directory to store the archived files
	MaxSize       int64         // Maximum size of the document file before archiving
	CheckInterval time.Duration // Time interval for checking the document file size
}

const (
	defaultMaxSize  = 10 * 1024 * 1024 // 10 MiB
	defaultInterval = 5 * time.Minute  // Check every 5 minutes
)

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		MaxSize:       defaultMaxSize,
		CheckInterval: defaultInterval,
	}
}
