// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server_test

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
	"h0llyw00dz-template/backend/internal/server"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func tlsConfig(cert tls.Certificate) *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
		Certificates: []tls.Certificate{cert},
	}
}

func clientTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
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
		log.Println("Server: Received request")
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
		log.Println("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	log.Println("Client: Establishing TLS connection")
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
	log.Println("[Packet Netw0rkz] Client: Sending encrypted request")
	req := "GET /test HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"
	encryptedReq := &bytes.Buffer{}
	err = clientStream.Encrypt(bytes.NewReader([]byte(req)), encryptedReq)
	if err != nil {
		t.Fatal(err)
	}
	encryptedReqHex := hex.EncodeToString(encryptedReq.Bytes())
	log.Printf("[Packet Netw0rkz] Client: Encrypted request (hex): %s", encryptedReqHex)
	_, err = conn.Write(encryptedReq.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	log.Println("[Packet Netw0rkz] Client: Encrypted request sent")

	// Read the encrypted response from the server
	log.Println("[Packet Netw0rkz] Server: Reading encrypted response")
	var encryptedResp []byte
	buffer := make([]byte, 1024)
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
	log.Printf("[Packet Netw0rkz] Server: Encrypted response (hex): %s", encryptedRespHex)
	log.Println("[Packet Netw0rkz] Server: Encrypted response received")

	// Decrypt the response
	log.Println("[Packet Netw0rkz] Server: Decrypting response")
	decryptedResp := &bytes.Buffer{}
	err = clientStream.Decrypt(bytes.NewReader(encryptedResp), decryptedResp)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("[Packet Netw0rkz] Server: Response decrypted")

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
		log.Println("Server: Received request")
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
		log.Println("Server: Starting server")
		errChan <- app.Listener(streamListener)
	}()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Create a TLS client configuration
	tlsClientConfig := clientTLSConfig()

	// Create a TLS connection to the server
	log.Println("Client: Establishing TLS connection")
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
	log.Println("[Packet Netw0rkz] Client: Sending encrypted request")
	req := "GET /test HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"
	encryptedReq := &bytes.Buffer{}
	err = clientStream.Encrypt(bytes.NewReader([]byte(req)), encryptedReq)
	if err != nil {
		t.Fatal(err)
	}
	encryptedReqHex := hex.EncodeToString(encryptedReq.Bytes())
	log.Printf("[Packet Netw0rkz] Client: Encrypted request (hex): %s", encryptedReqHex)
	_, err = conn.Write(encryptedReq.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	log.Println("[Packet Netw0rkz] Client: Encrypted request sent")

	// Read the encrypted response from the server
	log.Println("[Packet Netw0rkz] Server: Reading encrypted response")
	var encryptedResp []byte
	buffer := make([]byte, 1024)
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
	log.Printf("[Packet Netw0rkz] Server: Encrypted response (hex): %s", encryptedRespHex)
	log.Println("[Packet Netw0rkz] Server: Encrypted response received")

	// Decrypt the response
	log.Println("[Packet Netw0rkz] Server: Decrypting response")
	decryptedResp := &bytes.Buffer{}
	err = clientStream.Decrypt(bytes.NewReader(encryptedResp), decryptedResp)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("[Packet Netw0rkz] Server: Response decrypted")

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

	log.Printf("[Packet Netw0rkz] Boring TLS: Decrypted response: %s", decryptedResp.String())

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
		log.Println("Server: Received request")
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
		log.Println("Server: Starting server")
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
	log.Println("Client: Establishing TLS connection")
	_, err = tls.Dial("tcp", "localhost:8082", tlsClientConfig)
	if err == nil {
		t.Fatal("Expected TLS handshake to fail due to wrong protocol version")
	}
	log.Printf("Client: TLS handshake failed as expected: %v", err)

	// Check if the server returned an error
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}
