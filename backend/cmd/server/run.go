// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

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
	handler "h0llyw00dz-template/backend/internal/server"
)

// main is the entry point for the application. It initializes the application
// by setting up the Fiber web server, configuring middleware, and registering routes.
// It relies on environment variables to customize the application's behavior,
// and it starts the server with graceful shutdown capabilities.
func main() {
	appName, port, monitorPath, timeFormat, readTimeout, writeTimeout, shutdownTimeout := getEnvVariables()
	app := setupFiber(appName, readTimeout, writeTimeout)

	// Start the server with graceful shutdown and monitor
	startServer(app, appName, port, monitorPath, timeFormat, shutdownTimeout)
}

// getEnvVariables retrieves essential configuration settings from environment variables.
// It provides default values for the application name, port, monitoring path, time format, and timeouts
// to ensure the application has sensible defaults if environment variables are not set.
//
// The following environment variables are used:
//
//   - APP_NAME: The name of the application (default: "Gopher").
//
//   - PORT: The port number on which the server will listen (default: "8080").
//
//   - MONITOR_PATH: The path for the server monitoring endpoint (default: "/monitor").
//
//   - TIME_FORMAT: The format for logging timestamps (default: "unix").
//
//     Available options:
//
//   - "unix": Unix timestamp format (e.g., [1713355079]).
//
//   - "default": Default timestamp format (e.g., 2024/04/17 15:04:05).
//
//   - READ_TIMEOUT: The maximum duration for reading the entire request, including the body (default: "5s").
//
//   - WRITE_TIMEOUT: The maximum duration before timing out writes of the response (default: "5s").
//
//   - SHUTDOWN_TIMEOUT: The maximum duration to wait for active connections to finish during server shutdown (default: "5s").
func getEnvVariables() (appName, port, monitorPath, timeFormat string, readTimeout, writeTimeout, shutdownTimeout time.Duration) {
	// Get the APP_NAME, PORT, and MONITOR_PATH from environment variables or use default values.
	appName = getEnv("APP_NAME", "Gopher")
	port = getEnv("PORT", "8080")
	monitorPath = getEnv("MONITOR_PATH", "/monitor")
	// Get the TIME_FORMAT from environment variables or use default value
	// Note: List Time Format Available: unix,default
	timeFormat = getEnv("TIME_FORMAT", "unix")

	// Get the READ_TIMEOUT, WRITE_TIMEOUT, and SHUTDOWN_TIMEOUT from environment variables or use default values.
	// Note: These default timeout values (5 seconds) are set to help prevent potential deadlocks/hangs.
	readTimeoutStr := getEnv("READ_TIMEOUT", "5s")
	writeTimeoutStr := getEnv("WRITE_TIMEOUT", "5s")
	shutdownTimeoutStr := getEnv("SHUTDOWN_TIMEOUT", "5s")

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
// This function encapsulates the standard library's os.LookupEnv to provide defaults,
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
func TLSConfig(cert tls.Certificate, clientCertPool *x509.CertPool) *tls.Config {
	tlsHandler := &fiber.TLSHandler{}
	// Note: Go's standard TLS 1.3 implementation does not allow direct configuration of cipher suites.
	// This means that while one can specify cipher suites in Go code, the implementation will prioritize the use of
	// AES-based ciphers like TLS_AES_128_GCM_SHA256 or TLS_AES_256_GCM_SHA384 (both bad common cipher, not even allowed to use ChaCha20 especially XChaCha20 which more secure),
	// which may be slower than ChaCha20 which is faster on some platforms.
	tlsConfig := &tls.Config{
		MaxVersion: tls.VersionTLS13, // Explicit
		MinVersion: tls.VersionTLS13,
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
		Rand: handler.RandTLS(),
		// TODO: Handle "VerifyPeerCertificate" for Certificate Transparency.
	}

	// Only enable client auth if clientCertPool is not nil
	// TODO: Handle "GetClientCertificate" that might need.
	if clientCertPool != nil {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = clientCertPool
	}

	return tlsConfig
}
