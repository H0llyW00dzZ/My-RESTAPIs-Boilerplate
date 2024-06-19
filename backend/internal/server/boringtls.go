// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"bytes"
	"crypto/tls"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
	"net"
	"time"
)

// streamConn is a wrapper struct that combines a TLS 1.3 connection and a Hybrid Scheme (Stream) for encrypted communication.
//
// Note: Currently unused and marked as TODO. Implementation will be completed later (see https://tip.golang.org/doc/go1.23).
type streamConn struct {
	// Note: This is already connected because [stream.Stream] is the core of cryptographic operations.
	// It can be used to write/read over the network, for example, to store encrypted data in a database.
	*tls.Conn
	*stream.Stream
}

// Read reads encrypted data from the TLS connection, decrypts it using the Stream, and returns the decrypted data.
//
// Note: Exercise caution when calling this method in relation to section 10.10.3 of the TLS Encrypted Client Hello (ECH) draft
// (see https://www.ietf.org/archive/id/draft-ietf-tls-esni-18.html#section-10.10.3).
// The reason for implementing this is that it is legal for the server owner to do so. However, it requires careful consideration to use correctly.
//
// TIP: To mitigate this risk from section [10.10.3], servers can implement rate limiting or other security measures to control the number of decryption operations they perform within a given time frame.
// By monitoring and limiting the rate of decryption requests, servers can reduce the impact of potential DoS attacks while still fulfilling their role in the ECH protocol.
func (c *streamConn) Read(b []byte) (int, error) {
	var buffer bytes.Buffer
	err := c.Stream.Decrypt(c.Conn, &buffer)
	if err != nil {
		return 0, err
	}
	return buffer.Read(b)
}

// Write encrypts the provided data using the Stream and writes it to the TLS connection.
func (c *streamConn) Write(b []byte) (int, error) {
	var buffer bytes.Buffer
	err := c.Stream.Encrypt(bytes.NewReader(b), &buffer)
	if err != nil {
		return 0, err
	}
	return c.Conn.Write(buffer.Bytes())
}

// Close closes the underlying TLS connection.
func (c *streamConn) Close() error {
	return c.Conn.Close()
}

// LocalAddr returns the local network address.
func (c *streamConn) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *streamConn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated with the connection.
func (c *streamConn) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
func (c *streamConn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
func (c *streamConn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

// streamListener is a custom listener that wraps the original listener and creates a streamConn for each TLS connection
type streamListener struct {
	net.Listener
	tlsConfig *tls.Config
	stream    *stream.Stream
}

// NewStreamConn creates a new streamConn instance by wrapping a TLS connection and a Stream.
func NewStreamConn(tlsConn *tls.Conn, stream *stream.Stream) net.Conn {
	return &streamConn{
		Conn:   tlsConn,
		Stream: stream,
	}
}

// Accept waits for and returns the next connection to the listener, wrapped in a streamConn
func (l *streamListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Server(conn, l.tlsConfig)
	return NewStreamConn(tlsConn, l.stream), nil
}
