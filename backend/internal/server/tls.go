// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"bytes"
	"crypto/tls"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
)

// streamConn is a wrapper struct that combines a TLS 1.3 connection and a Hybrid Scheme (Stream) for encrypted communication.
//
// Note: Currently unused, and marked as TODO, will complete implementing this later
type streamConn struct {
	// Note: This is already connected because [stream.Stream] is the core of cryptographic operations.
	// It can be used to write/read over the network, for example, to store encrypted data in a database.
	*tls.Conn
	*stream.Stream
}

// Read reads encrypted data from the TLS connection, decrypts it using the Stream, and returns the decrypted data.
func (c *streamConn) Read() ([]byte, error) {
	var buffer bytes.Buffer
	err := c.Stream.Decrypt(c.Conn, &buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Write encrypts the provided data using the Stream and writes it to the TLS connection.
func (c *streamConn) Write(data []byte) error {
	return c.Stream.Encrypt(bytes.NewReader(data), c.Conn)
}
