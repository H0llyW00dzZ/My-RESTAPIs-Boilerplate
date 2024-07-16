// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package signature_test

import (
	"encoding/hex"
	"os"
	"testing"

	"h0llyw00dz-template/backend/internal/middleware/filesystem/crypto/signature"
)

func TestGenerateAndVerifyHMACSignature(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some content to the temporary file
	fileContent := []byte("This is a test file.")
	if _, err := tempFile.Write(fileContent); err != nil {
		t.Fatalf("Failed to write to temp file: %s", err)
	}
	if err := tempFile.Sync(); err != nil {
		t.Fatalf("Failed to sync file content: %s", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %s", err)
	}

	// Set the secret key
	secretKey := "gopher-testing-testing-testing"

	// Generate the HMAC signature from the file
	sig, err := signature.GenerateHMACSignatureFromFile(tempFile.Name(), secretKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC signature: %s", err)
	}

	// Verify the generated signature
	isValid, err := signature.VerifyHMACSignatureFromFile(tempFile.Name(), secretKey, hex.EncodeToString(sig))
	if err != nil {
		t.Fatalf("Failed to verify HMAC signature: %s", err)
	}
	if !isValid {
		t.Errorf("Expected signature to be valid, but it was invalid. Signature: %s", hex.EncodeToString(sig))
	}

	// Modify the file content and write it to the file
	modifiedContent := []byte("Modified content")
	if err := os.WriteFile(tempFile.Name(), modifiedContent, 0644); err != nil {
		t.Fatalf("Failed to write modified content to temp file: %s", err)
	}

	// Generate a new signature for the modified content
	modifiedSig, err := signature.GenerateHMACSignatureFromFile(tempFile.Name(), secretKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC signature for modified content: %s", err)
	}

	// Verify the signature for the modified content
	isValid, err = signature.VerifyHMACSignatureFromFile(tempFile.Name(), secretKey, hex.EncodeToString(modifiedSig))
	if err != nil {
		t.Fatalf("Failed to verify HMAC signature for modified content: %s", err)
	}
	if !isValid {
		t.Errorf("Expected signature for modified content to be valid, but it was invalid. Signature: %s", hex.EncodeToString(modifiedSig))
	}

	// Output signatures for debugging
	t.Logf("Original signature: %s", sig)
	t.Logf("Modified signature: %s", modifiedSig)
}

func TestGenerateAndVerifyHMACWithModifiedContent(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some content to the temporary file
	fileContent := []byte("This is a test file.")
	if _, err := tempFile.Write(fileContent); err != nil {
		t.Fatalf("Failed to write to temp file: %s", err)
	}
	if err := tempFile.Sync(); err != nil {
		t.Fatalf("Failed to sync file content: %s", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %s", err)
	}

	// Set the secret key
	secretKey := "gopher-testing-testing-testing"

	// Generate the HMAC signature from the file
	sig, err := signature.GenerateHMACSignatureFromFile(tempFile.Name(), secretKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC signature: %s", err)
	}
	hexSig := hex.EncodeToString(sig)

	// Modify the file content and write it to the file
	modifiedContent := []byte("Modified content")
	if err := os.WriteFile(tempFile.Name(), modifiedContent, 0644); err != nil {
		t.Fatalf("Failed to write modified content to temp file: %s", err)
	}

	// Attempt to verify the original signature against the modified file content
	isValid, err := signature.VerifyHMACSignatureFromFile(tempFile.Name(), secretKey, hexSig)
	if err != nil {
		t.Fatalf("Unexpected error verifying signature: %s", err)
	}
	if isValid {
		t.Errorf("Expected signature to be invalid due to modified content, but it was valid. Signature: %s", hexSig)
	}
}
