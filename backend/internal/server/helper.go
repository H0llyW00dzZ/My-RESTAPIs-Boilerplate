// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"h0llyw00dz-template/backend/internal/database"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// InitializeTables creates all the necessary tables in the database.
//
// Example usage:
//
//	// createUserTable creates the User table in the database if it doesn't exist.
//	func createUserTable(db database.Service) error {
//		// Define the SQL query to create the User table.
//		query := `
//		CREATE TABLE IF NOT EXISTS User (
//			id INT AUTO_INCREMENT PRIMARY KEY,
//			username VARCHAR(255) NOT NULL,
//			email VARCHAR(255) NOT NULL,
//			password VARCHAR(255) NOT NULL,
//			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
//		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
//		`
//
//		// Execute the SQL query without expecting any rows to be returned.
//		err := db.ExecWithoutRow(context.Background(), query)
//		if err != nil {
//			return err
//		}
//
//		// Log a success message indicating that the User table was initialized.
//		log.LogInfo("Successfully initialized the User table.")
//		return nil
//	}
//
//	// InitializeTables creates all the necessary tables in the database.
//	func InitializeTables(db database.Service) error {
//		// Call the createTables function with the database service and table creation functions.
//		return createTables(db,
//			createTable("User", createUserTable),
//			// Add more table creation functions as needed
//		)
//	}
func InitializeTables(db database.Service) error {
	// Note: This approach provides a more flexible and scalable way to initialize database tables compared to using an ORM system.
	// It allows for easy initialization or migration of tables, and can handle a large number of database schemas (e.g, 1 billion database schemas 🔥) without limitations.
	return createTables(db)

}

// createTable is a higher-order function that creates a table in the database.
// It takes the table name and a function that defines the table creation logic.
// It returns a function that accepts a database.Service and executes the table creation logic.
// If an error occurs during table creation, it wraps the error with a descriptive message.
func createTable(name string, fn func(database.Service) error) func(database.Service) error {
	return func(db database.Service) error {
		if err := fn(db); err != nil {
			return fmt.Errorf("failed to create %s table: %v", name, err)
		}
		return nil
	}
}

// createTables creates all the necessary tables in the database.
// It accepts a database.Service and a variadic list of functions that define the table creation logic.
// It iterates over the table creation functions and executes each one.
// If an error occurs during table creation, it returns the error.
func createTables(db database.Service, tables ...func(database.Service) error) error {
	for _, table := range tables {
		if err := table(db); err != nil {
			return err
		}
	}
	return nil
}

// isBrowserRequest checks if the given data represents a browser request.
func isBrowserRequest(data []byte) bool {
	// Check if the data starts with a valid HTTP method
	// Note: This is a raw packet, so it's different because it includes a space after the method.
	methods := [][]byte{
		[]byte(fiber.MethodGet + " "),
		[]byte(fiber.MethodHead + " "),
		[]byte(fiber.MethodPost + " "),
		[]byte(fiber.MethodPut + " "),
		[]byte(fiber.MethodDelete + " "),
		[]byte(fiber.MethodConnect + " "),
		[]byte(fiber.MethodOptions + " "),
		[]byte(fiber.MethodTrace + " "),
		[]byte(fiber.MethodPatch + " "),
	}

	for _, method := range methods {
		if bytes.HasPrefix(data, method) {
			return true
		}
	}

	return false
}

// MakeHTTPRequest is a helper function that makes an HTTP request using TLS 1.3.
//
// Note: This uses the standard library because it is only used for activation and certification, similar to them.
func (s *FiberServer) MakeHTTPRequest(req *http.Request) (*http.Response, error) {
	// Create a custom TLS configuration with TLS 1.3 enabled
	//
	// Note: The cipher/preferred cipher in Go's standard TLS 1.3 implementation does not allow direct configuration of cipher suites. See the note about TLSConfig in "run.go".
	// That's why it is kept like this, as it doesn't work when set to ChaCha.
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
	}

	// Create an HTTP client with the custom TLS configuration
	//
	// TODO: Implement a helper function for handling httpcookiejar/cookiejar. This might be needed for transport if the HTTP server response sends cookies.
	// Then Store the cookies in a file with a JSON format, and encrypt the JSON values. However, it requires storage of around 500MB ~ 1GB or lowerthan 500mb for the container.
	//
	// Also, note that when building an image (from this repo) for a container (e.g., for K8s), make sure to set the container to a non-root user (never set it to root).
	// Because This repository does not rely on the operating system or any other specific dependencies, only the resources of the machine (e.g., memory, cpu, disk) and networking.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			// Note: This should be enough for HTTP Client in the backend server.
			MaxIdleConns:          25,
			IdleConnTimeout:       5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ForceAttemptHTTP2:     true,
		},
		Timeout: 10 * time.Second,
	}

	// Send the HTTP request using the client
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// HTTPRequestMaker is a type that wraps the MakeHTTPRequest method.
type HTTPRequestMaker struct {
	MakeHTTPRequestFunc func(req *http.Request) (*http.Response, error)
}

// MakeHTTPRequest calls the wrapped MakeHTTPRequestFunc.
func (h *HTTPRequestMaker) MakeHTTPRequest(req *http.Request) (*http.Response, error) {
	return h.MakeHTTPRequestFunc(req)
}
