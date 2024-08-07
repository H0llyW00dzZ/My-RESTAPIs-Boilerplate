// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: The secure cryptographic random generator fixed size is moved here for easier maintenance.

package rand

import (
	"crypto/elliptic"
	"crypto/rand"
	"io"
)

// FixedSize32Bytes returns a custom [io.Reader] that provides a fixed-size random byte stream.
// The returned reader generates 32 random bytes each time it is read from.
// It uses the cryptographic random generator from the [crypto/rand] package to ensure secure randomness.
//
// The RandTLS function is suitable for use as the Rand field in [tls.Config] to provide
// a source of entropy for nonces and RSA blinding. It ensures that the TLS package
// always receives 32 random bytes when it requests random data.
//
// Example usage:
//
//	tlsConfig := &tls.Config{
//		// ...
//		Rand: rand.FixedSize32Bytes(),
//		// ...
//	}
//
// Note: This helper function is safe for use by multiple goroutines that call it simultaneously.
// Also note that the fixed reader of 32 random bytes is a well-known entropy size for nonces and RSA blinding. When captured in Wireshark,
// it is always unique. Plus, it is suitable for use by multiple goroutines because it provides an independent reader for each goroutine,
// and the size cannot be changed or increased.
func FixedSize32Bytes() io.Reader {
	return &fixedReader{
		size: 32,
	}
}

// fixedReader is a custom [io.Reader] implementation that provides a fixed-size random byte stream.
// It generates random bytes using the cryptographic random generator from the [crypto/rand] package.
type fixedReader struct {
	size int
}

// Read fills the provided byte slice p with random bytes up to the specified size.
// It returns the number of bytes read (n) and any error encountered.
//
// If the length of p is 0, Read returns immediately with n=0 and err=nil.
//
// If the length of p is less than the specified size, Read fills the entire buffer p
// with random bytes and returns the number of bytes read (n) and any error encountered.
//
// If the length of p is greater than or equal to the specified size, Read fills the first
// size bytes of p with random bytes and returns the number of bytes read (n) and any error encountered.
func (r *fixedReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(p) < r.size {
		// If the provided buffer is smaller than the fixed size,
		// read as much as possible and return the number of bytes read.
		return rand.Read(p)
	}

	return rand.Read(p[:r.size])
}

// FixedSizeECDSA returns a custom [io.Reader] that provides a fixed-size random byte stream
// suitable for generating ECDSA nonces. The returned reader generates a fixed number of random bytes
// each time it is read from, based on the provided elliptic curve.
// It uses the cryptographic random generator from the [crypto/rand] package to ensure secure randomness.
//
// The size of the random byte stream is determined by the curve:
//   - For curves with a bit size less than or equal to 256 (e.g., P-256), it returns the FixedSize32Bytes reader.
//   - For curves with a bit size greater than 256 (e.g., P-384, P-521), it generates 48 random bytes.
//
// Note: This helper function is safe for use by multiple goroutines that call it simultaneously.
// Also note that the P-521 curve might not be suitable for TLS or other applications because most internet
// on average use curves ranging from P-256 to P-384.
func FixedSizeECDSA(curve elliptic.Curve) io.Reader {
	if curve.Params().BitSize <= 256 {
		return FixedSize32Bytes()
	}
	return &fixedReader{
		size: 48,
	}
}
