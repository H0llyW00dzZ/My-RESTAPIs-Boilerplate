// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg_test

import (
	"bytes"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/gpg"
	"testing"
)

func TestKeybox_AddKey(t *testing.T) {
	kb := gpg.NewKeybox()

	err := kb.AddKey(testPublicKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kb.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(kb.Keys))
	}
}

func TestKeybox_SaveAndLoad(t *testing.T) {
	kb := gpg.NewKeybox()
	err := kb.AddKey(testPublicKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buffer bytes.Buffer
	err = kb.Save(&buffer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Output the JSON format for visual inspection.
	//
	// This format will not be corrupted and provides better readability when there are many keys.
	// For example:
	// [
	// "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...\n-----END PGP PUBLIC KEY BLOCK-----",
	// "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...\n-----END PGP PUBLIC KEY BLOCK-----"
	// ]
	t.Logf("JSON output: %s", buffer.String())

	loadedKb, err := gpg.Load(&buffer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loadedKb.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(loadedKb.Keys))
	}
}

func TestKeybox_GetEncryptor(t *testing.T) {
	kb := gpg.NewKeybox()
	err := kb.AddKey(testPublicKey)
	if err != nil {
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
	kb := gpg.NewKeybox()

	_, err := kb.GetEncryptor()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != gpg.ErrorCantEncrypt {
		t.Fatalf("expected ErrorCantEncrypt, got %v", err)
	}
}
