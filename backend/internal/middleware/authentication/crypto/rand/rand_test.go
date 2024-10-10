// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand_test

import (
	"crypto/elliptic"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	"testing"
)

// Note: If this function is not safe for concurrent use, it would lead to race conditions or produce the same result
// across multiple goroutines. This function is particularly useful in production environments where the system requires
// multiple goroutines (e.g., when 10000 goroutines are needed, it provides 10000 readers, each with its own instance so it always welcome).
func TestFixedSize32Bytes(t *testing.T) {
	r := rand.FixedSize32Bytes()

	// Test reading from the reader
	buf := make([]byte, 32)
	n, err := r.Read(buf)

	// Check the number of bytes read
	if n != 32 {
		t.Errorf("Expected to read 32 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test reading again to ensure it generates new random bytes
	buf2 := make([]byte, 32)
	_, err = r.Read(buf2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that the two reads generate different random bytes
	if string(buf) == string(buf2) {
		t.Error("Expected different random bytes on subsequent reads")
	}
}

// Note: If this function is not safe for concurrent use, it would lead to race conditions or produce the same result
// across multiple goroutines. This function is particularly useful in production environments where the system requires
// multiple goroutines (e.g., when 10000 goroutines are needed, it provides 10000 readers, each with its own instance so it always welcome).
func TestFixedReaderRead(t *testing.T) {
	r := rand.FixedSize32Bytes()

	// Test reading with a buffer smaller than the fixed size
	buf := make([]byte, 16)
	n, err := r.Read(buf)

	// Check the number of bytes read
	if n != 16 {
		t.Errorf("Expected to read 16 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test reading with a buffer larger than the fixed size
	buf = make([]byte, 64)
	n, err = r.Read(buf)

	// Check the number of bytes read
	if n != 32 {
		t.Errorf("Expected to read 32 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test reading with an empty buffer
	buf = make([]byte, 0)
	n, err = r.Read(buf)

	// Check the number of bytes read
	if n != 0 {
		t.Errorf("Expected to read 0 bytes, but read %d bytes", n)
	}

	// Check for any errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Note: If this function is not safe for concurrent use, it would lead to race conditions or produce the same result
// across multiple goroutines. This function is particularly useful in production environments where the system requires
// multiple goroutines (e.g., when 10000 goroutines are needed, it provides 10000 readers, each with its own instance so it always welcome).
func TestFixedSizeECDSA(t *testing.T) {
	curves := []elliptic.Curve{
		elliptic.P256(),
		elliptic.P384(),
		elliptic.P521(),
	}

	for _, curve := range curves {
		r := rand.FixedSizeECDSA(curve)

		// Determine the expected size based on the curve
		expectedSize := 32
		if curve.Params().BitSize > 256 {
			expectedSize = 48
		}

		// Test reading from the reader
		buf := make([]byte, expectedSize)
		n, err := r.Read(buf)

		// Check the number of bytes read
		if n != expectedSize {
			t.Errorf("Expected to read %d bytes for curve %s, but read %d bytes", expectedSize, curve.Params().Name, n)
		}

		// Check for any errors
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", curve.Params().Name, err)
		}

		// Test reading again to ensure it generates new random bytes
		buf2 := make([]byte, expectedSize)
		_, err = r.Read(buf2)
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", curve.Params().Name, err)
		}

		// Check that the two reads generate different random bytes
		if string(buf) == string(buf2) {
			t.Errorf("Expected different random bytes on subsequent reads for curve %s", curve.Params().Name)
		}
	}
}

func TestFixedSizeECC(t *testing.T) {
	curves := []struct {
		curve        elliptic.Curve
		expectedSize int
	}{
		{elliptic.P256(), 32},
		{elliptic.P384(), 48},
		{elliptic.P521(), 66},
	}

	for _, c := range curves {
		r := rand.FixedSizeECC(c.curve)

		// Test reading from the reader
		buf := make([]byte, c.expectedSize)
		n, err := r.Read(buf)

		// Check the number of bytes read
		if n != c.expectedSize {
			t.Errorf("Expected to read %d bytes for curve %s, but read %d bytes", c.expectedSize, c.curve.Params().Name, n)
		}

		// Check for any errors
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", c.curve.Params().Name, err)
		}

		// Test reading again to ensure it generates new random bytes
		buf2 := make([]byte, c.expectedSize)
		_, err = r.Read(buf2)
		if err != nil {
			t.Errorf("Unexpected error for curve %s: %v", c.curve.Params().Name, err)
		}

		// Check that the two reads generate different random bytes
		if string(buf) == string(buf2) {
			t.Errorf("Expected different random bytes on subsequent reads for curve %s", c.curve.Params().Name)
		}
	}
}
