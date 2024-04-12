// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package main

import (
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
	appName, port, monitorPath := getEnvVariables()
	app := setupFiber(appName)

	// Start the server with graceful shutdown and monitor
	startServer(app, appName, port, monitorPath)
}

// getEnvVariables retrieves essential configuration settings from environment variables.
// It provides default values for the application name, port, and monitoring path
// to ensure the application has sensible defaults if environment variables are not set.
func getEnvVariables() (appName, port, monitorPath string) {
	// Get the APP_NAME, PORT, and MONITOR_PATH from environment variables or use default values.
	appName = getEnv("APP_NAME", "Gopher")
	port = getEnv("PORT", "8080")
	monitorPath = getEnv("MONITOR_PATH", "/monitor")
	return
}

// setupFiber initializes a new Fiber application with custom configuration.
// It sets up the JSON encoder/decoder, case sensitivity, and strict routing,
// and applies the application name to the server headers.
func setupFiber(appName string) *fiber.App {
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
		ReadTimeout:      5 * time.Second,
	})
}

// startServer configures and starts the Fiber web server.
// It initializes logging, determines the server address, sets up a shutdown timeout,
// and calls the server start function with graceful shutdown handling.
func startServer(app *fiber.App, appName, port, monitorPath string) {
	// Initialize the logger with the AppName from the environment variable
	log.InitializeLogger(app.Config().AppName)

	// Define server address, shutdown timeout, and monitor path
	addr := fmt.Sprintf(":%s", port) // Use the port from the environment variable
	shutdownTimeout := 5 * time.Second

	// Create a new instance of FiberServer
	server := handler.NewFiberServer(app, appName, monitorPath)

	// Start the server with graceful shutdown and monitor
	handler.StartServer(server, addr, monitorPath, shutdownTimeout)
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
