// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package main

import (
	"fmt"
	setupTLS "h0llyw00dz-template/backend/internal/middleware/authentication/crypto/tls"
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
//
// Also, note that while the program is running and you can see the header "bound on host 0.0.0.0 and ...",
// the IP address "0.0.0.0" is not susceptible to exploits, attacks, or any other vulnerabilities.
// This is because "0.0.0.0" is my home. So be smart and refer to the source code and documentation to understand how it works you poggers.
func main() {
	// Retrieve configuration from environment variables
	config := getConfig()

	// Set up the Fiber application with the retrieved configuration
	app := setupFiber(config)

	// Start the server with the configured settings
	startServer(app, config)
}

// setupFiber initializes a new Fiber application with custom configuration.
// It sets up the JSON encoder/decoder, case sensitivity, and strict routing,
// and applies the application name to the server headers.
func setupFiber(config Config) *fiber.App {
	// TODO: Implement a server startup message mechanism similar to "Fiber" ASCII art,
	// with animation (e.g., similar to a streaming/bubble tea spinner) for multiple sites or large codebases.
	// The current static "Fiber" ASCII art only shows one site when there are multiple, which isn't ideal.
	// However, animated ASCII art may not be necessary right now, as it only works properly in terminals.
	return fiber.New(fiber.Config{
		ServerHeader: config.AppName,
		AppName:      config.AppName,
		// Note: Using the sonic JSON encoder/decoder provides better performance and is more memory-efficient
		// since Fiber is designed for zero allocation memory usage.
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
		CaseSensitive:    true,
		StrictRouting:    true,
		DisableKeepalive: false,
		ReadTimeout:      config.ReadTimeout,
		WriteTimeout:     config.WriteTimeout,
		// Note: It's important to set Prefork to false because if it's enabled and running in Kubernetes,
		// it may get killed by an Out-of-Memory (OOM) error due to a conflict with the Horizontal Pod Autoscaler (HPA).
		Prefork: false,
		// Which is suitable for streaming AI Response.
		StreamRequestBody:       true,
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
		// By default, it is set to 0.0.0.0/0 for local development; however, it can be bound to an ingress controller/proxy.
		// This can be a private IP range (e.g., 10.0.0.0/8).
		TrustedProxies: []string{"0.0.0.0/0"},
		// Trust X-Forwarded-For headers; additionally, this can be customized if using an ingress controller/proxy, especially Ingress Nginx.
		ProxyHeader: fiber.HeaderXForwardedFor, // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
	})
}

// startServer configures and starts the Fiber web server.
// It initializes logging, determines the server address, and calls the server start function
// with graceful shutdown handling.
//
// Note: Now that it supports HTTPS/TLS, it can be easily integrated with Kubernetes.
// For guidance on setting up HTTPS/TLS on Kubernetes, refer to:
//
// - https://kubernetes.io/docs/reference/kubectl/generated/kubectl_create/kubectl_create_secret_tls/
//
// For example, the following certificate can be used:
//
// - https://crt.sh/?q=d5b8a29e3eaf7413ee925dbb2ee9c9f9b6a73880fe0444704baaf71c1aa7feb3
//
// Note: The example certificate uses ECC (Elliptic Curve Cryptography), which is stable for internal mode and
// multiple clusters with many pods. It also helps alleviate the struggles that NGINX Ingress faces when handling
// a large number of concurrent requests, and it provides efficient bandwidth usage that saves cost.
//
// Also, note that startServer facilitates easy integration with HTTPS/TLS and supports ACME via cert-manager.io for Kubernetes.
// When running outside of Kubernetes (e.g., without an ingress), the PORT must be explicitly set to 443 for access via browser or other clients, as the default port is 8080.
// Make sure the certificate is correctly configured as well (e.g., the certificate chain, which is easy to handle in Go for chaining certificates).
// If the certificate is valid and properly configured, the server will run; otherwise, it won't run.
func startServer(app *fiber.App, config Config) {
	// Initialize logging with the application name and time format
	log.InitializeLogger(config.AppName, config.TimeFormat)

	// Define the server address using the specified port
	addr := fmt.Sprintf(":%s", config.Port)

	// Create a new instance of the server
	server := handler.NewFiberServer(app, config.AppName, config.MonitorPath)

	// Load TLS or mTLS certificates and keys from environment variables or command-line arguments ?
	//
	// TODO: Implement ACME?
	tlsConfig, err := setupTLS.LoadConfig()
	if err != nil {
		log.LogFatal(err)
	}

	// Start the server with graceful shutdown and monitor
	if tlsConfig != nil {
		// Start the server with TLS or mTLS ?
		handler.StartServer(server, addr, config.MonitorPath, config.ShutdownTimeout, tlsConfig, nil)
	} else {
		// Start the server without TLS
		handler.StartServer(server, addr, config.MonitorPath, config.ShutdownTimeout, nil, nil)
	}
}

// parseDuration converts a string to a time.Duration, logging an error and defaulting if necessary
func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.LogError(fmt.Errorf("invalid duration format: %s, using default 5s", durationStr))
		return 5 * time.Second
	}
	return duration
}
