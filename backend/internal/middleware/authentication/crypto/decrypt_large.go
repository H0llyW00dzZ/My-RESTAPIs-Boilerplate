// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/sha256"
	"io"
)

// DecryptLargeData decrypts large data using AES decryption with the same derived encryption key used during encryption.
// It reads the encrypted data from the provided io.Reader and writes the decrypted data to the provided io.Writer.
// It verifies the signature of the encrypted data before decrypting.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the decryption key.
//
// TODO: Improve this.
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
