// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream_test

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"

	"golang.org/x/crypto/chacha20poly1305"
)

// Moved here because VSCode keeps crashing without any reason and then causes a blue screen.
// Windows it's so bad
const TempSizeData = 10 * 1024 * 1024 // 10 MB

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
	plaintext := make([]byte, TempSizeData)
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

func TestHybridEncryptDecryptStreamLargeDataWithHMACEnabled(t *testing.T) {
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

	// Generate a random HMAC key.
	hmacKey := make([]byte, 32)
	_, err = rand.Read(hmacKey)
	if err != nil {
		t.Fatalf("Failed to generate HMAC key: %v", err)
	}

	// Enable HMAC authentication.
	s.EnableHMAC(hmacKey)

	// Generate a large plaintext data.
	plaintext := make([]byte, TempSizeData)
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

// Let's Say this test simulates a scenario where the encrypted data buffer is vulnerable to a buffer overflow attack,
// which is exploitable in most languages (e.g., C/C++, Assembly, Python, Java, Ruby) but not in Go.
// Go is considered safe and suitable for cryptographic operations because it provides built-in protection against buffer overflow vulnerabilities.
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
	} else {
		t.Logf("Decryption failed as expected: %v", err)
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
			t.Logf("Encryption failed as expected: %v", err)
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
			t.Logf("Decryption failed as expected: %v", err)
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
			t.Logf("Decryption failed as expected: %v", err)
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
			t.Logf("Decryption failed as expected: %v", err)
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
			t.Logf("Decryption failed as expected: %v", err)
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
			t.Logf("Decryption failed as expected: %v", err)
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

	// Calculate the HMAC digest of the encrypted data.
	hmacDigest, err := s.Digest(bytes.NewReader(encryptedData))
	if err != nil {
		t.Fatalf("Failed to calculate HMAC digest: %v", err)
	}

	// Simulate an invalid HMAC digest size by truncating the encrypted data.
	invalidEncryptedData := encryptedData[len(encryptedData)-len(hmacDigest)-1:]

	// Attempt to decrypt the data with the invalid HMAC digest size.
	decryptedBuffer := new(bytes.Buffer)
	err = s.Decrypt(bytes.NewBuffer(invalidEncryptedData), decryptedBuffer)
	if err == nil {
		t.Errorf("Expected decryption error due to invalid HMAC digest size, but got nil.")
	} else if err.Error() != "XChacha20Poly1305: invalid HMAC digest size" {
		t.Errorf("Expected error message 'XChacha20Poly1305: invalid HMAC digest size', but got: %v", err)
	} else {
		t.Logf("Decryption failed as expected: %v", err)
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
	} else if err.Error() != "XChacha20Poly1305: encrypted chunk size mismatch" {
		t.Errorf("Expected error message 'XChacha20Poly1305: encrypted chunk size mismatch', but got: %v", err)
	} else {
		t.Logf("Decryption failed as expected: %v", err)
	}
}

// Let's Say this test simulates a scenario where the encrypted data buffer is vulnerable to a buffer overflow attack,
// which is exploitable in most languages (e.g., C/C++, Assembly, Python, Java, Ruby) but not in Go.
// Go is considered safe and suitable for cryptographic operations because it provides built-in protection against buffer overflow vulnerabilities.
func TestHybridEncryptDecryptStreamHasBeenCompromised(t *testing.T) {
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

	// Simulate unauthorized modification of the encrypted data.
	//
	// Let's say this Data has been Compromised.
	// Note: Without HMAC it's starting from 2
	encryptedData[2] ^= 0xFF // Flip the second byte of the encrypted data.

	// Decrypt the data.
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

func TestHybridEncryptDecryptStreamWithHMACDigestInvalidkey(t *testing.T) {
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

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()
	encryptedBuffer = bytes.NewBuffer(encryptedData)

	// Calculate the HMAC digest of the encrypted data.
	hmacDigest, err := s.Digest(bytes.NewReader(encryptedData))
	if err != nil {
		t.Fatalf("Failed to calculate HMAC digest: %v", err)
	}

	// Generate a different HMAC key.
	wrongHMACKey := make([]byte, 32)
	_, err = rand.Read(wrongHMACKey)
	if err != nil {
		t.Fatalf("Failed to generate wrong HMAC key: %v", err)
	}

	// Digest the data using the wrong HMAC key.
	s.EnableHMAC(wrongHMACKey)

	// Verify the HMAC digest.
	encryptedBuffer = bytes.NewBuffer(encryptedData)
	calculatedHMACDigest, err := s.Digest(encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to calculate HMAC digest for verification: %v", err)
	}

	if subtle.ConstantTimeCompare(calculatedHMACDigest, hmacDigest) == 1 {
		t.Errorf("HMAC digest verification succeeded with the wrong HMAC key.")
	} else {
		// Note: This output contains the raw encrypted ciphertext when HMAC authentication is enabled.
		// The ciphertext can be bound to an [io.Writer] output, such as a string builder, rune builder, or byte builder,
		// and can be used for various protocols. However, using this combination of AES-CTR and XChaCha20-Poly1305 for TLS/SSL is not recommended
		// as it may be slower compared to using pure XChaCha20-Poly1305.
		t.Logf("HMAC digest verification failed as expected %x, Got: %x", hmacDigest, calculatedHMACDigest)
	}
}

func TestDecryptUnexpectedChunk(t *testing.T) {
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

	var decryptedOutput bytes.Buffer

	// Test case: Decrypt encrypted data with an unexpected chunk
	//
	// Note: The maximum chunk buffer size is 2 and can be used for the HMAC Tag as the current identifier.
	// If set to 3 or more, it will lead to an unexpected NonceSizeX error. (see https://i.imgur.com/p8imLJp.png how it work)
	invalidEncryptedData := []byte("invalid-encrypted-data")
	shortBufferSize := 1
	decryptedOutput.Reset()
	decryptedOutput.Grow(shortBufferSize)
	invalidEncryptedInput := bytes.NewReader(invalidEncryptedData[:shortBufferSize])

	err = s.Decrypt(invalidEncryptedInput, &decryptedOutput)
	if err == nil {
		t.Errorf("Expected error due to buffer too short, but got nil.")
	} else if err.Error() != "XChacha20Poly1305: Unexpected Chunk Buffer Size" {
		t.Errorf("Expected error message 'XChacha20Poly1305: Unexpected Chunk Buffer Size', but got: %v", err)
	} else {
		t.Logf("Decryption failed as expected: %v", err)
	}
}

func TestHybridDecryptStreamXChaCha20NonceSizeXTooShort(t *testing.T) {
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
	plaintext := []byte("Hello, World! This is a test of a XChaCha20-Poly1305 NonceX that is too short.")

	// Encrypt the data.
	inputBuffer := bytes.NewBuffer(plaintext)
	encryptedBuffer := new(bytes.Buffer)
	err = s.Encrypt(inputBuffer, encryptedBuffer)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Ensure the encrypted data buffer's read position is reset to the beginning.
	encryptedData := encryptedBuffer.Bytes()

	// Extract the chunk size from the encrypted data.
	chunkSizeBuf := encryptedData[:2]
	chunkSize := binary.BigEndian.Uint16(chunkSizeBuf)

	// Calculate the total size of the chunk (chunk size + nonce size).
	totalChunkSize := int(chunkSize) + chacha20poly1305.NonceSizeX

	// Modify the encrypted data to have a buffer size smaller than the total chunk size.
	if len(encryptedData) >= totalChunkSize {
		// Create a reader for a subset of the encrypted data.
		shortBufferSize := 2 + chacha20poly1305.NonceSizeX/2 // Chunk size + half of the nonce size
		shortBufferReader := bytes.NewReader(encryptedData[:shortBufferSize])

		// Create a buffer to store the decrypted data.
		decryptedBuffer := new(bytes.Buffer)

		// Attempt to decrypt the data with the short buffer.
		err = s.Decrypt(shortBufferReader, decryptedBuffer)
		if err == nil {
			t.Errorf("Expected error due to buffer too short, but got nil.")
		} else if err.Error() != "XChacha20Poly1305: Unexpected NonceSizeX" {
			t.Errorf("Expected error message 'XChacha20Poly1305: Unexpected NonceSizeX', but got: %v", err)
		} else {
			t.Logf("Decryption failed as expected: %v", err)
		}
	}
}
