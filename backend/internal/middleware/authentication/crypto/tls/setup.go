// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"h0llyw00dz-template/env"
	"os"
)

var (
	// ErrorMTLS is returned when the CA certificates cannot be appended to the certificate pool.
	ErrorMTLS = errors.New("crypto/mtls: Failed to append CA certificates")

	// ErrorMTLSCAEmpty is returned when the CA certificate file is not provided.
	ErrorMTLSCAEmpty = errors.New("crypto/mtls: CA certificate file not provided")
)

var (
	// tlsCertFile holds the path to the server TLS certificate file.
	tlsCertFile = env.GetEnv(env.SERVERCERTTLS, "")

	// tlsKeyFile holds the path to the server TLS key file.
	tlsKeyFile = env.GetEnv(env.SERVERKEYTLS, "")

	// caCertFile holds the path to the CA certificate file.
	caCertFile = env.GetEnv(env.SERVERCATLS, "")
)

var (
	// enableMTLS indicates whether mutual TLS is enabled, based on the environment variable.
	enableMTLS = env.GetEnv(env.ENABLEMTLS, "") == "true"
)

// LoadConfig loads TLS configuration based on environment variables.
func LoadConfig() (*tls.Config, error) {
	if tlsCertFile != "" && tlsKeyFile != "" {
		// Note: Fiber uses ECC is significantly faster compared to Nginx uses ECC, which struggles to handle a billion concurrent requests.
		cert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load key pair: %w", err)
		}

		// Note: For ECC the OCSP, it's optional if explicitly set to TLSv1.3 and used in internal mode.
		// However, if it's used externally and allows TLSv1.2, then OCSP should be configured, provided that
		// there is knowledge on how to set it up.
		//
		// For an example of OCSP stapling and TLSv1.2 configuration (using "nginx.ingress.kubernetes.io/backend-protocol: HTTPS", enable-ocsp) that follows best practices for securing websites, see:
		//
		// - https://www.immuniweb.com/ssl/git.b0zal.io/KRIX2G2F/ (most all green)
		// - https://www.immuniweb.com/ssl/api.b0zal.io/VPdKSN3p/ (most all green)
		// - https://decoder.link/sslchecker/git.b0zal.io/443
		// - https://decoder.link/sslchecker/b0zal.io/443 (from this repository boilerplate is used for sandbox development exposed to public/prods)
		// - https://decoder.link/sslchecker/api.b0zal.io/443 (from this repository boilerplate is used for sandbox development exposed to public/prods)
		//
		// Additionally, note that if "enable-ocsp" is set to true in the Ingress Nginx ConfigMap, OCSP Stapling remains optional.
		// This is because when Nginx passes requests to HTTPS/TLS related to this service without terminating it,
		// as long as the certificate is the same for both Ingress and this service, OCSP Stapling can still be utilized.
		// If your cluster has any Kubernetes network mechanism that doesn't work with these configurations (e.g., nginx.ingress.kubernetes.io/backend-protocol: HTTPS, enable-ocsp),
		// then there may be an issue with your Kubernetes network configuration.
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		// This boolean is determined by the environment variable ENABLE_MTLS using env.GetEnv, which performs a lookup.
		// Unlike other environment variables (e.g., tlsCertFile), it does not directly use the value ENABLE_MTLS=true.
		if enableMTLS {
			caCertPool, err := loadCA()
			if err != nil {
				return nil, err
			}
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = caCertPool
		} else {
			tlsConfig.ClientAuth = tls.NoClientCert
		}

		return tlsConfig, nil
	}

	return nil, nil
}

// loadCA loads the CA certificates for client verification.
func loadCA() (*x509.CertPool, error) {
	if caCertFile == "" {
		return nil, ErrorMTLSCAEmpty
	}
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("error loading CA certificates: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, ErrorMTLS
	}

	return caCertPool, nil
}
