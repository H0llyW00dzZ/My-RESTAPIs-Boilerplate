// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"

	log "h0llyw00dz-template/backend/internal/logger"
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
	appName = env.GetEnv(env.APPNAME, "Gopher")
	port = env.GetEnv(env.PORT, "8080")
	monitorPath = env.GetEnv(env.MONITORPATH, "/monitor")
	// Get the TIME_FORMAT from environment variables or use default value
	// Note: List Time Format Available: unix,default
	timeFormat = env.GetEnv(env.TIMEFORMAT, "unix")

	// Get the READ_TIMEOUT, WRITE_TIMEOUT, and SHUTDOWN_TIMEOUT from environment variables or use default values.
	// Note: These default timeout values (5 seconds) are set to help prevent potential deadlocks/hangs.
	readTimeoutStr := env.GetEnv(env.READTIMEOUT, "5s")
	writeTimeoutStr := env.GetEnv(env.SHUTDOWNTIMEOUT, "5s")
	shutdownTimeoutStr := env.GetEnv(env.SHUTDOWNTIMEOUT, "5s")

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
		// When running behind an ingress controller/proxy, disable "EnableIPValidation"
		// because the ingress controller/proxy will forward the real IP anyway from the header, which is already valid.
		EnableIPValidation:      false,
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
// Also, note that startServer facilitates easy integration with HTTPS/TLS and supports ACME via cert-manager.io for Kubernetes.
func startServer(app *fiber.App, appName, port, monitorPath, timeFormat string, shutdownTimeout time.Duration) {
	// Initialize the logger with the AppName from the environment variable
	log.InitializeLogger(app.Config().AppName, timeFormat)

	// Define server address
	addr := fmt.Sprintf(":%s", port) // Use the port from the environment variable

	// Create a new instance of FiberServer
	server := handler.NewFiberServer(app, appName, monitorPath)

	// Load TLS certificate and key from environment variables or command-line arguments
	//
	// TODO: ACME Implementations ?
	tlsCertFile := env.GetEnv(env.SERVERCERTTLS, "")
	tlsKeyFile := env.GetEnv(env.SERVERKEYTLS, "")

	var tlsConfig *tls.Config
	if tlsCertFile != "" && tlsKeyFile != "" {
		// Note: Fiber uses ECC is significantly faster compared to Nginx uses ECC, which struggles to handle a billion concurrent requests.
		cert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
		if err != nil {
			log.LogError(err)
			os.Exit(1)
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	// Start the server with graceful shutdown and monitor
	if tlsConfig != nil {
		// Start the server with TLS
		handler.StartServer(server, addr, monitorPath, shutdownTimeout, tlsConfig, nil)
	} else {
		// Start the server without TLS
		handler.StartServer(server, addr, monitorPath, shutdownTimeout, nil, nil)
	}
}
