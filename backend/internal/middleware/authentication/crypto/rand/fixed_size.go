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
//
// Deprecated: Use FixedSizeECC instead.
func FixedSizeECDSA(curve elliptic.Curve) io.Reader {
	if curve.Params().BitSize <= 256 {
		return FixedSize32Bytes()
	}
	return &fixedReader{
		size: 48,
	}
}

// FixedSizeECC returns an [io.Reader] that provides a fixed-size random byte stream,
// suitable for generating nonces for elliptic curve cryptography (ECC). It can be used
// for both ECDSA and ECDH operations. The size of the random byte stream is determined
// by the elliptic curve's bit size, ensuring that the number of random bytes is sufficient
// for secure cryptographic operations.
//
// The function calculates the byte size needed for the given curve by rounding up the bit size
// to the nearest byte boundary. This ensures that even if the bit size is not a multiple of 8,
// the byte size will be sufficient.
//
// Example byte sizes for common curves:
//   - P-224: Bit size is 224. Byte size is (224 + 7) / 8 = 28 bytes.
//   - P-256: Bit size is 256. Byte size is (256 + 7) / 8 = 32 bytes.
//   - P-384: Bit size is 384. Byte size is (384 + 7) / 8 = 48 bytes.
//   - P-521: Bit size is 521. Byte size is (521 + 7) / 8 = 66 bytes.
//
// Note: This function is safe for use by multiple goroutines simultaneously.
func FixedSizeECC(curve elliptic.Curve) io.Reader {
	// This effectively Go rounds up (performs integer division), ensuring the correct number of bytes.
	//
	// Playground: https://go.dev/play/p/6oPu30iYILF
	//
	// Note: This may differ from calculator results, which include decimals.
	// Avoid using a calculator for this calculation, as it can be confusing;
	// this approach ensures the correct number of bytes for cryptographic purposes.
	bitSize := curve.Params().BitSize
	byteSize := (bitSize + 7) / 8
	return &fixedReader{
		size: byteSize,
	}
}

// FixedSizeRSA returns an [io.Reader] that provides a fixed-size random byte stream,
// suitable for generating nonces or blinding values for RSA operations. The size of
// the random byte stream is determined by the RSA modulus size in bytes.
//
// This function ensures that the number of random bytes is sufficient for secure
// cryptographic operations with RSA keys.
//
// Example usage:
//   - For a 2048-bit RSA key, the byte size is 2048 / 8 = 256 bytes.
//
// Note: This function is safe for use by multiple goroutines simultaneously.
// Consider using ECC instead of RSA, as RSA keys are larger and can consume more
// bandwidth, especially in HTTPS/TLS scenarios. For similar security levels, ECC
// (e.g., P-256) is more efficient than RSA with a 3072-bit modulus. Essentially,
// ECC P-256 offers similar security to RSA with a 3072-bit modulus, while ECC P-224
// is comparable to RSA with a 2048-bit modulus, which is commonly used in HTTPS/TLS.
func FixedSizeRSA(modulusBits int) io.Reader {
	// Calculate the byte size needed for the RSA modulus.
	byteSize := modulusBits / 8
	return &fixedReader{
		size: byteSize,
	}
}
