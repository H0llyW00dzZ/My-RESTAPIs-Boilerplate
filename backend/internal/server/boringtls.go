// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package server

import (
	"bytes"
	"crypto/tls"
	"errors"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
	"net"
	"time"
)

// streamConn is a wrapper struct that combines a TLS 1.3 connection and a Hybrid Scheme (Stream) for encrypted communication.
// It is designed on top of the TLS connection, providing enhanced security for the communication channel,
// unlike most TLS implementations in the world (e.g, not allowed to use another ciphertext).
//
// Note: Currently unused and marked as TODO. Implementation will be completed later (see https://tip.golang.org/doc/go1.23).
type streamConn struct {
	// Note: This is already connected because [stream.Stream] is the core of cryptographic operations.
	// It can be used to write/read over the network, for example, to store encrypted data in a database.
	// Also note that this wrapper is safe, instead of using an unsafe pointer (e.g., https://pkg.go.dev/unsafe). When it binds * into *tls.Conn and *stream.Stream,
	// it already binds everything for the client and server. It doesn't need to be modified as it is already stable.
	// Even if modified, it only copies the standard library and adds the stream encrypter/decrypter along with stream.New for the cipher suites that used for key sharing.
	// This design allows the use of another cipher (due to TLS 1.3 being already secure) for experimental purposes (e.g., implementing another cipher) because it is written in Go, which is safe for cryptography.
	// However, for modifications made by copying, it only works for Go applications.
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
	// Note: This should be correct for TLS 1.3, and it's safe for Go due to the following reasons:
	//
	// - The Stream instance itself is designed to be thread-safe and can be safely shared across multiple goroutines.
	//   It does not maintain any mutable state that could cause race conditions or interference between goroutines.
	//
	// - Additionally, don't use QUIC connections (UDP) as they are not safe for multiple goroutines (see https://pkg.go.dev/crypto/tls@master#QUICConn) and might not suitable for ciphertext chacha (due it's udp).
	n, err := c.Conn.Read(b)
	if err != nil {
		return 0, err
	}
	var decryptedBuf bytes.Buffer
	err = c.Stream.Decrypt(bytes.NewReader(b[:n]), &decryptedBuf)
	if err != nil {
		return 0, err
	}
	copy(b, decryptedBuf.Bytes())
	return decryptedBuf.Len(), nil
}

// Write encrypts the provided data using the Stream and writes it to the TLS connection.
//
// TODO: Improve this function to support storing encrypted data/values from the client into a database (e.g sensitive data), which is a safer approach for storing encrypted data.
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
//
// Note: This is suitable due to TLS 1.3's improved handling of protocols (e.g., keys, handshake)
// compared to TLS 1.2, which is more complex and less efficient. It's no wonder TLS 1.2 is more susceptible to DoS attacks.
// However, this implementation is not yet complete as Go 1.23 has not been released.
func NewStreamConn(tlsConn *tls.Conn, stream *stream.Stream) net.Conn {
	return &streamConn{
		Conn:   tlsConn,
		Stream: stream,
	}
}

// Accept waits for and returns the next connection to the listener, wrapped in a streamConn.
func (l *streamListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// Peek into the connection to check if it's a browser request
	buf := make([]byte, stream.ChunkSize)
	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Check if the request is from a browser
	if isBrowserRequest(buf[:n]) {
		// Send a response indicating that the browser is unsupported
		//
		// Note: This is a hardcoded raw response, which may not be suitable depending on the client request (e.g., when the client sends JSON, the response must be in JSON format).
		// Since this packet is communicated through your network/modem router, it's possible but not worthwhile due to the complexity involved.
		response := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nUnsupported browser. Please use a compatible client.\r\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			conn.Close()
			return nil, err
		}
		conn.Close()
		return nil, errors.New("Boring TLS: unsupported browser")
	}

	// Reset the connection back to its original state
	conn = &rewoundConn{Conn: conn, buf: buf[:n]}

	tlsConn := tls.Server(conn, l.tlsConfig)
	return NewStreamConn(tlsConn, l.stream), nil
}

// NewStreamListener creates a new streamListener instance.
//
// Note: This for Server-Side, For clients that can be used for private communication, such as real-time chat or other features (e.g, authentication),
// the implementation must be done outside of this server due to the nature of the streamListener.
func NewStreamListener(listener net.Listener, tlsConfig *tls.Config, stream *stream.Stream) net.Listener {
	return &streamListener{
		Listener:  listener,
		tlsConfig: tlsConfig,
		stream:    stream,
	}
}

// rewoundConn is a net.Conn that allows reading from a buffer before reading from the underlying connection.
type rewoundConn struct {
	net.Conn
	buf       []byte
	bufCursor int
}

// Read reads data from the connection.
//
// It first checks if there is any buffered data available from the previous peek operation.
// If buffered data is available, it copies the data from the buffer to the provided byte slice
// and advances the buffer cursor accordingly. If the buffered data is exhausted, it proceeds
// to read from the underlying connection.
//
// Details:
//
//   - If there is buffered data available (c.bufCursor < len(c.buf)), the method copies the data
//     from the buffer starting at the current cursor position (c.buf[c.bufCursor:]) to the provided
//     byte slice 'b'. It then advances the buffer cursor by the number of bytes copied.
//
//   - If the buffered data is exhausted (c.bufCursor >= len(c.buf)), the method proceeds to read
//     from the underlying connection using c.Conn.Read(b). This allows seamless reading of data
//     from the connection once the buffered data has been consumed.
func (c *rewoundConn) Read(b []byte) (int, error) {
	if c.bufCursor < len(c.buf) {
		n := copy(b, c.buf[c.bufCursor:])
		c.bufCursor += n
		return n, nil
	}
	return c.Conn.Read(b)
}
