// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

// Note: This "helper_tls_test.go" file contains helper functions used for testing purpose.
// Reason for extraction: Making testing easier by ensuring TLS tests are as production-like as possible (e.g., TLS Config must set InsecureSkipVerify to false).

package server_test

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/server"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Note: This environment is for testing TLS 1.3. It's crucial to make testing TLS 1.3 as production-like as possible.
// Setting InsecureSkipVerify to true would significantly hinder testing, as it wouldn't mimic real-world TLS behavior.
// Also note that this must be a valid domain name that is bound to the host. While domain names are relatively inexpensive to acquire,
// it's essential to use a valid one for accurate TLS 1.3 testing.
//
// Demo: api-beta.btz.pm
// Server Backend: Heroku (Due it's free and perfect for demo/test about TLS)
// Server Frontend: Cloudflare (Paid $10 to get ACM)
// Scan Result (This site are accurate): https://decoder.link/sslchecker/api-beta.btz.pm/443
var testHostName = os.Getenv("TEST_HOSTNAME") // Use Real domain (e.g, testing-tls.go.dev)

func copySysCertPoolFromFile(certFilePath string) (*x509.CertPool, error) {
	// Read the CA certificate from the file
	caCert, err := os.ReadFile(certFilePath)
	if err != nil {
		return nil, err
	}

	// Create a new certificate pool
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	// Append the CA certificate to the pool
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("error appending CA certificate to pool")
	}

	return certPool, nil
}

func createCertPoolFromFile(certFilePath string) (*x509.CertPool, error) {
	// Read the CA certificate from the file
	caCert, err := os.ReadFile(certFilePath)
	if err != nil {
		return nil, err
	}

	// Create a new certificate pool
	certPool := x509.NewCertPool()

	// Append the CA certificate to the pool
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("error appending CA certificate to pool")
	}

	return certPool, nil
}

func tlsServerConfig(cert tls.Certificate) *tls.Config {
	log.InitializeLogger("Boring TLS 1.3 Testing", "")
	tlsHandler := &fiber.TLSHandler{}
	RootCA, _ := createCertPoolFromFile("boring-RootCA.pem")
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			// Note: These are classical elliptic curves for TLS 1.3 key exchange.
			// For experimental purposes related to post-quantum hybrid design, refer to:
			// https://datatracker.ietf.org/doc/html/draft-ietf-tls-hybrid-design-10
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		Certificates:   []tls.Certificate{cert},
		RootCAs:        RootCA,
		GetCertificate: tlsHandler.GetClientInfo,
		// Note: This doesn't need to be explicitly set to "tls.RequireAndVerifyClientCert" because the Go TLS standard library
		// defaults to verifying client certificates when ClientCAs is set.
		// Also, note that ClientCAs refers to the chain of Certificate Authorities Pool that made & signed by RootCAs, which is why it's different from RootCAs.
		ClientAuth: tls.RequireAndVerifyClientCert,
		Rand:       server.RandTLS(),
	}
}

func clientTLSConfig() *tls.Config {
	log.InitializeLogger("Boring TLS 1.3 Testing", "")
	certPool, _ := createCertPoolFromFile("boring-ca.pem")
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			// Note: These are classical elliptic curves for TLS 1.3 key exchange.
			// For experimental purposes related to post-quantum hybrid design, refer to:
			// https://datatracker.ietf.org/doc/html/draft-ietf-tls-hybrid-design-10
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		ClientCAs:  certPool,
		ServerName: testHostName,
	}
}

const (
	AheadTime24Hours = 24 * time.Hour
	AheadTime7Days   = 7 * 24 * time.Hour
	AheadTime30Days  = 30 * 24 * time.Hour
	Expired          = -time.Hour * 24
)

// MockHTTPClient is a mock implementation of the HTTP client.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do mocks the Do function of the HTTP client.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

func generateSelfSignedCertECDSA() (*x509.Certificate, *ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: testHostName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(AheadTime24Hours),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, privateKey, nil
}

func generateSelfSignedCertRSA() (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: testHostName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(AheadTime7Days),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, privateKey, nil
}

func generateSelfSignedCertEd25519() (*x509.Certificate, ed25519.PrivateKey, error) {
	privateKey := ed25519.NewKeyFromSeed(make([]byte, ed25519.SeedSize))
	publicKey := privateKey.Public().(ed25519.PublicKey)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: testHostName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(AheadTime30Days),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, privateKey, nil
}

func generateSelfSignedCertEd25519WithExpired() (*x509.Certificate, ed25519.PrivateKey, error) {
	privateKey := ed25519.NewKeyFromSeed(make([]byte, ed25519.SeedSize))
	publicKey := privateKey.Public().(ed25519.PublicKey)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: testHostName,
		},
		NotBefore: time.Now().Add(Expired),      // Set NotBefore to 24 hours ago
		NotAfter:  time.Now().Add(-time.Minute), // Set NotAfter to 1 minute ago

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, privateKey, nil
}

// createTestCertificateWithSCTs creates a test certificate with SCTs for testing purposes.
func createTestCertificateWithSCTs(t *testing.T) (*x509.Certificate, *server.SCTResponse) {
	// Generate an ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA private key: %v", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("Failed to generate serial number: %v", err)
	}

	// Create SCT data with a valid timestamp
	timestamp := uint64(time.Now().Unix())
	logID := make([]byte, 32)
	_, err = rand.Read(logID)
	if err != nil {
		t.Fatalf("Failed to generate log ID: %v", err)
	}

	extensions := []byte{0x00} // Empty extensions

	sctData := make([]byte, 0, 44)
	sctData = append(sctData, byte(server.CTVersion1))
	sctData = append(sctData, logID...)
	sctData = append(sctData, byte(timestamp>>56),
		byte(timestamp>>48),
		byte(timestamp>>40),
		byte(timestamp>>32),
		byte(timestamp>>24),
		byte(timestamp>>16),
		byte(timestamp>>8),
		byte(timestamp))
	sctData = append(sctData, byte(len(extensions)))
	sctData = append(sctData, extensions...)

	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, sctData)
	if err != nil {
		t.Fatalf("Failed to generate signature: %v", err)
	}

	sctData = append(sctData, signature...)

	sctResponse := &server.SCTResponse{
		SCTVersion: server.CTVersion1,
		ID:         base64.StdEncoding.EncodeToString(logID),
		Timestamp:  timestamp,
		Extensions: base64.StdEncoding.EncodeToString(extensions),
		Signature:  base64.StdEncoding.EncodeToString(signature),
	}

	// Create a self-signed certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Issuer: pkix.Name{
			CommonName: "Gopher",
		},
		Subject: pkix.Name{
			CommonName: testHostName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(AheadTime24Hours),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		ExtraExtensions: []pkix.Extension{
			{
				Id:       server.OIDExtensionCTSCT,
				Critical: false,
				Value:    sctData,
			},
		},
	}

	// Create the self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create self-signed certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	return cert, sctResponse
}

// createTestCertificateValidSCTs creates a test certificate valid SCTs for testing purposes.
func createTestCertificateValidSCTs(t *testing.T) (*x509.Certificate, *server.SCTResponse) {
	// Generate an ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA private key: %v", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("Failed to generate serial number: %v", err)
	}

	// Create a self-signed certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Issuer: pkix.Name{
			CommonName: "Gopher",
		},
		Subject: pkix.Name{
			CommonName: testHostName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(AheadTime24Hours),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
	}

	// Create the self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create self-signed certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	// Create SCT data with a valid timestamp
	timestamp := uint64(time.Now().Unix())
	var data []byte
	data = append(data, cert.Raw...)
	data = append(data, byte(server.CTVersion1))
	data = append(data, []byte("test-ct-log")...)
	data = append(data, byte(timestamp>>56),
		byte(timestamp>>48),
		byte(timestamp>>40),
		byte(timestamp>>32),
		byte(timestamp>>24),
		byte(timestamp>>16),
		byte(timestamp>>8),
		byte(timestamp))
	data = append(data, []byte("")...)

	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, data)
	if err != nil {
		t.Fatalf("Failed to generate signature: %v", err)
	}

	sctResponse := &server.SCTResponse{
		SCTVersion: server.CTVersion1,
		ID:         "test-ct-log",
		Timestamp:  timestamp,
		Extensions: "",
		Signature:  base64.StdEncoding.EncodeToString(signature),
	}

	return cert, sctResponse
}

func generateSCTExtension(SCData server.SCTData) (pkix.Extension, error) {
	// Marshal the SCT data into ASN.1 format
	sctBytes, err := asn1.Marshal(SCData)
	if err != nil {
		return pkix.Extension{}, err
	}

	// Create the SCT extension
	sctExtension := pkix.Extension{
		Id:       server.OIDExtensionCTSCT,
		Critical: false,
		Value:    sctBytes,
	}

	return sctExtension, nil
}
