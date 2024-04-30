// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import "io"

// EncryptLargeData encrypts large data using AES encryption with a derived encryption key.
// It reads the data from the provided io.Reader and writes the encrypted data to the provided io.Writer.
// It generates a signature for the encrypted data and appends it to the output.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the encryption key.
//
// TODO: Improve this.
func EncryptLargeData(src io.Reader, dst io.Writer, useArgon2 bool) error {
	return processLargeData(src, dst, useArgon2, encrypt)
}
