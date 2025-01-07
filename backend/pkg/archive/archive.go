// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// truncateFile truncates the specified file to start fresh.
func truncateFile(file string) error { return os.Truncate(file, 0) }

// archiveDoc archives the specified document file by compressing it into a tar.gz archive.
// It creates a new archive file with a timestamp appended to the filename.
// The function supports streaming and truncating the file while other callers are writing to it.
func archiveDoc(docFile, archiveDir string) (err error) {
	// Generate the archive filename with a timestamp.
	archiveFileName := fmt.Sprintf("%s_%s.tar.gz", filepath.Base(docFile), time.Now().Format("20060102150405"))
	archiveFilePath := filepath.Join(archiveDir, archiveFileName)

	// Create the archive directory if it doesn't exist.
	if err := os.MkdirAll(archiveDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating archive directory: %v", err)
	}

	// Create the archive file for writing.
	archive, err := os.Create(archiveFilePath)
	if err != nil {
		return fmt.Errorf("error creating archive file: %v", err)
	}
	defer archive.Close()

	// Create a gzip writer to compress the archive.
	gzipWriter := gzip.NewWriter(archive)
	defer gzipWriter.Close()

	// Create a tar writer to write the document file to the archive.
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Open the document file for reading.
	file, err := os.Open(docFile)
	if err != nil {
		return fmt.Errorf("error opening document file: %v", err)
	}
	defer file.Close()

	// Get the file information of the document file.
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting document file info: %v", err)
	}

	// Write the tar header for the document file.
	//
	// TODO: Do we need to consider adding more header fields for the tar archive?
	//       (e.g., Additional header fields could include ownership, permissions, or other metadata.)
	header := &tar.Header{
		Name:    filepath.Base(docFile),
		Mode:    0600,
		Size:    fileInfo.Size(),
		ModTime: time.Now(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("error writing tar header: %v", err)
	}

	// Create a wait group to synchronize the streaming and truncation.
	//
	// TODO: This still needs improvement. Might Must using the Context Mechanism.
	var wg sync.WaitGroup
	wg.Add(2)

	// Start a goroutine to stream the document file contents to the tar writer.
	go func() {
		defer wg.Done()
		if _, err := io.Copy(tarWriter, file); err != nil {
			err = fmt.Errorf("error writing document file to archive: %v", err)
		}
	}()

	// Start a goroutine to truncate the document file after reaching the maximum size.
	go func() {
		defer wg.Done()
		// Truncate the document file to start fresh.
		if err := truncateFile(docFile); err != nil {
			err = fmt.Errorf("error truncating document file: %v", err)
		}
	}()

	// Wait for both goroutines to finish.
	wg.Wait()

	return err
}
