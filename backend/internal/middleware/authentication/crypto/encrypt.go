// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"

	// Import the godotenv package for loading environment variables from a .env file
	// The "_" blank identifier is used to import the package for its side effects (auto-loading .env file)
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/argon2"
)

var (
	// secryptkey holds the secret encryption key retrieved from the environment variable "SECRETCRYPT_KEY"
	secryptkey = os.Getenv("SECRETCRYPT_KEY")
	// signkey holds the secret signing key retrieved from the environment variable "SIGN_KEY"
	signkey = os.Getenv("SIGN_KEY")
)

var (
	// ErrorInvalidCipherText is a custom error variable that represents an error
	// which occurs when the ciphertext (encrypted text) is invalid or malformed.
	ErrorInvalidCipherText = errors.New("invalid ciphertext")
	// ErrorInvalidSignature is a custom error variable that represents an error
	// which occurs when the signature is invalid or does not match the expected signature.
	ErrorInvalidSignature = errors.New("invalid signature")
)

// signData generates an HMAC signature for the given data using the signing key.
func signData(data []byte) []byte {
	mac := hmac.New(sha256.New, []byte(signkey))
	mac.Write(data)
	return mac.Sum(nil)
}

// verifySignature verifies the HMAC signature of the given data using the signing key.
func verifySignature(data, signature []byte) bool {
	expectedMAC := signData(data)
	return hmac.Equal(signature, expectedMAC)
}

// deriveKey derives an encryption key using Argon2 key derivation function or returns the secryptkey directly.
func deriveKey(salt []byte, useArgon2 bool) []byte {
	// Note: Using Argon2 is expensive (100MB+ per encrypt/decrypt) the cost, which is not recommended.
	// I might try to introduce cryptographic techniques to implement a similar but cheaper approach and a new cipher from scratch later.
	if useArgon2 {
		return argon2.IDKey([]byte(secryptkey), salt, 1, 64*1024, 4, 32)
	}
	return []byte(secryptkey)
}

// encrypt encrypts the given data using AES encryption with the provided key and returns the ciphertext.
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return append(nonce, ciphertext...), nil
}

// decrypt decrypts the given ciphertext using AES decryption with the provided key and returns the plaintext.
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrorInvalidCipherText
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptData encrypts the given data using AES encryption with a derived encryption key and signs the ciphertext.
// It returns the base64-encoded ciphertext and signature.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the encryption key.
func EncryptData(data string, useArgon2 bool) (string, string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	key := deriveKey(salt, useArgon2)
	ciphertext, err := encrypt([]byte(data), key)
	if err != nil {
		return "", "", err
	}

	encryptedData := append(salt, ciphertext...)
	signature := signData(encryptedData)

	encodedData := base64.StdEncoding.EncodeToString(encryptedData)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	return encodedData, encodedSignature, nil
}

// DecryptData decrypts the given encrypted data using AES decryption with the same derived encryption key used during encryption.
// It expects the encrypted data and signature to be base64-encoded.
// It verifies the signature before decrypting the data.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the decryption key.
func DecryptData(encryptedData, signature string, useArgon2 bool) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return "", err
	}

	if !verifySignature(encryptedBytes, signatureBytes) {
		return "", ErrorInvalidSignature
	}

	salt := encryptedBytes[:16]
	ciphertext := encryptedBytes[16:]

	key := deriveKey(salt, useArgon2)
	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// processLargeData is a higher-order function that processes large data using the provided processor function.
// It reads the data from the provided io.Reader and writes the processed data to the provided io.Writer.
// The processor function is responsible for encrypting or decrypting the data.
// It generates a signature for the processed data and appends it to the output.
func processLargeData(src io.Reader, dst io.Writer, useArgon2 bool, processor func([]byte, []byte) ([]byte, error)) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(salt, useArgon2)

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

	signature := signData(hash.Sum(nil))
	if _, err := dst.Write(signature); err != nil {
		return err
	}

	return nil
}

// EncryptLargeData encrypts large data using AES encryption with a derived encryption key.
// It reads the data from the provided io.Reader and writes the encrypted data to the provided io.Writer.
// It generates a signature for the encrypted data and appends it to the output.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the encryption key.
func EncryptLargeData(src io.Reader, dst io.Writer, useArgon2 bool) error {
	return processLargeData(src, dst, useArgon2, encrypt)
}

// DecryptLargeData decrypts large data using AES decryption with the same derived encryption key used during encryption.
// It reads the encrypted data from the provided io.Reader and writes the decrypted data to the provided io.Writer.
// It verifies the signature of the encrypted data before decrypting.
// If useArgon2 is true, it uses Argon2 key derivation function to derive the decryption key.
func DecryptLargeData(src io.Reader, dst io.Writer, useArgon2 bool) error {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(src, salt); err != nil {
		return err
	}

	key := deriveKey(salt, useArgon2)

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

	if !verifySignature(hash.Sum(nil), signature) {
		return ErrorInvalidSignature
	}

	return nil
}
