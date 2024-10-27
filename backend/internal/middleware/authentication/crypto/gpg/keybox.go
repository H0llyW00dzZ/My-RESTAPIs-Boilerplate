// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/bytedance/sonic"
)

// Keybox manages a collection of keys that can be stored and retrieved.
type Keybox struct {
	Keys []crypto.Key `json:"-"`
}

// NewKeybox creates a new Keybox instance.
func NewKeybox() *Keybox {
	return &Keybox{
		Keys: []crypto.Key{},
	}
}

// AddKey adds a new key to the Keybox, supporting multiple purposes.
func (kb *Keybox) AddKey(armoredKey string) error {
	key, err := crypto.NewKeyFromArmored(armoredKey)
	if err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}

	kb.Keys = append(kb.Keys, *key)
	return nil
}

// Save saves the Keybox to an [io.Writer] in JSON format.
//
// Note: Since it allow supports multiple purposes, it's recommended to store it in a file (e.g., over the network smiliar encrypt stream), network storage, or a database that can handle this object.
// Avoid using GPG key handling mechanisms that store keys directly in memory (bad), as it inefficient for a large number of keys.
func (kb *Keybox) Save(w io.Writer) error {
	var armoredKeys []string
	for _, key := range kb.Keys {
		armored, err := key.Armor()
		if err != nil {
			return fmt.Errorf("failed to armor key: %w", err)
		}
		armoredKeys = append(armoredKeys, armored)
	}

	data, err := sonic.Marshal(armoredKeys)
	if err != nil {
		return fmt.Errorf("failed to marshal keybox: %w", err)
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("failed to write keybox: %w", err)
	}

	return nil
}

// Load loads a Keybox from an [io.Reader] in JSON format.
//
// Note: Since it allow supports multiple purposes, it's recommended to store it in a file (e.g., over the network smiliar encrypt stream), network storage, or a database that can handle this object.
// Avoid using GPG key handling mechanisms that store keys directly in memory (bad), as it inefficient for a large number of keys.
func Load(r io.Reader) (*Keybox, error) {
	var armoredKeys []string
	if err := sonic.ConfigDefault.NewDecoder(r).Decode(&armoredKeys); err != nil {
		return nil, fmt.Errorf("failed to decode keybox: %w", err)
	}

	kb := NewKeybox()
	for _, armoredKey := range armoredKeys {
		if err := kb.AddKey(armoredKey); err != nil {
			return nil, err
		}
	}

	return kb, nil
}

// GetEncryptor creates an Encryptor from the keys in the Keybox that can be used for encryption.
func (kb *Keybox) GetEncryptor() (*Encryptor, error) {
	var encryptKeys []string
	for _, key := range kb.Keys {
		if key.CanEncrypt() {
			armored, err := key.Armor()
			if err != nil {
				return nil, fmt.Errorf("failed to armor key: %w", err)
			}
			encryptKeys = append(encryptKeys, armored)
		}
	}

	if len(encryptKeys) == 0 {
		return nil, ErrorCantEncrypt
	}

	return NewEncryptor(encryptKeys)
}
