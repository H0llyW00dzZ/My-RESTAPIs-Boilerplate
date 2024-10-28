// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg_test

import (
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/gpg"
	"strings"
	"testing"
)

func TestKeybox_AddKey(t *testing.T) {
	// Sample public key
	publicKeys := []string{
		testPublicKey,
		// test duplicate public key
		testPublicKey,
	}

	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := kb.AddKey(publicKeys); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kb.KeyCount() != 1 {
		t.Fatalf("expected 1 key, got %d", kb.KeyCount())
	}
}

func TestKeybox_SaveAndLoad(t *testing.T) {
	// Sample public key
	publicKeys := []string{
		// Support multiple public key
		testPublicKey,
		testPublicKeyRSA2048,
	}

	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add multiple keys to the Keybox
	if err := kb.AddKey(publicKeys); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// This is compatible with most standard libraries due to its I/O operations.
	var buffer strings.Builder
	if err := kb.Save(&buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Output the JSON format for visual inspection.
	//
	// This format will not be corrupted and provides better readability when there are many keys.
	// For example (JSON):
	// {
	//	"uuid": "1ee90424-2892-4df4-bad8-522a5a5dade6",
	// 	"keys": [
	// 	  {
	// 		"fingerprint": "ABC123DEF456...",
	// 		"creation_date": "2023-10-28T00:00:00Z",
	// 		"armored_key": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...\n-----END PGP PUBLIC KEY BLOCK-----"
	// 	  },
	// 	  {
	// 		"fingerprint": "XYZ789GHI012...",
	// 		"creation_date": "2024-01-15T00:00:00Z",
	// 		"armored_key": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...\n-----END PGP PUBLIC KEY BLOCK-----"
	// 	  }
	// 	]
	// }
	//
	// Also the UUID can later be used as an identifier or for other purposes as needed.
	// It is secure, and even if someone tries to guess it, it remains a random number.
	t.Logf("JSON output: %s", buffer.String())

	loadedKb, err := gpg.Load(strings.NewReader(buffer.String()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loadedKb.KeyCount() != 2 {
		t.Fatalf("expected 2 key, got %d", loadedKb.KeyCount())
	}
}

func TestKeybox_GetEncryptor(t *testing.T) {
	// Sample public key
	publicKeys := []string{
		testPublicKey,
	}

	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := kb.AddKey(publicKeys); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	encryptor, err := kb.GetEncryptor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if encryptor == nil {
		t.Fatal("expected encryptor, got nil")
	}
}

func TestKeybox_GetEncryptor_NoKeys(t *testing.T) {
	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = kb.GetEncryptor() // Remove redeclaration of err
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != gpg.ErrorCantEncrypt {
		t.Fatalf("expected ErrorCantEncrypt, got %v", err)
	}
}

func TestKeybox_EncryptBeforeSave(t *testing.T) {
	// Sample public key
	publicKeys := []string{
		// Support multiple public key
		testPublicKey,
		testPublicKeyRSA2048,
	}

	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add multiple keys to the Keybox
	if err := kb.AddKey(publicKeys); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	encryptor, err := kb.GetEncryptor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Use a buffer to simulate file writing
	var buffer strings.Builder

	// Encrypt and save the Keybox
	if err := kb.EncryptBeforeSave(&buffer, encryptor); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Log the encrypted JSON output for inspection
	//
	// This format will not be corrupted and provides better readability when there are many keys.
	// For example (JSON):
	// {
	// 	"uuid": "1ee90424-2892-4df4-bad8-522a5a5dade6",
	// 	"keys": [
	// 	  {
	// 		"encrypted": "-----BEGIN PGP MESSAGE-----\n...\n-----END PGP MESSAGE-----"
	// 	  },
	// 	  {
	// 		"encrypted": "-----BEGIN PGP MESSAGE-----\n...\n-----END PGP MESSAGE-----"
	// 	  }
	// 	]
	// }
	// Also the UUID can later be used as an identifier or for other purposes as needed.
	// It is secure, and even if someone tries to guess it, it remains a random number.
	t.Logf("Encrypted JSON output: %s", buffer.String())

	// Load the Keybox from the buffer to ensure it is correctly saved
	loadedKb, err := gpg.Load(strings.NewReader(buffer.String()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if the loaded Keybox has the same number of keys
	if loadedKb.KeyCount() != 2 {
		t.Fatalf("expected 2 key, got %d", loadedKb.KeyCount())
	}
}
