// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package stream provides a hybrid encryption scheme that combines AES-CTR and XChaCha20-Poly1305 algorithms
// for secure encryption and decryption of data streams. It also supports optional HMAC authentication for added
// integrity and authenticity.
//
// Compatibility:
//
// It's important to note that the stream package may not be directly compatible with browsers or other non-Go
// environments. The specific ciphertext format used by this package relies on the implementation details of the
// AES-CTR and XChaCha20-Poly1305 algorithms in Go, along with custom chunking and HMAC authentication mechanisms.
//
// The stream package is primarily designed for internal use within Go applications. It provides a secure and
// efficient way to encrypt and decrypt data streams, making it suitable for internal or private communication
// between Go services or components.
//
// Recommended Usage:
//
// The stream package is recommended for use in internal or private communication scenarios within a Go application
// or between Go services. It can be used to secure sensitive data transmission, such as:
//
// - Communication between microservices in a distributed system
// - Encryption of data stored in databases or files
// - Secure transmission of data between different components of a Go application
// - VPN (Virtual Private Network)
package stream
