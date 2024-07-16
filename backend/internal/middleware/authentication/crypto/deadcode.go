// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// DecryptLargeData decrypts large data using AES decryption with the same derived encryption key used during encryption.
// It reads the encrypted data from the provided io.Reader and writes the decrypted data to the provided io.Writer.
// It verifies the signature of the encrypted data before decrypting.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the decryption key.
//
// Depcreated: Use streaming chunk instead.
func DecryptLargeData(src io.Reader, dst io.Writer, useArgon2 bool, secryptKey, signKey string) error {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(src, salt); err != nil {
		return err
	}

	key := deriveKey(salt, useArgon2, secryptKey)

	hash := sha256.New()
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		decrypted, err := decrypt(buf[:n], key)
		if err != nil {
			return err
		}

		if _, err := dst.Write(decrypted); err != nil {
			return err
		}

		hash.Write(buf[:n])
	}

	signature := make([]byte, sha256.Size)
	if _, err := io.ReadFull(src, signature); err != nil {
		return err
	}

	if !verifySignature(hash.Sum(nil), signature, signKey) {
		return ErrorInvalidSignature
	}

	return nil
}

// EncryptLargeData encrypts large data using AES encryption with a derived encryption key.
// It reads the data from the provided io.Reader and writes the encrypted data to the provided io.Writer.
// It generates a signature for the encrypted data and appends it to the output.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the encryption key.
//
// Depcreated: Use streaming chunk instead.
func EncryptLargeData(src io.Reader, dst io.Writer, useArgon2 bool, secryptKey, signKey string) error {
	return processLargeData(src, dst, useArgon2, secryptKey, signKey, encrypt)
}

// processLargeData is a higher-order function that processes large data using the provided processor function.
// It reads the data from the provided io.Reader and writes the processed data to the provided io.Writer.
// The processor function is responsible for encrypting or decrypting the data.
// It generates a signature for the processed data and appends it to the output.
//
// Depcreated: Use streaming chunk instead.
func processLargeData(src io.Reader, dst io.Writer, useArgon2 bool, secryptKey, signKey string, processor func([]byte, []byte) ([]byte, error)) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(salt, useArgon2, secryptKey)

	if _, err := dst.Write(salt); err != nil {
		return err
	}

	hash := sha256.New()
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		processed, err := processor(buf[:n], key)
		if err != nil {
			return err
		}

		if _, err := dst.Write(processed); err != nil {
			return err
		}

		hash.Write(processed)
	}

	signature := signData(hash.Sum(nil), signKey)
	if _, err := dst.Write(signature); err != nil {
		return err
	}

	return nil
}
