// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package tls provides functionality for configuring TLS settings in a Go application.
// It supports loading TLS certificates and keys from environment variables and optionally
// configuring mutual TLS (mTLS) by loading CA certificates for client verification.
//
// This package is designed to be used in conjunction with the h0llyw00dz-template environment
// configuration, which defines the necessary environment variables for TLS setup.
//
// Environment Variables:
//   - env.SERVERCERTTLS: The path to the server's TLS certificate file.
//   - env.SERVERKEYTLS: The path to the server's TLS private key file.
//   - env.SERVERCATLS: The path to the CA certificate file for mTLS.
//   - env.ENABLEMTLS: A flag to enable mutual TLS (mTLS) if set to "true".
//
// Usage:
//
// To configure TLS, call LoadConfig, which returns a configured [tls.Config].
// If mTLS is enabled, it will also configure client certificate verification.
//
// Example:
//
//	tlsConfig, err := tls.LoadConfig()
//	if err != nil {
//	    log.Fatalf("Failed to load TLS configuration: %v", err)
//	}
//	// Use tlsConfig in your server setup
//
// Error Handling:
//
// The package defines ErrorMTLS for handling cases where CA certificates cannot be appended
// to the certificate pool, indicating issues with CA certificate loading or format.
package tls
