// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package archive_test

import (
	"archive/tar"
	"compress/gzip"
	"h0llyw00dz-template/backend/pkg/archive"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestMonitorAndArchiveLog(t *testing.T) {
	logFile := "./test.log"
	archiveDir := "./archives"
	maxSize := int64(1024) // 1KiB for testing purposes

	// Create a sample log file
	if err := createSampleLogFile(logFile, 2048); err != nil { // 2KiB log file
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

var initArchive sync.Once

func TestArchiveOnce(t *testing.T) {
	logFile := "./test.log"
	archiveDir := "./archives-once"
	maxSize := int64(1024) // 1KiB for testing purposes

	// Create a sample log file
	if err := createSampleLogFile(logFile, 2048); err != nil { // 2KiB log file
		t.Fatalf("error creating sample log file: %v", err)
	}
	defer os.Remove(logFile)

	// Configure the archiving process
	config := archive.Config{
		MaxSize:        maxSize,
		CheckInterval:  time.Second * 1, // Check every second for testing purposes
		FileNameFormat: "%s_%s",
		TimeFormat:     "20060102150405",
	}

	// Start monitoring and archiving the log file using [sync.Once]
	initArchive.Do(func() {
		go archive.Do(logFile, archiveDir, config)
	})

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

func getArchiveFileSize(file string) (int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		return 0, err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	header, err := tarReader.Next()
	if err != nil {
		return 0, err
	}

	return header.Size, nil
}

// Note: This test is a simulation because testing goroutine synchronization mechanism is challenging, unlike in production.
func TestStreamingAndMultipleArchiving(t *testing.T) {
	logFile := "./test.log"
	archiveDir := "./multiple-archives"
	maxSize := int64(1024) // 1KiB for testing purposes

	// Create a sample log file
	if err := createSampleLogFile(logFile, 512); err != nil { // 0.5KiB log file
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

	// Start streaming data to the log file
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Errorf("error opening log file for writing: %v", err)
			return
		}
		defer file.Close()

		// Write data to the log file in chunks
		chunkSize := int64(512)
		data := make([]byte, chunkSize)
		for i := 0; i < 5; i++ {
			if _, err := file.Write(data); err != nil {
				t.Errorf("error writing to log file: %v", err)
				return
			}
			time.Sleep(time.Millisecond * 500) // Simulate a delay between writes
		}
	}()

	// Wait for the streaming to finish
	wg.Wait()

	// Wait for a short duration to allow archiving to happen
	time.Sleep(time.Second * 3)

	// Check if the log file was archived twice
	archiveFile := filepath.Join(archiveDir, "test.log_*.tar.gz")
	matches, err := filepath.Glob(archiveFile)
	if err != nil {
		t.Fatalf("error finding archive files: %v", err)
	}

	if len(matches) != 3 {
		t.Errorf("expected 3 archive files, got %d", len(matches))
	}

	// Check the contents of the archive files
	var totalSize int64
	for _, file := range matches {
		size, err := getArchiveFileSize(file)
		if err != nil {
			t.Errorf("error getting archive file size: %v", err)
		}
		totalSize += size
	}

	expectedSize := int64(3072) // 5 chunks of 512 bytes each
	if totalSize != expectedSize {
		t.Errorf("expected total archive size %d, got %d", expectedSize, totalSize)
	}

	// Clean up the archive directory
	os.RemoveAll(archiveDir)
}

func TestArchiveError(t *testing.T) {
	tests := []struct {
		name           string
		logFile        string
		archiveDir     string
		maxSize        int64
		docFile        string
		expectedErrMsg string
	}{
		{
			name:           "Empty DocFile",
			logFile:        "./test.log",
			archiveDir:     "./archives-error",
			maxSize:        1024,
			docFile:        "",
			expectedErrMsg: "expected an error when DocFile is empty",
		},
		{
			name:           "Empty ArchiveDir",
			logFile:        "./test.log",
			archiveDir:     "",
			maxSize:        1024,
			docFile:        "./test.log",
			expectedErrMsg: "expected an error when ArchiveDir is empty",
		},
		{
			name:           "Empty DocFile & ArchiveDir",
			logFile:        "./test.log",
			archiveDir:     "",
			maxSize:        1024,
			docFile:        "",
			expectedErrMsg: "expected an error when DocFile & ArchiveDir is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a sample log file
			if err := createSampleLogFile(tt.logFile, 2048); err != nil { // 2KiB log file
				t.Fatalf("error creating sample log file: %v", err)
			}
			defer os.Remove(tt.logFile)

			// Configure the archiving process
			config := archive.Config{
				MaxSize:        tt.maxSize,
				CheckInterval:  time.Second * 1, // Check every second for testing purposes
				FileNameFormat: "%s_%s",
				TimeFormat:     "%d",
			}

			// Perform the test
			err := archive.Do(tt.docFile, tt.archiveDir, config)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Error(tt.expectedErrMsg)
				}
			}

			// Clean up the archive directory
			os.RemoveAll(tt.archiveDir)
		})
	}
}
