// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream_test

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"testing"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
)

func TestHybridEncryptDecryptStream(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(encryptedBuffer, decryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext. Got: %s, Want: %s", decryptedData, plaintext)
	}
}

func TestHybridEncryptDecryptStreamWithApiKey(t *testing.T) {
	// Predefined API keys or secret keys, which should be securely stored and retrieved.
	aesKey := []byte("gopher-testing-testing-testinggg")
	chachaKey := []byte("gopher-testing-testing-testinggg")

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(encryptedBuffer, decryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext. Got: %s, Want: %s", decryptedData, plaintext)
	}
}

func TestHybridEncryptDecryptStreamLargeData(t *testing.T) {
	// Note: Works well testing on AMD Ryzen 9 3900x 12-Core Processor (24 CPUs) RAM 32GB
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate a large plaintext data.
	plaintextSize := 10 * 1024 * 1024 // 10 MB
	plaintext := make([]byte, plaintextSize)
	_, err = rand.Read(plaintext)
	if err != nil {
		t.Fatalf("Failed to generate plaintext: %v", err)
	}

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(encryptedBuffer, decryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext.")
	}
}

func TestHybridEncryptDecryptStreamWithHMAC(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate a random HMAC key.
	hmacKey := make([]byte, 32)
	_, err = rand.Read(hmacKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC key: %v", err)
	}

	// Enable HMAC authentication.
	s.EnableHMAC(hmacKey)

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system with HMAC.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Calculate the HMAC digest of the encrypted data.
	hmacDigest, err := s.Digest(bytes.NewReader(encryptedData))
	if err != nil {
		t.Fatalf("Failed to calculate HMAC digest: %v", err)
	}

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(encryptedBuffer, decryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext. Got: %s, Want: %s", decryptedData, plaintext)
	}

	// Verify the HMAC digest.
	encryptedBuffer = bytes.NewBuffer(encryptedData)
	verifiedHMACDigest, err := s.Digest(encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to calculate HMAC digest for verification: %v", err)
	}

	t.Logf("Verified Checksum: %x", verifiedHMACDigest)

	if !bytes.Equal(verifiedHMACDigest, hmacDigest) {
		t.Errorf("HMAC verification failed. Expected: %x, Got: %x", hmacDigest, verifiedHMACDigest)
	}
}

func TestHybridEncryptDecryptStreamWithHMACHasBeenCompromised(t *testing.T) {
	// Generate random keys for AES and XChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate a random HMAC key.
	hmacKey := make([]byte, 32)
	_, err = rand.Read(hmacKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC key: %v", err)
	}

	// Enable HMAC authentication.
	s.EnableHMAC(hmacKey)

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system with HMAC.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Simulate unauthorized modification of the encrypted data.
	//
	// Let's say this Data has been Compromised.
	encryptedData[1] ^= 0xFF // Flip the first byte of the encrypted data.

	// Decrypt the data without calculating the HMAC digest (skipping step 2 and 3).
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(bytes.NewBuffer(encryptedData), decryptedBuffer)
	if err == nil {
		t.Errorf("Decryption succeeded despite unauthorized modification.")
	} else {
		t.Logf("Decryption failed as expected: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data matches original plaintext despite unauthorized modification.")
	}
}

// Test an additional layer of security on top of the strong (3-key) authentication.
func TestHybridEncryptDecryptStreamWithWrongHMACKey(t *testing.T) {
	// Generate random keys for AES and XChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate a random HMAC key.
	hmacKey := make([]byte, 32)
	_, err = rand.Read(hmacKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC key: %v", err)
	}

	// Enable HMAC authentication.
	s.EnableHMAC(hmacKey)

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system with the wrong HMAC key.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Generate a different HMAC key.
	wrongHMACKey := make([]byte, 32)
	_, err = rand.Read(wrongHMACKey)
	if err != nil {
		t.Fatalf("Failed to generate wrong HMAC key: %v", err)
	}

	// Decrypt the data using the wrong HMAC key.
	s.EnableHMAC(wrongHMACKey)
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(encryptedBuffer, decryptedBuffer)
	if err == nil {
		t.Errorf("Decryption succeeded with the wrong HMAC key.")
	}
}

// Test without calculating or collecting the digest or checksum.
func TestHybridEncryptDecryptStreamWithHMACWithoutDigestorChecksumCollected(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate a random HMAC key.
	hmacKey := make([]byte, 32)
	_, err = rand.Read(hmacKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC key: %v", err)
	}

	// Enable HMAC authentication.
	s.EnableHMAC(hmacKey)

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of the hybrid encryption system with HMAC, without collecting a digest or checksum.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Decrypt the data.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(encryptedBuffer, decryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare the decrypted data to the original plaintext.
	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(decryptedData, plaintext) {
		t.Errorf("Decrypted data does not match original plaintext. Got: %s, Want: %s", decryptedData, plaintext)
	}

}

// errorReader is a custom reader that always returns an error.
// It is used for testing purposes in low-level operations related to cryptography.
type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

// errorWriter is a custom writer that always returns an error.
// It is used for testing purposes in low-level operations related to cryptography.
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("simulated write error")
}

func TestHybridEncryptDecryptStreamErrorHandling(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Test encryption error handling.
	t.Run("EncryptionError", func(t *testing.T) {
		// Create an error-producing input reader.
		errorReader := &errorReader{}

		// Attempt to encrypt data with the error-producing reader.
		encryptedBuffer := new(bytes.Buffer)
		err = s.Encrypt(errorReader, encryptedBuffer)
		if err == nil {
			t.Errorf("Expected encryption error, but got nil.")
		}
	})

	// Test decryption error handling.
	t.Run("DecryptionError", func(t *testing.T) {
		// Create an error-producing input reader.
		errorReader := &errorReader{}

		// Attempt to decrypt data with the error-producing reader.
		decryptedBuffer := new(bytes.Buffer)
		err = s.Decrypt(errorReader, decryptedBuffer)
		if err == nil {
			t.Errorf("Expected decryption error, but got nil.")
		}
	})

	// Test HMAC verification error handling.
	t.Run("HMACVerificationError", func(t *testing.T) {
		// Generate a random HMAC key.
		hmacKey := make([]byte, 32)
		_, err = rand.Read(hmacKey)
		if err != nil {
			t.Fatalf("Failed to generate HMAC key: %v", err)
		}

		// Enable HMAC authentication.
		s.EnableHMAC(hmacKey)

		// Simulate plaintext data to encrypt.
		plaintext := []byte("Hello, World! This is a test of HMAC verification error handling.")

		// Encrypt the data.
		inputBuffer := bytes.NewBuffer(plaintext)
		encryptedBuffer := new(bytes.Buffer)
		err = s.Encrypt(inputBuffer, encryptedBuffer)
		if err != nil {
			t.Fatalf("Failed to encrypt data: %v", err)
		}

		// Ensure the encrypted data buffer's read position is reset to the beginning.
		encryptedData := encryptedBuffer.Bytes()
		encryptedBuffer = bytes.NewBuffer(encryptedData)

		// Simulate unauthorized modification of the encrypted data.
		encryptedData[len(encryptedData)-1] ^= 0xFF // Flip the last byte of the encrypted data.

		// Attempt to decrypt the modified data.
		decryptedBuffer := new(bytes.Buffer)
		err = s.Decrypt(bytes.NewBuffer(encryptedData), decryptedBuffer)
		if err == nil {
			t.Errorf("Expected HMAC verification error, but got nil.")
		}
	})

	// Test readChunkMetadata error handling.
	t.Run("ReadChunkMetadataError", func(t *testing.T) {
		// Create an error-producing input reader.
		errorReader := &errorReader{}
		decryptedBuffer := new(bytes.Buffer)
		// Attempt to read chunk metadata with the error-producing reader.
		err = s.Decrypt(errorReader, decryptedBuffer)
		if err == nil {
			t.Errorf("Expected readChunkMetadata error, but got nil.")
		}
	})

	// Test readChunkMetadata EOF handling.
	t.Run("ReadChunkMetadataEOF", func(t *testing.T) {
		// Create an input reader with incomplete chunk metadata to simulate EOF.
		incompleteMetadata := []byte{0x00} // Only one byte instead of the required chunk size and nonce
		incompleteReader := bytes.NewReader(incompleteMetadata)
		decryptedBuffer := new(bytes.Buffer)

		// Attempt to read chunk metadata from the empty reader.
		err = s.Decrypt(incompleteReader, decryptedBuffer)
		if err == io.EOF {
			t.Errorf("Expected io.EOF error, but got: %v", err)
		}
	})

	// Test readChunkMetadata partial read handling.
	t.Run("ReadChunkMetadataPartialRead", func(t *testing.T) {
		// Create an input reader with incomplete chunk metadata.
		incompleteMetadata := []byte{0x00} // Only one byte instead of the required chunk size and nonce
		incompleteReader := bytes.NewReader(incompleteMetadata)
		decryptedBuffer := new(bytes.Buffer)

		// Attempt to read chunk metadata from the incomplete reader.
		err = s.Decrypt(incompleteReader, decryptedBuffer)
		if err == nil {
			t.Errorf("Expected readChunkMetadata error for partial read, but got nil.")
		}
	})
}

func TestHybridDecryptStreamInvalidHMACDigestSize(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Generate a random HMAC key.
	hmacKey := make([]byte, 32)
	_, err = rand.Read(hmacKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC key: %v", err)
	}

	// Enable HMAC authentication.
	s.EnableHMAC(hmacKey)

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of invalid HMAC digest size.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()

	// Simulate an invalid HMAC digest size by truncating the encrypted data.
	invalidEncryptedData := encryptedData[:len(encryptedData)-1]

	// Attempt to decrypt the data with the invalid HMAC digest size.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(bytes.NewBuffer(invalidEncryptedData), decryptedBuffer)
	if err == nil {
		t.Errorf("Expected decryption error due to invalid HMAC digest size, but got nil.")
	} else if err.Error() != "unexpected EOF" {
		t.Errorf("Expected 'unexpected EOF' error, but got: %v", err)
	}
}

func TestHybridDecryptStreamEncryptedChunkSizeMismatch(t *testing.T) {
	// Generate random keys for AES and ChaCha20-Poly1305.
	aesKey := make([]byte, 32)    // AES-256 requires a 32-byte key.
	chachaKey := make([]byte, 32) // XChaCha20-Poly1305 uses a 32-byte key.

	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatalf("Failed to generate XChaCha20-Poly1305 key: %v", err)
	}

	// Create a new Stream instance.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream instance: %v", err)
	}

	// Simulate plaintext data to encrypt.
	plaintext := []byte("Hello, World! This is a test of encrypted chunk size mismatch.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()

	// Simulate an encrypted chunk size mismatch by modifying the chunk size.
	encryptedData[1] ^= 0xFF // Flip the first byte of the chunk size.

	// Attempt to decrypt the data with the encrypted chunk size mismatch.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(bytes.NewBuffer(encryptedData), decryptedBuffer)
	if err == nil {
		t.Errorf("Expected decryption error due to encrypted chunk size mismatch, but got nil.")
	} else if err.Error() != "unexpected EOF" {
		t.Errorf("Expected 'unexpected EOF' error, but got: %v", err)
	}
}
