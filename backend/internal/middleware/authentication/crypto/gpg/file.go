// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// extractFilename checks if the input or output is a file and extracts the filename.
//
// This effectively Go detects actual files (file disk not a memory again) in I/O.
//
// # Result:
//
//	Decrypt Operation - Success
//
// # General State:
//   - File Name: test_output_1960559248.txt
//   - MIME: false
//   - Message Integrity Protection: true
//   - Symmetric Encryption Algorithm: AES256.CFB
//   - German Encryption Standards: false
//
// Note: This helper function uses [os.File], which connects to the filesystem.
// If files are handled differently (other way), they may reside entirely in memory and not actual on disk.
func extractFilename(i io.Reader, o io.Writer, suffix string) (string, error) {
	// Check if the input is a file.
	if file, ok := i.(*os.File); ok {
		return filepath.Base(file.Name()), nil
	}

	// If the input is not a file, check if the output is a file.
	if outFile, ok := o.(*os.File); ok {
		outName := filepath.Base(outFile.Name())
		if strings.HasSuffix(outName, suffix) {
			return strings.TrimSuffix(outName, suffix), nil
		}
		return "", fmt.Errorf("output file does not have the expected %s suffix", suffix)
	}

	// If neither is a file, return an empty filename.
	return "", nil
}
