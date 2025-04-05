// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package vault

import (
	"context"
	"fmt"
)

// buildTransitPath constructs the full Transit API path with optional customization.
func (v *VClient) buildTransitPath(operation, keyName string) string {
	path := fmt.Sprintf("%s/%s/%s", v.parameters.TransitPath, operation, keyName)
	return path
}

// Encrypt encrypts data using Vault's Transit Secrets Engine.
func (v *VClient) Encrypt(ctx context.Context, keyName string, plaintext []byte) ([]byte, error) {
	encryptData := map[string]any{
		"plaintext": plaintext,
	}

	encryptResp, err := v.client.Logical().Write(v.buildTransitPath("encrypt", keyName), encryptData)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt data: %w", err)
	}

	ciphertext, ok := encryptResp.Data["ciphertext"].(string)
	if !ok {
		return nil, fmt.Errorf("ciphertext not found in response")
	}

	return []byte(ciphertext), nil
}

// Decrypt decrypts data using Vault's Transit Secrets Engine.
func (v *VClient) Decrypt(ctx context.Context, keyName string, ciphertext []byte) ([]byte, error) {
	decryptData := map[string]any{
		"ciphertext": ciphertext,
	}

	decryptResp, err := v.client.Logical().Write(v.buildTransitPath("decrypt", keyName), decryptData)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt data: %w", err)
	}

	plaintext, ok := decryptResp.Data["plaintext"].([]byte)
	if !ok {
		return nil, fmt.Errorf("plaintext not found in response")
	}

	return plaintext, nil
}
