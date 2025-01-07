// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package archive

import (
	"fmt"
	"os"
	"time"
)

// Do monitors the specified document file and archives it when its size reaches the maximum size.
// It continuously checks the document file size based on the configured check interval.
//
// Example Usage:
//
//	// Configure the archiving process
//	configArchive := archive.Config{
//		MaxSize:       int64(sizeBytes),
//		CheckInterval: sizeInterval,
//	}
//
//	go func() {
//		if err := archive.Do(diskStorageFiberLog, diskStorageFiberLogArchiveDir, configArchive); err != nil {
//			// Handle error you poggers
//		}
//	}()
//
// Note: Due to the streaming nature of the archiving process and its ability to connect with code that writes files using streaming mechanisms,
// ensure that the document or log file is properly closed when it is no longer writable. Additionally, this archiving mechanism has zero memory allocation overhead
// because it utilizes streaming methods.
func Do(docFile, archiveDir string, configs ...Config) error {
	config := DefaultConfig()
	if len(configs) > 0 {
		config = configs[0]
	}
	config.docFile = docFile
	config.archiveDir = archiveDir

	archiver := NewArchiver(config)

	for {
		// Get the file information of the document file.
		fileInfo, err := os.Stat(config.docFile)
		if err != nil {
			time.Sleep(config.CheckInterval)
			return fmt.Errorf("error checking document file: %v", err)
		}

		// Check if the document file size exceeds the maximum size.
		if fileInfo.Size() >= config.MaxSize {
			if err := archiver.ArchiveDoc(); err != nil {
				return err
			}
		}

		// Sleep for the configured check interval before the next check.
		time.Sleep(config.CheckInterval)
	}
}
