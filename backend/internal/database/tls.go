// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package database

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// loadRootCA loads the root CA certificate from the environment variable.
func loadRootCA() (*x509.CertPool, error) {
	rootCABase64 := tlsCAs
	if rootCABase64 == "" {
		return nil, fmt.Errorf("EXTRA_CERTS_TLS environment variable is not set")
	}

	rootCABytes, err := base64.StdEncoding.DecodeString(rootCABase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode root CA: %v", err)
	}

	rootCAs := x509.NewCertPool()
	if ok := rootCAs.AppendCertsFromPEM(rootCABytes); !ok {
		return nil, fmt.Errorf("failed to append root CA to cert pool")
	}

	return rootCAs, nil
}
