// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package archive_test

import (
	"h0llyw00dz-template/backend/pkg/archive"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMonitorAndArchiveLog(t *testing.T) {
	logFile := "./test.log"
	archiveDir := "./archives"
	maxSize := int64(1024) // 1KiB for testing purposes

	// Create a sample log file
	err := createSampleLogFile(logFile, 2048) // 2KiB log file
	if err != nil {
		t.Fatalf("error creating sample log file: %v", err)
	}
	defer os.Remove(logFile)

	// Configure the archiving process
	config := archive.Config{
		MaxSize:       maxSize,
		CheckInterval: time.Second * 1, // Check every second for testing purposes
	}

	// Start monitoring and archiving the log file
	go archive.Do(logFile, archiveDir, config)

	// Wait for a short duration to allow archiving to happen
	time.Sleep(time.Second * 2)

	// Check if the log file was archived
	archiveFile := filepath.Join(archiveDir, "test.log_*.tar.gz")
	matches, err := filepath.Glob(archiveFile)
	if err != nil {
		t.Fatalf("error finding archive file: %v", err)
	}

	if len(matches) == 0 {
		t.Error("log file was not archived")
	}

	// Clean up the archive directory
	os.RemoveAll(archiveDir)
}

func createSampleLogFile(logFile string, size int64) error {
	file, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(make([]byte, size)); err != nil {
		return err
	}

	return nil
}
