// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"
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
	Keys []KeyMetadata `json:"keys"`
}

// NewKeybox creates a new Keybox instance.
func NewKeybox() *Keybox {
	return &Keybox{
		Keys: []KeyMetadata{},
	}
}

// AddKey adds a new key to the Keybox, supporting multiple purposes.
func (kb *Keybox) AddKey(armoredKey string) error {
	key, err := crypto.NewKeyFromArmored(armoredKey)
	if err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}

	creationDate := key.GetEntity().PrimaryKey.CreationTime.UTC().Format(time.RFC3339)

	keyInfo := KeyMetadata{
		Fingerprint:  key.GetFingerprint(),
		CreationDate: creationDate,
		ArmoredKey:   armoredKey,
	}

	kb.Keys = append(kb.Keys, keyInfo)
	return nil
}

// Save saves the Keybox to an [io.Writer] in JSON format.
//
// Note: Since it allow supports multiple purposes, it's recommended to store it in a file (e.g., over the network smiliar encrypt stream), network storage, or a database that can handle this object.
// Avoid using GPG key handling mechanisms that store keys directly in memory (bad), as it inefficient for a large number of keys.
func (kb *Keybox) Save(w io.Writer) error {
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

	_, err := io.Copy(w, pr)
	if err != nil {
		return fmt.Errorf("failed to copy data to writer: %w", err)
	}

	return nil
}

// Load loads a Keybox from an [io.Reader] in JSON format.
//
// Note: Since it allow supports multiple purposes, it's recommended to store it in a file (e.g., over the network smiliar encrypt stream), network storage, or a database that can handle this object.
// Avoid using GPG key handling mechanisms that store keys directly in memory (bad), as it inefficient for a large number of keys.
func Load(r io.Reader) (*Keybox, error) {
	// Now we can perform this operation over the network, especially when using Kubernetes. It's very smooth sailing.
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		if _, err := io.Copy(pw, r); err != nil {
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
