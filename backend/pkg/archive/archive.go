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
	"time"
)

// Archiver represents the archiving functionality with customizable options.
type Archiver struct{ *Config }

// truncateFile truncates the specified file to start fresh.
func (a *Archiver) truncateFile() error { return os.Truncate(a.Config.DocFile, 0) }

// File archives the specified document file by compressing it into a tar.gz archive.
// It creates a new archive file with a formatted filename based on the Archiver's fileNameFormat.
// The function supports streaming and truncating the file while other callers are writing to it.
//
// Note: For better performance, it depends on the disk. In Kubernetes, it also depends on the CSI driver storage.
// If using a CSI driver storage with S3 or S3-compatible storage in Kubernetes, it might require running the container (this deployment) as root,
// as it depends on the implementation of the CSI driver storage for S3 or S3-compatible storage.
// The compression ratio for tar.gz is such that if the file contains 10Gi, it will become approximately 1Gi in the tar.gz archive.
// However, when extracting the archive, it will revert back to its original size of 10Gi.
// Due to the streaming nature of the function, the speed depends on the disk. In production, archiving a 10Gi file typically takes only a few minute.
func (a *Archiver) File() (err error) {
	var timestamp string
	if a.Config.TimeFormat == defaultTimeFormat {
		timestamp = fmt.Sprintf(defaultTimeFormat, time.Now().Unix())
	} else {
		timestamp = time.Now().Format(a.Config.TimeFormat)
	}

	// Generate the archive filename based on the fileNameFormat.
	archiveFileName := fmt.Sprintf(a.Config.FileNameFormat+".tar.gz", filepath.Base(a.Config.DocFile), timestamp)
	archiveFilePath := filepath.Join(a.Config.ArchiveDir, archiveFileName)

	// Open the document file for reading.
	// This is the correct placement for both non-streaming and streaming.
	// The previous commit was only for streaming when writing data to disk.
	file, err := os.Open(a.Config.DocFile)
	if err != nil {
		return fmt.Errorf("error opening document file: %v", err)
	}
	defer func() {
		// In case an error occurs during file closure, this Trick Go deferred function
		// captures the error and assigns it to the named return value "err".
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("error closing document file: %v", closeErr)
		}
	}()

	// Get the file information of the document file.
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting document file info: %v", err)
	}

	// Create the archive directory if it doesn't exist.
	if err := os.MkdirAll(a.Config.ArchiveDir, os.ModePerm); err != nil {
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

	// Write the tar header for the document file.
	//
	// TODO: Do we need to consider adding more header fields for the tar archive?
	//       (e.g., Additional header fields could include ownership, permissions, or other metadata.)
	header := &tar.Header{
		Name:    filepath.Base(a.Config.DocFile),
		Mode:    0600,
		Size:    fileInfo.Size(),
		ModTime: time.Now(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("error writing tar header: %v", err)
	}

	// Create a channel to signal when streaming is done.
	//
	// TODO: This still needs improvement. Might Must using the Context Mechanism.
	done := make(chan error)

	// Start another goroutine to stream the document file contents to the tar writer.
	go func() {
		// Using a channel for signaling completion might be better than using a context in this case.
		// Due this requires another goroutine because it cannot directly interact with the caller's main goroutine in a non-blocking manner.
		_, err := io.Copy(tarWriter, file)
		done <- err
	}()

	// Wait for the streaming to complete.
	if err := <-done; err != nil {
		return fmt.Errorf("error writing document file to archive: %v", err)
	}

	// Truncate the document file to start fresh.
	if err := a.truncateFile(); err != nil {
		return fmt.Errorf("error truncating document file: %v", err)
	}

	return err
}
