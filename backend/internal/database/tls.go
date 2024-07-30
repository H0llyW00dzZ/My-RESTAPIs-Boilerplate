// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package database

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// loadMySQLRootCA loads the MySQL root CA certificate from the environment variable.
//
// Note: It is now extracted into two functions for CA certificates because other databases (e.g., Redis) might have a different issuer even if the root CA is the same.
// The reason for extracting it into two functions is that it makes it easier to set up in a load-balancing cloud environment.
// It is also recommended to use CA chains (e.g, Root CA + Subs CA Without Leaf CA) instead of only the root CA.
//
// Also note that, regarding load-balancing or none, since this repository is written in the Fiber framework with some custom advances,
// the speed of the site depends on the database connection. For example, if both MySQL and Redis are in close proximity to this repository (let's say pods),
// the site can load faster when a visitor/client visits the website.
func loadMySQLRootCA() (*x509.CertPool, error) {
	rootCABase64 := mysqltlsCAs
	if rootCABase64 == "" {
		return nil, fmt.Errorf("MYSQL_CERTS_TLS environment variable is not set")
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

// loadRedisRootCA loads the Redis root CA certificate from the environment variable.
//
// Note: It is now extracted into two functions for CA certificates because other databases (e.g., MySQL) might have a different issuer even if the root CA is the same.
// The reason for extracting it into two functions is that it makes it easier to set up in a load-balancing cloud environment.
// It is also recommended to use CA chains (e.g, Root CA + Subs CA Without Leaf CA) instead of only the root CA.
//
// Also note that, regarding load-balancing or none, since this repository is written in the Fiber framework with some custom advances,
// the speed of the site depends on the database connection. For example, if both MySQL and Redis are in close proximity to this repository (let's say pods),
// the site can load faster when a visitor/client visits the website.
func loadRedisRootCA() (*x509.CertPool, error) {
	rootCABase64 := redistlsCAs
	if rootCABase64 == "" {
		return nil, fmt.Errorf("REDIS_CERTS_TLS environment variable is not set")
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
