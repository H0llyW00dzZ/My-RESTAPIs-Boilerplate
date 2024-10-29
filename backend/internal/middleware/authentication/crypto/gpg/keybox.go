// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	"io"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/bytedance/sonic"
)

// KeyMetadata contains metadata about a GPG/OpenPGP key, including its fingerprint, creation date, and armored representation.
//
// TODO: Adding support for YAML format ?
type KeyMetadata struct {
	Fingerprint  string `json:"fingerprint"`
	CreationDate string `json:"creation_date"`
	ArmoredKey   string `json:"armored_key"`
}

// Keybox manages a collection of keys that can be stored and retrieved.
type Keybox struct {
	UUID      string        `json:"uuid"`
	TotalKeys int           `json:"total_keys"`
	Keys      []KeyMetadata `json:"keys"`
}

// KeyMetadataEncrypted contains metadata about a GPG/OpenPGP key, including its encrypted representation.
type KeyMetadataEncrypted struct {
	Encrypted string `json:"encrypted"`
}

// NewKeybox creates a new Keybox instance.
func NewKeybox() (*Keybox, error) {
	uuid, err := rand.GenerateFixedUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}
	return &Keybox{
		UUID: uuid,
		Keys: []KeyMetadata{},
	}, nil
}

// AddKey adds a new key to the Keybox, supporting multiple purposes.
func (kb *Keybox) AddKey(armoredKey []string) error {
	// Set of existing fingerprints to avoid duplicates.
	existingFingerprints := make(map[string]bool)
	for _, keyInfo := range kb.Keys {
		existingFingerprints[keyInfo.Fingerprint] = true
	}

	for _, armoredKeys := range armoredKey {
		key, err := crypto.NewKeyFromArmored(armoredKeys)
		if err != nil {
			return fmt.Errorf("invalid key: %w", err)
		}

		fingerprint := key.GetFingerprint()
		if existingFingerprints[fingerprint] {
			continue // Skip duplicate keys
		}

		creationDate := key.GetEntity().PrimaryKey.CreationTime.UTC().Format(time.RFC3339)

		armoredWithHeader, err := kb.armorKeyWithHeader(*key)
		if err != nil {
			return fmt.Errorf("failed to add custom header: %w", err)
		}

		keyInfo := KeyMetadata{
			Fingerprint:  key.GetFingerprint(),
			CreationDate: creationDate,
			ArmoredKey:   armoredWithHeader,
		}

		existingFingerprints[fingerprint] = true
		kb.Keys = append(kb.Keys, keyInfo)
		kb.TotalKeys++ // Increment the total key count
	}
	return nil
}

// Save saves the Keybox to an [io.Writer] in JSON format.
//
// Note: Since it supports multiple purposes, it's recommended to store it in a file, network storage, or a database that can handle this object.
// Avoid using GPG key handling mechanisms that store keys directly in memory, as it is inefficient for a large number of keys.
//
// For example, handling many keys or UIDs (e.g., revoked UIDs) in armored format can be inefficient
// because they still require processing during export/import.
func (kb *Keybox) Save(o io.Writer) error {
	// Now we can perform this operation over the network, especially when using Kubernetes. It's very smooth sailing.
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		data, err := sonic.MarshalIndent(kb, "", "  ")
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to marshal keybox: %w", err))
			return
		}

		if _, err := pw.Write(data); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to write keybox: %w", err))
			return
		}
	}()

	if _, err := io.Copy(o, pr); err != nil {
		return fmt.Errorf("failed to copy data to writer: %w", err)
	}

	return nil
}

// Load loads a Keybox from an [io.Reader] in JSON format.
//
// Note: Since it supports multiple purposes, it's recommended to store it in a file, network storage, or a database that can handle this object.
// Avoid using GPG key handling mechanisms that store keys directly in memory, as it is inefficient for a large number of keys.
//
// For example, handling many keys or UIDs (e.g., revoked UIDs) in armored format can be inefficient
// because they still require processing during export/import.
func Load(i io.Reader) (*Keybox, error) {
	// Now we can perform this operation over the network, especially when using Kubernetes. It's very smooth sailing.
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		if _, err := io.Copy(pw, i); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to copy data from reader: %w", err))
			return
		}
	}()

	var kb Keybox
	if err := sonic.ConfigDefault.NewDecoder(pr).Decode(&kb); err != nil {
		return nil, fmt.Errorf("failed to decode keybox: %w", err)
	}

	return &kb, nil
}

// GetEncryptor creates an Encryptor from the keys in the Keybox that can be used for encryption.
func (kb *Keybox) GetEncryptor() (*Encryptor, error) {
	var encryptKeys []string
	for _, keyInfo := range kb.Keys {
		key, err := crypto.NewKeyFromArmored(keyInfo.ArmoredKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse armored key: %w", err)
		}
		if key.CanEncrypt() {
			encryptKeys = append(encryptKeys, keyInfo.ArmoredKey)
		}
	}

	if len(encryptKeys) == 0 {
		return nil, ErrorCantEncrypt
	}

	return NewEncryptor(encryptKeys)
}

// KeyCount returns the number of keys in the Keybox.
//
// This method provides a safe and easy way to get the count.
func (kb *Keybox) KeyCount() int { return len(kb.Keys) }

// TODO: Implement automated version detection for improved versioning
const keyBoxVersion = "v0.0.0-beta"
const customHeader = "From GPG/OpenPGP Keybox ⛵ 📦 🔐 🗝️  Written In Go by H0llyW00dzZ"

func (kb *Keybox) armorKeyWithHeader(key crypto.Key) (string, error) {
	armored, err := key.ArmorWithCustomHeaders(customHeader, keyBoxVersion)
	if err != nil {
		return "", fmt.Errorf("failed to armor key: %w", err)
	}
	return armored, nil
}

// EncryptBeforeSave encrypts all keys in the Keybox and writes the encrypted Keybox to an [io.Writer].
//
// This method first encrypts each key stored in the Keybox using the provided Encryptor. It then serializes
// the Keybox, including the encrypted keys, into JSON format and writes it to the provided [io.Writer].
//
// Note:
//   - This method ensures that keys are securely encrypted (effective for private keys) before being saved or transmitted over the network.
//   - The encryption process uses the public keys contained within the provided Encryptor.
//   - It is important to ensure that the Encryptor is properly initialized with valid public keys capable of encryption.
func (kb *Keybox) EncryptBeforeSave(o io.Writer, encryptor *Encryptor) error {
	encryptedKeys := []KeyMetadataEncrypted{}

	for _, keyInfo := range kb.Keys {
		encryptedKey, err := encryptor.encryptArmored(keyInfo.ArmoredKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt key: %w", err)
		}

		encryptedKeys = append(encryptedKeys, KeyMetadataEncrypted{
			Encrypted: encryptedKey,
		})
	}

	// Create a temporary structure to hold the UUID and encrypted keys
	type EncryptedKeybox struct {
		UUID      string                 `json:"uuid"`
		TotalKeys int                    `json:"total_keys"`
		Keys      []KeyMetadataEncrypted `json:"keys"`
	}

	encryptedKeybox := EncryptedKeybox{
		UUID:      kb.UUID,
		TotalKeys: kb.TotalKeys,
		Keys:      encryptedKeys,
	}

	// Now we can perform this operation over the network, especially when using Kubernetes. It's very smooth sailing.
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		data, err := sonic.MarshalIndent(encryptedKeybox, "", "  ")
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to marshal encrypted keybox: %w", err))
			return
		}

		if _, err := pw.Write(data); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to write encrypted keybox: %w", err))
			return
		}
	}()

	if _, err := io.Copy(o, pr); err != nil {
		return fmt.Errorf("failed to copy data to writer: %w", err)
	}

	return nil
}

// DeleteKey removes a key from the Keybox based on its fingerprint.
//
// This method searches for a key with the specified fingerprint within the Keybox.
// If the key is found, it is removed from the internal slice of keys, and the total
// key count is decremented. If the key is not found, an error is returned.
//
// Note:
//   - Ensure that the fingerprint provided is accurate and corresponds to a key
//     currently stored in the Keybox.
func (kb *Keybox) DeleteKey(fingerprints []string) error {
	fingerprintMap := make(map[string]bool)
	for _, fp := range fingerprints {
		fingerprintMap[fp] = true
	}

	originalCount := kb.TotalKeys
	kb.Keys = filterKeys(kb.Keys, fingerprintMap)
	kb.TotalKeys = kb.KeyCount()

	if originalCount == kb.TotalKeys {
		return fmt.Errorf("key with fingerprint %s not found", fingerprints)
	}

	return nil
}

func filterKeys(keys []KeyMetadata, fingerprintMap map[string]bool) []KeyMetadata {
	filtered := keys[:0] // Reuse the original slice
	for _, key := range keys {
		if !fingerprintMap[key.Fingerprint] {
			filtered = append(filtered, key)
		}
	}
	return filtered
}
