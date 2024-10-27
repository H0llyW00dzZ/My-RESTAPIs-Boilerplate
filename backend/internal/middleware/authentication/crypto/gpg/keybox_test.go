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
	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := kb.AddKey(testPublicKey); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kb.KeyCount() != 1 {
		t.Fatalf("expected 1 key, got %d", kb.KeyCount())
	}
}

func TestKeybox_SaveAndLoad(t *testing.T) {
	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := kb.AddKey(testPublicKey); err != nil {
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
	t.Logf("JSON output: %s", buffer.String())

	loadedKb, err := gpg.Load(strings.NewReader(buffer.String()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loadedKb.KeyCount() != 1 {
		t.Fatalf("expected 1 key, got %d", loadedKb.KeyCount())
	}
}

func TestKeybox_GetEncryptor(t *testing.T) {
	kb, err := gpg.NewKeybox()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := kb.AddKey(testPublicKey); err != nil {
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
