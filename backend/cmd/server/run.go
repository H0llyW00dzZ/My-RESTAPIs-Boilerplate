// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"

	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	handler "h0llyw00dz-template/backend/internal/server"
	"h0llyw00dz-template/env"
)

// main is the entry point for the application. It initializes the application
// by setting up the Fiber web server, configuring middleware, and registering routes.
// It relies on environment variables to customize the application's behavior,
// and it starts the server with graceful shutdown capabilities.
//
// Also, note that while the program is running and you can see the header "bound on host 0.0.0.0 and ...",
// the IP address "0.0.0.0" is not susceptible to exploits, attacks, or any other vulnerabilities.
// This is because "0.0.0.0" is my home. So be smart and refer to the source code and documentation to understand how it works you poggers.
func main() {
	appName, port, monitorPath, timeFormat, readTimeout, writeTimeout, shutdownTimeout := getEnvVariables()
	app := setupFiber(appName, readTimeout, writeTimeout)

	// Start the server with graceful shutdown and monitor
	startServer(app, appName, port, monitorPath, timeFormat, shutdownTimeout)
}

// getEnvVariables retrieves essential configuration settings from environment variables.
// It provides default values for the application name, port, monitoring path, time format, and timeouts
// to ensure the application has sensible defaults if environment variables are not set.
func getEnvVariables() (appName, port, monitorPath, timeFormat string, readTimeout, writeTimeout, shutdownTimeout time.Duration) {
	// Get the APP_NAME, PORT, and MONITOR_PATH from environment variables or use default values.
	appName = getEnv(env.APPNAME, "Gopher")
	port = getEnv(env.PORT, "8080")
	monitorPath = getEnv(env.MONITORPATH, "/monitor")
	// Get the TIME_FORMAT from environment variables or use default value
	// Note: List Time Format Available: unix,default
	timeFormat = getEnv(env.TIMEFORMAT, "unix")

	// Get the READ_TIMEOUT, WRITE_TIMEOUT, and SHUTDOWN_TIMEOUT from environment variables or use default values.
	// Note: These default timeout values (5 seconds) are set to help prevent potential deadlocks/hangs.
	readTimeoutStr := getEnv(env.READTIMEOUT, "5s")
	writeTimeoutStr := getEnv(env.SHUTDOWNTIMEOUT, "5s")
	shutdownTimeoutStr := getEnv(env.SHUTDOWNTIMEOUT, "5s")

	// Parse the timeout values into time.Duration
	readTimeout, _ = time.ParseDuration(readTimeoutStr)
	writeTimeout, _ = time.ParseDuration(writeTimeoutStr)
	shutdownTimeout, _ = time.ParseDuration(shutdownTimeoutStr)

	return
}

// setupFiber initializes a new Fiber application with custom configuration.
// It sets up the JSON encoder/decoder, case sensitivity, and strict routing,
// and applies the application name to the server headers.
func setupFiber(appName string, readTimeout, writeTimeout time.Duration) *fiber.App {
	return fiber.New(fiber.Config{
		ServerHeader: appName,
		AppName:      appName,
		// Note: Using the sonic JSON encoder/decoder provides better performance and is more memory-efficient
		// since Fiber is designed for zero allocation memory usage.
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
		CaseSensitive:    true,
		StrictRouting:    true,
		DisableKeepalive: false,
		ReadTimeout:      readTimeout,
		WriteTimeout:     writeTimeout,
		// Note: It's important to set Prefork to false because if it's enabled and running in Kubernetes,
		// it may get killed by an Out-of-Memory (OOM) error due to a conflict with the Horizontal Pod Autoscaler (HPA).
		Prefork: false,
		// Which is suitable for streaming AI Response.
		StreamRequestBody: true,
	})
}

// startServer configures and starts the Fiber web server.
// It initializes logging, determines the server address, and calls the server start function
// with graceful shutdown handling.
func startServer(app *fiber.App, appName, port, monitorPath, timeFormat string, shutdownTimeout time.Duration) {
	// Initialize the logger with the AppName from the environment variable
	log.InitializeLogger(app.Config().AppName, timeFormat)

	// Define server address
	addr := fmt.Sprintf(":%s", port) // Use the port from the environment variable

	// Create a new instance of FiberServer
	server := handler.NewFiberServer(app, appName, monitorPath)

	// Start the server with graceful shutdown and monitor
	//
	// TODO: Implement environment mode. For example, when the environment is set to "dev" or "local",
	// it will switch to "Listen (Non HTTPS)". Otherwise, it will force a change from Listen to ListenTLS
	// for public access that can be accessed by a browser. For the Go application itself (only accessed by the Go application, which is pretty useful for authentication), it will switch
	// to a combination of Listener and StreamListener (automatically and transparently encrypting and decrypting,
	// similar to Certificate Transparency my Boring TLS Certificate) to use TLS 1.3 protocols.
	//
	// Note: When running in Kubernetes, this is an easy configuration with cert-manager.io for environment mode (as currently implemented). It just uses a secret.
	handler.StartServer(server, addr, monitorPath, shutdownTimeout, nil, nil)
}

// getEnv reads an environment variable and returns its value.
// If the environment variable is not set, it returns a specified default value.
// This function encapsulates the standard library's [os.LookupEnv] to provide defaults,
// following the common Go idiom of "make the zero value useful".
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// TLSConfig creates and configures a TLS configuration for the server.
// It sets the minimum TLS version to TLS 1.3 and defines preferred curve preferences.
// It also uses the [fiber.GetClientInfo] function from the [fiber.TLSHandler] to get client information.
//
// Example Usage in startupserver:
//
//	 // Note: myCert, myCertKey load from environment, if it's running on Kubernetes it would be easy for secret and secure
//	gopherCert, err := tls.LoadX509KeyPair(myCert, myCertKey)
//	if err != nil {
//		// handle error
//	}
//	gopherTLSConfig := TLSConfig(gopherCert, nil) // clientCert nil
//	handler.StartServer(server, addr, monitorPath, shutdownTimeout, gopherTLSConfig, nil) // Boring TLS 1.3 nil due it's my own protocol and currently unavailable.
//
// Example ListenMutualTLS Usage in startupserver:
//
//	 // Note: myCert, myCertKey, clientCertFile load from environment, if it's running on Kubernetes it would be easy for secret and secure
//	gopherCertPool, err := tls.LoadX509KeyPair(myCert, myCertKey)
//	if err != nil {
//		// handle error
//	}
//
//	// Load client CA certificate (optional)
//	var clientCertPool *x509.CertPool
//	clientCertBytes, err := os.ReadFile(clientCertFile)
//
//	if err == nil {
//	    clientCertPool = x509.NewCertPool()
//	    clientCertPool.AppendCertsFromPEM(clientCertBytes)
//	}
//
//	gopherTLSConfig := TLSConfig(gopherCertPool, clientCertPool)
//	handler.StartServer(server, addr, monitorPath, shutdownTimeout, gopherTLSConfig, nil) // Boring TLS 1.3 nil due it's my own protocol and currently unavailable.
//
// Note: This design is well-written and idiomatic, unlike designs that spliting functions (e.g., those related to TLS like "ListenTLS" "ListenMutualTLS" or whatever it is).
func TLSConfig(cert tls.Certificate, leafCA, subCA, rootCA *x509.Certificate, clientCertPool *x509.CertPool) *tls.Config {
	tlsHandler := &fiber.TLSHandler{}
	// Note: Go's standard TLS 1.3 implementation does not allow direct configuration of cipher suites.
	// This means that while one can specify cipher suites in Go code, the implementation will prioritize the use of
	// AES-based ciphers like TLS_AES_128_GCM_SHA256 or TLS_AES_256_GCM_SHA384 (both bad common cipher, not even allowed to use ChaCha20 especially XChaCha20 which more secure),
	// which may be slower than ChaCha20 which is faster on some platforms.
	//
	// However, the actual cipher suite used in a TLS connection is determined by the client's preferences and capabilities.
	// When the client's preference is set to prioritize ChaCha20-based cipher suites like TLS_CHACHA20_POLY1305_SHA256,
	// and the server supports it, they will negotiate and agree to use that cipher suite for the encrypted communication.
	//
	// Example (Tested on Firefox Browser, which allows customization of cipher suites for TLS 1.3):
	// Network Tool: Wireshark Interface -> Link-Layer Header BSD Loopback
	//
	//  TLSv1.3 Record Layer: Handshake Protocol: Client Hello (SNI=localhost) (Client Browser)
	//  Cipher Suites (9 suites)
	//  Cipher Suite: TLS_CHACHA20_POLY1305_SHA256 (0x1303)
	//  Cipher Suite: TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256 (0xcca9)
	//  Cipher Suite: TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256 (0xcca8)
	//  Cipher Suite: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 (0xc02b)
	//  Cipher Suite: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 (0xc02f)
	//  Cipher Suite: TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 (0xc02c)
	//  Cipher Suite: TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 (0xc030)
	//  Cipher Suite: TLS_RSA_WITH_AES_128_GCM_SHA256 (0x009c)
	//  Cipher Suite: TLS_RSA_WITH_AES_256_GCM_SHA384 (0x009d)
	//
	// TLSv1.3 Record Layer: Handshake Protocol: Server Hello (This Repo)
	// Cipher Suite: TLS_CHACHA20_POLY1305_SHA256 (0x1303)
	// TLSv1.3 Record Layer: Change Cipher Spec Protocol: Change Cipher Spec
	//
	// In this example, the client (Firefox Browser) sends a Client Hello message with its supported cipher suites, prioritizing
	// TLS_CHACHA20_POLY1305_SHA256. The server responds with a Server Hello message, agreeing to use
	// TLS_CHACHA20_POLY1305_SHA256 based on the client's preferences.
	//
	// Acknowledgment:
	// ChaCha20 is known for its excellent performance, particularly on mobile devices and low-end processors.
	// By using ChaCha20, clients can potentially achieve better encryption and decryption speeds compared to
	// using AES-based ciphers, resulting in improved overall performance.
	tlsConfig := &tls.Config{
		// Note: This Explicitly setting the maximum and minimum TLS versions can improve the negotiation process.
		MaxVersion: tls.VersionTLS13, // Explicit
		MinVersion: tls.VersionTLS13,
		// Note: CurvePreferences works well when set, for example, with "tls.X25519" as the first preference.
		// However, when setting CipherSuites for TLS 1.3 or PreferServerCipherSuites (which is not actually deprecated),
		// it won't work properly if CipherSuites or PreferServerCipherSuites are specified because Go's standard TLS 1.3 implementation
		// does not allow direct configuration of cipher suites, even for the client side (e.g., http client in Go).
		// It's better to keep it like this, as it will depend on the client's preferences. For example,
		// when TLS_CHACHA20_POLY1305_SHA256 is set as the top/first preference, the server will choose "TLS_CHACHA20_POLY1305_SHA256" based on the client's preference.
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		Certificates:   []tls.Certificate{cert},
		GetCertificate: tlsHandler.GetClientInfo,
		// Note: This safe for multiple goroutines each time it is called, ensuring that each goroutine gets its own independent reader
		// The fixedReader itself does not maintain any mutable state, making it safe for concurrent use.
		Rand: rand.FixedSize32Bytes(),
		// TODO: Handle "VerifyPeerCertificate" for Certificate Transparency.
	}

	// Create a certificate pool (basically CA chains) for the CA certificates
	// Note: A correct implementation:
	// leafCA (first), subCA (second), rootCA (third)
	//
	// Bad Practice:
	// Using cat command for append it.
	caCertPool := x509.NewCertPool()
	if leafCA != nil {
		caCertPool.AddCert(leafCA)
	}
	if subCA != nil {
		caCertPool.AddCert(subCA)
	}
	if rootCA != nil {
		caCertPool.AddCert(rootCA)
	}

	// Set the RootCAs field in the TLS config
	tlsConfig.RootCAs = caCertPool

	// Only enable client auth if clientCertPool is not nil
	// TODO: Handle "GetClientCertificate" that might need.
	// Note: This different, it for mTLS, unlike "caCertPool" that for HTTPS
	if clientCertPool != nil {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = clientCertPool
	}

	return tlsConfig
}
