// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
	"h0llyw00dz-template/backend/internal/server"
	"io"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func tlsConfig(cert tls.Certificate) *tls.Config {
	log.InitializeLogger("Boring TLS 1.3 Testing", "")
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			// Note: These are classical elliptic curves for TLS 1.3 key exchange.
			// For experimental purposes related to post-quantum hybrid design, refer to:
			// https://datatracker.ietf.org/doc/html/draft-ietf-tls-hybrid-design-10
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		Certificates: []tls.Certificate{cert},
	}
}

func clientTLSConfig() *tls.Config {
	log.InitializeLogger("Boring TLS 1.3 Testing", "")
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			// Note: These are classical elliptic curves for TLS 1.3 key exchange.
			// For experimental purposes related to post-quantum hybrid design, refer to:
			// https://datatracker.ietf.org/doc/html/draft-ietf-tls-hybrid-design-10
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		InsecureSkipVerify: true,
		ServerName:         "localhost",
	}
}

// Note: This is just a test that demonstrates a working example of using TLS 1.3 along with an additional encryption layer.
// It is still unfinished. If finished, it would require writing a lot of functions when using a custom cipher for the cipher suite (might be copied from std/dependency injection).
func TestStreamServer(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	// Note: This test kinda slow (tested on windows) due 2 cipher text, if pure ChaCha20-Poly1305 might faster
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		log.LogInfo("Server: Received request")
		return c.SendString("Hello, World!")
	})

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a listener
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}

	// Wrap the listener with streamListener
	streamListener := server.NewStreamListener(listener, tlsServerConfig, s)

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		log.LogInfo("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	log.LogInfo("Client: Establishing TLS connection")
	conn, err := tls.Dial("tcp", "localhost:8080", tlsClientConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Create a new Stream instance for the client
	clientStream, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Send an encrypted request to the server
	log.LogInfo("[Packet Netw0rkz] Client: Sending encrypted request")
	req := "GET /test HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"
	encryptedReq := &bytes.Buffer{}
	err = clientStream.Encrypt(bytes.NewReader([]byte(req)), encryptedReq)
	if err != nil {
		t.Fatal(err)
	}
	encryptedReqHex := hex.EncodeToString(encryptedReq.Bytes())
	log.LogInfof("[Packet Netw0rkz] Client: Encrypted request (hex): %s", encryptedReqHex)
	_, err = conn.Write(encryptedReq.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Client: Encrypted request sent")

	// Read the encrypted response from the server
	log.LogInfo("[Packet Netw0rkz] Server: Reading encrypted response")
	var encryptedResp []byte
	buffer := make([]byte, stream.ChunkSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			t.Fatal(err)
		}
		encryptedResp = append(encryptedResp, buffer[:n]...)
		if n < len(buffer) {
			break
		}
	}
	encryptedRespHex := hex.EncodeToString(encryptedResp)
	log.LogInfof("[Packet Netw0rkz] Server: Encrypted response (hex): %s", encryptedRespHex)
	log.LogInfo("[Packet Netw0rkz] Server: Encrypted response received")

	// Decrypt the response
	log.LogInfo("[Packet Netw0rkz] Server: Decrypting response")
	decryptedResp := &bytes.Buffer{}
	err = clientStream.Decrypt(bytes.NewReader(encryptedResp), decryptedResp)
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Server: Response decrypted")

	// Check the decrypted response
	expectedHeaders := []string{
		"HTTP/1.1 200 OK",
		"Content-Type: text/plain; charset=utf-8",
		"Content-Length: 13",
	}
	expectedBody := "Hello, World!"

	respLines := strings.Split(decryptedResp.String(), "\r\n")
	for _, header := range expectedHeaders {
		if !contains(respLines, header) {
			t.Errorf("missing expected header: %q", header)
		}
	}

	if !contains(respLines, expectedBody) {
		t.Errorf("missing expected body: %q", expectedBody)
	}

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

func contains(lines []string, target string) bool {
	for _, line := range lines {
		if line == target {
			return true
		}
	}
	return false
}

// Note: Explicit HTTPS "Content-Length" it's nil
func TestStreamServerExplicitHTTPS(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	// Note: This test kinda slow (tested on windows) due 2 cipher text, if pure ChaCha20-Poly1305 might faster
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		log.LogInfo("Server: Received request")
		if c.Protocol() == "https" {
			return c.SendString("Hello, World! (via TLS)")
		}
		return c.SendString("Hello, World!")
	})

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a listener
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}

	// Wrap the listener with streamListener
	streamListener := server.NewStreamListener(listener, tlsServerConfig, s)

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		log.LogInfo("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	log.LogInfo("Client: Establishing TLS connection")
	conn, err := tls.Dial("tcp", "localhost:8081", tlsClientConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Create a new Stream instance for the client
	clientStream, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Send an encrypted request to the server
	log.LogInfo("[Packet Netw0rkz] Client: Sending encrypted request")
	req := "GET /test HTTP/1.1\r\nHost: localhost:8081\r\n\r\n"
	encryptedReq := &bytes.Buffer{}
	err = clientStream.Encrypt(bytes.NewReader([]byte(req)), encryptedReq)
	if err != nil {
		t.Fatal(err)
	}
	encryptedReqHex := hex.EncodeToString(encryptedReq.Bytes())
	log.LogInfof("[Packet Netw0rkz] Client: Encrypted request (hex): %s", encryptedReqHex)
	_, err = conn.Write(encryptedReq.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Client: Encrypted request sent")

	// Read the encrypted response from the server
	log.LogInfo("[Packet Netw0rkz] Server: Reading encrypted response")
	var encryptedResp []byte
	buffer := make([]byte, stream.ChunkSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			t.Fatal(err)
		}
		encryptedResp = append(encryptedResp, buffer[:n]...)
		if n < len(buffer) {
			break
		}
	}
	encryptedRespHex := hex.EncodeToString(encryptedResp)
	log.LogInfof("[Packet Netw0rkz] Server: Encrypted response (hex): %s", encryptedRespHex)
	log.LogInfo("[Packet Netw0rkz] Server: Encrypted response received")

	// Decrypt the response
	log.LogInfo("[Packet Netw0rkz] Server: Decrypting response")
	decryptedResp := &bytes.Buffer{}
	err = clientStream.Decrypt(bytes.NewReader(encryptedResp), decryptedResp)
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Server: Response decrypted")

	// Check the decrypted response
	expectedHeaders := []string{
		"HTTP/1.1 200 OK",
		"Content-Type: text/plain; charset=utf-8",
	}
	expectedBody := "Hello, World! (via TLS)"

	respLines := strings.Split(decryptedResp.String(), "\r\n")
	for _, header := range expectedHeaders {
		if !contains(respLines, header) {
			t.Errorf("missing expected header: %q", header)
		}
	}

	if !contains(respLines, expectedBody) {
		t.Errorf("missing expected body: %q", expectedBody)
	}

	log.LogInfof("[Packet Netw0rkz] Boring TLS: Decrypted response: %s", decryptedResp.String())

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

func TestStreamClientWrongProtocol(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		log.LogInfo("Server: Received request")
		return c.SendString("Hello, World!")
	})

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a listener
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		t.Fatal(err)
	}

	// Wrap the listener with streamListener
	streamListener := server.NewStreamListener(listener, tlsServerConfig, s)

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		log.LogInfo("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TLS client configuration with TLS 1.2
	tlsClientConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
		InsecureSkipVerify: true,
		ServerName:         "localhost",
	}

	// Create a TLS connection to the server
	log.LogInfo("Client: Establishing TLS connection")
	_, err = tls.Dial("tcp", "localhost:8082", tlsClientConfig)
	if err == nil {
		t.Fatal("Expected TLS handshake to fail due to wrong protocol version")
	}
	log.LogInfof("Client: TLS handshake failed as expected: %v", err)

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

func TestStreamServerStupidMiddleman(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New(
		fiber.Config{},
	)

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		log.LogInfo("Server: Received request")
		return c.SendString("Hello, World!")
	})

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a listener
	listener, err := net.Listen("tcp", ":8083")
	if err != nil {
		t.Fatal(err)
	}

	// Wrap the listener with streamListener
	streamListener := server.NewStreamListener(listener, tlsServerConfig, s)

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		log.LogInfo("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	log.LogInfo("Client: Establishing TLS connection")
	conn, err := tls.Dial("tcp", "localhost:8083", tlsClientConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Send a plain request to the server
	log.LogInfo("[Packet Netw0rkz] Client: Sending plain request")
	req := "GET /test HTTP/1.1\r\nHost: localhost:8083\r\n\r\n"
	_, err = conn.Write([]byte(req))
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Client: Plain request sent")

	// Simulate a stupid middleman intercepting the plain request
	log.LogInfo("[Packet Netw0rkz] Middleman: Intercepted plain request")
	log.LogInfof("[Packet Netw0rkz] Middleman: Plain request: %s", req)

	// Read the response from the server
	log.LogInfo("[Packet Netw0rkz] Server: Reading response")
	var resp []byte
	buffer := make([]byte, stream.ChunkSize)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}
	resp = append(resp, buffer[:n]...)
	log.LogInfof("[Packet Netw0rkz] Server: Encrypted response: %x", resp)
	log.LogInfo("[Packet Netw0rkz] Server: Encrypted response received")

	// Simulate a stupid middleman intercepting the encrypted response
	log.LogInfo("[Packet Netw0rkz] Middleman: Intercepted encrypted response")
	log.LogInfof("[Packet Netw0rkz] Middleman: Encrypted response: %x", resp)

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

func TestStreamServerExplicitHTTPSUnixPacket(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix systems")
	}

	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	// Note: This test kinda slow (tested on windows) due 2 cipher text, if pure ChaCha20-Poly1305 might faster
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		log.LogInfo("Server: Received request")
		if c.Protocol() == "https" {
			return c.SendString("Hello, Unix! (via TLS)")
		}
		return c.SendString("Hello, Unix!")
	})

	// Create a Unix domain socket path
	socketPath := "/tmp/test.sock"

	// Create a listener using Unix domain socket
	listener, err := net.Listen("unixpacket", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// Wrap the listener with streamListener
	streamListener := server.NewStreamListener(listener, nil, s)

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		log.LogInfo("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a Unix domain socket connection to the server
	log.LogInfo("Client: Establishing Unix domain socket connection")
	conn, err := net.Dial("unixpacket", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Create a new Stream instance for the client
	clientStream, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Send an encrypted request to the server
	log.LogInfo("[Packet Netw0rkz] Client: Sending encrypted request")
	req := "GET /test HTTP/1.1\r\nHost: localhost\r\n\r\n"
	encryptedReq := &bytes.Buffer{}
	err = clientStream.Encrypt(bytes.NewReader([]byte(req)), encryptedReq)
	if err != nil {
		t.Fatal(err)
	}
	encryptedReqHex := hex.EncodeToString(encryptedReq.Bytes())
	log.LogInfof("[Packet Netw0rkz] Client: Encrypted request (hex): %s", encryptedReqHex)
	_, err = conn.Write(encryptedReq.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Client: Encrypted request sent")

	// Read the encrypted response from the server
	log.LogInfo("[Packet Netw0rkz] Server: Reading encrypted response")
	var encryptedResp []byte
	buffer := make([]byte, stream.ChunkSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			t.Fatal(err)
		}
		encryptedResp = append(encryptedResp, buffer[:n]...)
		if n < len(buffer) {
			break
		}
	}
	encryptedRespHex := hex.EncodeToString(encryptedResp)
	log.LogInfof("[Packet Netw0rkz] Server: Encrypted response (hex): %s", encryptedRespHex)
	log.LogInfo("[Packet Netw0rkz] Server: Encrypted response received")

	// Decrypt the response
	log.LogInfo("[Packet Netw0rkz] Server: Decrypting response")
	decryptedResp := &bytes.Buffer{}
	err = clientStream.Decrypt(bytes.NewReader(encryptedResp), decryptedResp)
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Server: Response decrypted")

	// Check the decrypted response
	expectedHeaders := []string{
		"HTTP/1.1 200 OK",
		"Content-Type: text/plain; charset=utf-8",
	}
	expectedBody := "Hello, Unix! (via TLS)" // forgot

	respLines := strings.Split(decryptedResp.String(), "\r\n")
	for _, header := range expectedHeaders {
		if !contains(respLines, header) {
			t.Errorf("missing expected header: %q", header)
		}
	}

	if !contains(respLines, expectedBody) {
		t.Errorf("missing expected body: %q", expectedBody)
	}

	log.LogInfof("[Packet Netw0rkz] Boring TLS: Decrypted response: %s", decryptedResp.String())

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

func TestStreamConnDeadlines(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a listener
	listener, err := net.Listen("tcp", ":8084")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- err
			return
		}
		defer conn.Close()

		// Establish a TLS connection
		tlsConn := tls.Server(conn, tlsServerConfig)
		defer tlsConn.Close()

		// Create a new streamConn instance
		streamConn := server.NewStreamConn(tlsConn, s)

		// Set the read deadline to 1 second from now
		readDeadline := time.Now().Add(time.Second)
		if err := streamConn.SetReadDeadline(readDeadline); err != nil {
			errChan <- err
			return
		}

		// Read from the connection
		buffer := make([]byte, stream.ChunkSize)
		_, err = streamConn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				errChan <- nil // Read timeout occurred as expected
			} else {
				errChan <- err
			}
			return
		}

		// Set the write deadline to 1 second from now
		writeDeadline := time.Now().Add(time.Second)
		if err := streamConn.SetWriteDeadline(writeDeadline); err != nil {
			errChan <- err
			return
		}

		// Write to the connection
		_, err = streamConn.Write([]byte("Hello, World!"))
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				errChan <- nil // Write timeout occurred as expected
			} else {
				errChan <- err
			}
			return
		}

		errChan <- nil
	}()

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	conn, err := tls.Dial("tcp", "localhost:8084", tlsClientConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Set the overall deadline to 2 seconds from now
	deadline := time.Now().Add(2 * time.Second)
	if err := conn.SetDeadline(deadline); err != nil {
		t.Fatal(err)
	}

	// Send an encrypted request to the server
	log.LogInfo("[Packet Netw0rkz] Client: Sending encrypted request")
	req := "GET /test HTTP/1.1\r\nHost: localhost:8084\r\n\r\n"
	encryptedReq := &bytes.Buffer{}
	err = s.Encrypt(bytes.NewReader([]byte(req)), encryptedReq)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write(encryptedReq.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Client: Encrypted request sent")

	// Wait for the read deadline to expire
	time.Sleep(2 * time.Second)

	// Read the encrypted response from the server
	log.LogInfo("[Packet Netw0rkz] Server: Reading encrypted response")

	buffer := make([]byte, stream.ChunkSize)
	_, err = conn.Read(buffer)
	if err == nil {
		t.Fatal("Expected read timeout error")
	}
	if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
		t.Fatalf("Unexpected error: %v", err)
	}
	log.LogInfo("[Packet Netw0rkz] Server: Read timeout occurred as expected")

	// Wait for the server to finish
	err = <-errChan
	if err != nil {
		t.Fatal(err)
	}

	// Set the write deadline to the past
	pastDeadline := time.Now().Add(-time.Second)
	if err := conn.SetWriteDeadline(pastDeadline); err != nil {
		t.Fatal(err)
	}

	// Attempt to write to the connection
	_, err = conn.Write([]byte("Another message"))
	if err == nil {
		t.Fatal("Expected write timeout error")
	}
	if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
		t.Fatalf("Unexpected error: %v", err)
	}
	log.LogInfo("[Packet Netw0rkz] Client: Write timeout occurred as expected")
}

func TestStreamConnSetDeadline(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a listener
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Start the server
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- err
			return
		}
		defer conn.Close()

		// Establish a TLS connection
		tlsConn := tls.Server(conn, tlsServerConfig)
		defer tlsConn.Close()

		// Create a new streamConn instance
		streamConn := server.NewStreamConn(tlsConn, s)

		// Set the overall deadline to 2 seconds from now
		deadline := time.Now().Add(2 * time.Second)
		if err := streamConn.SetDeadline(deadline); err != nil {
			errChan <- err
			return
		}

		// Read from the connection
		buffer := make([]byte, stream.ChunkSize)
		_, err = streamConn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				errChan <- nil // Overall deadline occurred as expected
			} else {
				errChan <- err
			}
			return
		}

		// Write to the connection
		_, err = streamConn.Write([]byte("Hello, World!"))
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				errChan <- nil // Overall deadline occurred as expected
			} else {
				errChan <- err
			}
			return
		}

		errChan <- nil
	}()

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	conn, err := tls.Dial("tcp", "localhost:8085", tlsClientConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Set the overall deadline to 1 second from now
	deadline := time.Now().Add(time.Second)
	if err := conn.SetDeadline(deadline); err != nil {
		t.Fatal(err)
	}

	// Send an encrypted request to the server
	log.LogInfo("[Packet Netw0rkz] Client: Sending encrypted request")
	req := "GET /test HTTP/1.1\r\nHost: localhost:8085\r\n\r\n"
	encryptedReq := &bytes.Buffer{}
	err = s.Encrypt(bytes.NewReader([]byte(req)), encryptedReq)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write(encryptedReq.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	log.LogInfo("[Packet Netw0rkz] Client: Encrypted request sent")

	// Wait for the overall deadline to expire
	time.Sleep(2 * time.Second)

	// Read the encrypted response from the server
	log.LogInfo("[Packet Netw0rkz] Server: Reading encrypted response")
	buffer := make([]byte, stream.ChunkSize)
	_, err = conn.Read(buffer)
	if err == nil {
		t.Fatal("Expected overall deadline error")
	}
	if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
		t.Fatalf("Unexpected error: %v", err)
	}
	log.LogInfo("[Packet Netw0rkz] Server: Overall deadline occurred as expected")

	// Wait for the server to finish
	err = <-errChan
	if err != nil {
		t.Fatal(err)
	}
}

func TestStreamServerWithoutAdditionalEncrypt(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	// Note: This test is kind of slow (tested on Windows) due to 2 ciphertexts; if pure ChaCha20-Poly1305 is used, it might be faster
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		log.LogInfo("Server: Received request")
		if c.Secure() {
			return c.JSON(fiber.Map{
				"message": "Hello, World! (via TLS)",
			})
		}
		return c.SendString("Hello, World!")
	})

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a listener
	listener, err := net.Listen("tcp", ":8086")
	if err != nil {
		t.Fatal(err)
	}

	// Wrap the listener with streamListener
	streamListener := server.NewStreamListener(listener, tlsServerConfig, s)

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		log.LogInfo("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	log.LogInfo("Client: Establishing TLS connection")
	conn, err := tls.Dial("tcp", "localhost:8086", tlsClientConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Wrap the TLS connection with streamConn
	streamConn := server.NewStreamConn(conn, s)

	// Create the plain HTTP request
	plainReq := "GET /test HTTP/1.1\r\nHost: localhost:8086\r\n\r\n"

	// Send the plain request to the server (it will be automatically encrypted)
	log.LogInfo("[Packet Netw0rkz] Client: Sending request")
	_, err = streamConn.Write([]byte(plainReq))
	if err != nil {
		t.Fatal(err)
	}

	log.LogInfo("[Packet Netw0rkz] Client: Request sent")

	// Read the response from the server (it will be automatically decrypted)
	log.LogInfo("[Packet Netw0rkz] Server: Reading response")
	var response []byte
	buffer := make([]byte, stream.ChunkSize)
	n, err := streamConn.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}
	response = append(response, buffer[:n]...)

	log.LogInfof("[Packet Netw0rkz] Boring TLS: Decrypted response: %s", string(response))

	// Extract the JSON body from the response
	responseLines := strings.Split(string(response), "\r\n")
	var jsonBody string
	for i, line := range responseLines {
		if line == "" {
			jsonBody = strings.Join(responseLines[i+1:], "\r\n")
			break
		}
	}

	// Check if the JSON body contains the expected message
	expectedMessage := "Hello, World! (via TLS)"
	var responseJSON fiber.Map
	err = sonic.Unmarshal([]byte(jsonBody), &responseJSON)
	if err != nil {
		t.Fatal(err)
	}
	if message, ok := responseJSON["message"]; ok {
		if message != expectedMessage {
			t.Errorf("Expected response message to be '%s', but got '%s'", expectedMessage, message)
		}
	} else {
		t.Error("Response JSON does not contain the 'message' field")
	}

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

// Note: The speed is not bad; however, this is not fully implemented.
func TestStreamServerWithCustomTransport(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	// Note: When this is bound into a TLS connection by modifying from copies the standard library,
	// for example, adding the stream encrypter/decrypter along with stream.New for the cipher suites,
	// it can create a high-level secure TLS connection without relying on quantum-resistant algorithm approaches,
	// because the most important about crytographic it's ciphertext and TLS/DTLS for transmitting.
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		if c.Secure() {
			// Note: This works well for automatic encryption/decryption during transport transparently.
			// However, do not try this on front-end apps such as browsers,
			// as it may not be compatible due to the specific cipher used and protocols. If it's still a Go application, it is compatible and works well (e.g., keys, handshake).
			// Even with TLS 1.3, not all browsers will work if used for HTTPS front-end, even on Firefox (in Firefox, it works; however, it only encrypts and is unable to decrypt), due to the cipher.
			log.LogInfof("Server: Received request %s", c.Body()) // Decrypted transparently...
			log.LogInfo("Server: Encrypting response")            // Encrypting transparently...
			return c.JSON(fiber.Map{
				"message": "Hello, World! (via TLS)",
			})
		}
		return c.SendString("Hello, World!")
	})

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a regular TCP listener
	ln, err := net.Listen("tcp", ":8087")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	// Wrap the TCP listener with the streamListener
	streamLn := server.NewStreamListener(ln, tlsServerConfig, s)

	// Start the server with the streamListener
	go app.Listener(streamLn)

	// Define different curve preferences for each transport
	curvePreferences := [][]tls.CurveID{
		{tls.CurveP384, tls.CurveP521, tls.X25519, tls.CurveP256},
		{tls.CurveP521, tls.X25519, tls.CurveP256, tls.CurveP384},
		{tls.X25519, tls.CurveP384, tls.CurveP521, tls.CurveP256},
	}

	// Create multiple custom transports with different curve preferences for Boring TLS 1.3 protocol
	// Note: This is suitable for Go applications; however, do not try it in a browser as it may not be compatible due to the specific cipher used and protocols.
	// If it's still a Go application, it is compatible and works well (e.g., keys, handshake).
	transports := make([]*http.Transport, len(curvePreferences))
	for i, curves := range curvePreferences {
		transports[i] = &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := &net.Dialer{}
				conn, err := dialer.DialContext(context.Background(), network, addr)
				if err != nil {
					return nil, err
				}

				tlsConn := tls.Client(conn, &tls.Config{
					MinVersion:         tls.VersionTLS13,
					InsecureSkipVerify: true,
					ServerName:         "localhost",
					CurvePreferences:   curves,
				})

				if err := tlsConn.HandshakeContext(context.Background()); err != nil {
					conn.Close()
					return nil, err
				}

				log.LogInfo(fmt.Sprintf("Client %d: Established TLS connection", i+1))
				return server.NewStreamConn(tlsConn, s), nil
			},
		}
	}

	// Create multiple clients with different transports
	clients := make([]*http.Client, len(transports))
	for i, transport := range transports {
		clients[i] = &http.Client{
			Transport: transport,
		}
	}

	// Make requests to the server using each client
	for i, client := range clients {
		log.LogInfo(fmt.Sprintf("Client %d: Sending request", i+1))

		// Create a request with a body
		requestBody := []byte(fmt.Sprintf("Request body from client %d", i+1)) // Encrypting transparently...
		req, err := http.NewRequest("GET", "https://"+ln.Addr().String()+"/test", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// Check the response
		// When the client reaches this point, the response is automatically decrypted transparently,
		// just like when the server reaches "c.Secure" during the transport. So The packet is encrypted during transmission.
		expectedBody := `{"message":"Hello, World! (via TLS)"}`
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		log.LogInfo(fmt.Sprintf("Client %d: Received response", i+1))
		log.LogInfo(fmt.Sprintf("Client %d: Decrypting response", i+1)) // Decrypted transparently...

		if string(body) != expectedBody {
			// If the response doesn't match the expected response when the server sends "{"message":"Hello, World! (via TLS)"}",
			// then the client is unable to decrypt it even transparently (similar to when testing in Firefox, where the browser is unable to decrypt the response for the client).
			t.Errorf("Expected response body to be '%s', but got '%s'", expectedBody, string(body))
		}
	}
}

func TestUnsupportedBrowserRequest(t *testing.T) {
	// Generate AES key and ChaCha20 key
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rand.Read(chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Stream instance
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	// Create a TLS configuration for the server
	tlsServerConfig := tlsConfig(cert)

	// Create a listener
	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		t.Fatal(err)
	}

	// Wrap the listener with streamListener
	streamListener := server.NewStreamListener(listener, tlsServerConfig, s)

	// Create a channel to receive the server error
	errChan := make(chan error)

	// Start the server
	go func() {
		log.LogInfo("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TCP connection to the server
	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Send an HTTP request simulating a browser
	req := "GET / HTTP/1.1\r\nHost: localhost:8080\r\nUser-Agent: GopherBrowser/5.0\r\n\r\n"
	_, err = conn.Write([]byte(req))
	if err != nil {
		t.Fatal(err)
	}

	// Read the response from the server
	var resp []byte
	buffer := make([]byte, stream.ChunkSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		resp = append(resp, buffer[:n]...)
	}

	// Check the response
	expectedResponse := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nUnsupported browser. Please use a compatible client.\r\n"
	if string(resp) != expectedResponse {
		t.Errorf("Unexpected response. Expected:\n%s\nGot:\n%s", expectedResponse, string(resp))
	}

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

func TestPipeStreamConn(t *testing.T) {
	// Generate test keys
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)

	// Create a Stream instance
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream: %v", err)
	}

	// Create a test TLS connection
	clientConn, serverConn := net.Pipe()
	tlsClientConn := tls.Client(clientConn, clientTLSConfig())

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	tlsServerConn := tls.Server(serverConn, tlsConfig(cert))

	// Create a streamConn instance for the client
	clientStreamConn := server.NewStreamConn(tlsClientConn, s)

	// Create a streamConn instance for the server
	serverStreamConn := server.NewStreamConn(tlsServerConn, s)

	// Test data
	testData := []byte("Hello, World!")

	// Write test data to the server streamConn
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := serverStreamConn.Write(testData)
		if err != nil {
			t.Errorf("Failed to write data to server streamConn: %v", err)
		}
		serverStreamConn.Close()
	}()

	// Read the decrypted data from the client streamConn
	decryptedData := make([]byte, stream.ChunkSize)
	n, err := clientStreamConn.Read(decryptedData)
	if err != nil {
		t.Fatalf("Failed to read decrypted data from client streamConn: %v", err)
	}
	decryptedData = decryptedData[:n]

	// Log the decrypted data
	log.LogInfof("Decrypted data: %s", decryptedData)

	// Compare the decrypted data with the original test data
	if !bytes.Equal(decryptedData, testData) {
		t.Errorf("Decrypted data does not match the original data")
	}

	// Wait for the goroutine to finish
	wg.Wait()
}

func TestPipeStreamOutside(t *testing.T) {
	// Generate test keys
	aesKey := make([]byte, 32)
	chachaKey := make([]byte, 32)

	// Create a Stream instance
	s, err := stream.New(aesKey, chachaKey)
	if err != nil {
		t.Fatalf("Failed to create Stream: %v", err)
	}

	// Create a test TLS connection
	clientConn, serverConn := net.Pipe()
	tlsClientConn := tls.Client(clientConn, clientTLSConfig())

	// Load the self-signed certificate and key
	cert, err := tls.LoadX509KeyPair("boring-cert.pem", "boring-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	tlsServerConn := tls.Server(serverConn, tlsConfig(cert))

	// Create a streamConn instance for the server
	serverStreamConn := server.NewStreamConn(tlsServerConn, s)

	// Test data
	testData := []byte("Hello, World!")

	// Write test data to the server streamConn
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := serverStreamConn.Write(testData)
		if err != nil {
			t.Errorf("Failed to write data to server streamConn: %v", err)
		}
		serverStreamConn.Close()
	}()

	// Read the encrypted data from the client connection
	encryptedData := make([]byte, stream.ChunkSize)
	n, err := tlsClientConn.Read(encryptedData)
	if err != nil {
		t.Fatalf("Failed to read encrypted data from client connection: %v", err)
	}
	encryptedData = encryptedData[:n]

	// Log the encrypted data
	log.LogInfof("Encrypted data: %x", encryptedData)

	// Wait for the goroutine to finish
	wg.Wait()
}
