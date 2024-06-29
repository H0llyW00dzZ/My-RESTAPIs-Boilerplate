// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server_test

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid/stream"
	"h0llyw00dz-template/backend/internal/server"
	"math/big"
	"testing"
	"time"
)

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
			CommonName: "localhost",
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
			CommonName: "localhost",
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

// TestExtractSCTsFromCertificate tests the ExtractSCTsFromCertificate method of the CTVerifier.
func TestExtractSCTsFromCertificate(t *testing.T) {
	// Create a test certificate with SCTs
	cert, _ := createTestCertificateWithSCTs(t)

	// Create a new CTVerifier
	ctVerifier := new(server.CTVerifier)

	// Test case 1: Extract SCTs from a valid certificate
	t.Run("ExtractSCTsFromValidCertificate", func(t *testing.T) {
		scts, err := ctVerifier.ExtractSCTsFromCertificate(cert)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(scts) == 0 {
			t.Error("Expected SCTs, but got none")
		}
	})
}

// TestVerifySCT tests the VerifySCT method of the CTVerifier.
func TestVerifySCT(t *testing.T) {
	// Create a test certificate with SCTs
	cert, sctResponse := createTestCertificateValidSCTs(t)

	// Create a new CTVerifier
	ctVerifier := new(server.CTVerifier)

	// Test case 1: Verify a valid SCT
	t.Run("VerifyValidSCT", func(t *testing.T) {
		err := ctVerifier.VerifySCT(sctResponse, cert)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
	// Test case 2: Verify certificate transparency for a certificate without SCTs
	t.Run("VerifyCertificateTransparencyForCertificateWithoutSCTs", func(t *testing.T) {
		// Create a test certificate without SCTs
		certWithoutSCTs, _ := createTestCertificateValidSCTs(t)

		err := ctVerifier.VerifyCertificateTransparency(certWithoutSCTs)
		if err == nil {
			t.Error("Expected an error, but got none")
		}
	})
}

// createTestCertificateValidSCTsForLTS creates a test certificate valid SCTs for testing purposes.
func createTestCertificateValidSCTsForLTS(t *testing.T) (*x509.Certificate, crypto.PrivateKey) {
	// Generate a random serial number
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("Failed to generate serial number: %v", err)
	}

	// Set the validity period for the certificate
	notBefore := time.Now()
	notAfter := notBefore.Add(AheadTime24Hours)

	// Create the certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Issuer: pkix.Name{
			CommonName: "Gopher",
		},
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
	}

	// Generate a new ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Create a self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	// Create SCT data with a valid timestamp
	timestamp := uint64(time.Now().Unix())
	logID := make([]byte, 32)
	_, err = rand.Read(logID)
	if err != nil {
		t.Fatalf("Failed to generate log ID: %v", err)
	}

	extensions := []byte{} // Empty extensions

	sctData := make([]byte, 0, 44+len(cert.Raw))
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
	sctData = append(sctData, cert.Raw...)

	// Sign the SCT data
	hash := crypto.SHA256
	hasher := hash.New()
	hasher.Write(sctData)
	hashed := hasher.Sum(nil)
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hashed)
	if err != nil {
		t.Fatalf("Failed to generate signature: %v", err)
	}

	// Construct the SCT extension value
	sctExtensionValue := make([]byte, 0, 44+len(signature))
	sctExtensionValue = append(sctExtensionValue, byte(server.CTVersion1))
	sctExtensionValue = append(sctExtensionValue, logID...)
	sctExtensionValue = append(sctExtensionValue, byte(timestamp>>56),
		byte(timestamp>>48),
		byte(timestamp>>40),
		byte(timestamp>>32),
		byte(timestamp>>24),
		byte(timestamp>>16),
		byte(timestamp>>8),
		byte(timestamp))
	sctExtensionValue = append(sctExtensionValue, byte(len(extensions)))
	sctExtensionValue = append(sctExtensionValue, extensions...)
	sctExtensionValue = append(sctExtensionValue, signature...)

	// Attach the SCT to the certificate
	template.ExtraExtensions = append(template.ExtraExtensions, pkix.Extension{
		Id:       server.OIDExtensionCTSCT,
		Critical: false,
		Value:    sctExtensionValue,
	})

	// Create a new certificate with the SCT extension
	certWithSCT, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate with SCT: %v", err)
	}

	// Parse the new certificate
	finalCert, err := x509.ParseCertificate(certWithSCT)
	if err != nil {
		t.Fatalf("Failed to parse certificate with SCT: %v", err)
	}

	return finalCert, privateKey
}

// TestVerifyCertificateTransparencyInTLSConnection tests the certificate transparency verification in a TLS connection.
//
// Note: This method currently works only with ECDSA certificates, not with other certificate.
func TestVerifyCertificateTransparencyInTLSConnection(t *testing.T) {
	// Create a test certificate with SCTs
	cert, privateKey := createTestCertificateValidSCTsForLTS(t)

	// Create a server TLS configuration with the test certificate
	serverTLSConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		Rand: server.RandTLS(),
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert.Raw},
				PrivateKey:  privateKey,
			},
		},
	}

	// Create a client TLS configuration with CT verification
	clientTLSConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// Create a new CTVerifier
			ctVerifier := new(server.CTVerifier)

			// Perform Certificate Transparency checks
			for _, chain := range verifiedChains {
				for _, cert := range chain {
					if err := ctVerifier.VerifyCertificateTransparency(cert); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	// Start a TLS server with the test certificate
	listener, err := tls.Listen("tcp", "localhost:443", serverTLSConfig)
	if err != nil {
		t.Fatalf("Failed to create TLS listener: %v", err)
	}
	defer listener.Close()

	// Start a goroutine to handle the TLS connection
	errChan := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- fmt.Errorf("failed to accept TLS connection: %v", err)
			return
		}
		defer conn.Close()

		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			errChan <- fmt.Errorf("expected a TLS connection")
			return
		}

		// Read the message from the client
		buffer := make([]byte, stream.ChunkSize)
		n, err := tlsConn.Read(buffer)
		if err != nil {
			errChan <- fmt.Errorf("failed to read from TLS connection: %v", err)
			return
		}
		message := string(buffer[:n])
		t.Logf("Server received: %s", message)

		// Send a response back to the client
		response := "Hello, client!"
		_, err = tlsConn.Write([]byte(response))
		if err != nil {
			errChan <- fmt.Errorf("failed to write to TLS connection: %v", err)
			return
		}
		t.Logf("Server sent: %s", response)

		errChan <- nil
	}()

	// Test case 1: Connect to the TLS server with a valid certificate
	t.Run("ConnectWithValidCertificate", func(t *testing.T) {
		conn, err := tls.Dial("tcp", listener.Addr().String(), clientTLSConfig)
		if err != nil {
			t.Errorf("Failed to establish TLS connection: %v", err)
			return
		}
		defer conn.Close()

		// Send a message to the server
		message := "Hello, server!"
		_, err = conn.Write([]byte(message))
		if err != nil {
			t.Errorf("Failed to write to TLS connection: %v", err)
			return
		}
		t.Logf("Client sent: %s", message)

		// Read the response from the server
		buffer := make([]byte, stream.ChunkSize)
		n, err := conn.Read(buffer)
		if err != nil {
			t.Errorf("Failed to read from TLS connection: %v", err)
			return
		}
		response := string(buffer[:n])
		t.Logf("Client received: %s", response)
	})

	// Wait for the server goroutine to finish
	if err := <-errChan; err != nil {
		t.Errorf("Server error: %v", err)
	}
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

// TestInvalidSCT tests the Invalid method of the CTVerifier.
func TestInvalidSCT(t *testing.T) {
	// Generate a dummy SCT data
	timestamp := make([]byte, 8)
	binary.BigEndian.PutUint64(timestamp, uint64(time.Now().Unix()))

	// Create a new CTVerifier
	ctVerifier := new(server.CTVerifier)

	// Generate a random serial number
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Set the certificate subject and issuer
	subject := pkix.Name{
		Organization: []string{"Example Inc."},
		CommonName:   "example.com",
	}

	// Set the certificate validity period
	notBefore := time.Now()
	notAfter := notBefore.Add(24 * time.Hour)

	// Test case 1: Invalid SCT data (Insufficient Length)
	t.Run("InvalidSCTDataInsufficientLength", func(t *testing.T) {
		// Add an SCT extension to the certificate template
		sctExtension, err := generateSCTExtension(
			server.SCTData{
				Version:    1,
				LogID:      []byte("invalid"),
				Timestamp:  asn1.RawValue{FullBytes: timestamp},
				Extensions: []byte("invalid"),
				Signature:  []byte("invalid"),
			},
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Generate a new RSA private key
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Generate a self-signed certificate with the SCT extension
		template := x509.Certificate{
			SerialNumber:    serialNumber,
			Subject:         subject,
			Issuer:          subject,
			NotBefore:       notBefore,
			NotAfter:        notAfter,
			KeyUsage:        x509.KeyUsageDigitalSignature,
			ExtKeyUsage:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			ExtraExtensions: []pkix.Extension{sctExtension},
		}

		// Create a self-signed certificate
		certx, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Parse the certificate
		x509Cert, err := x509.ParseCertificate(certx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		_, err = ctVerifier.ExtractSCTsFromCertificate(x509Cert)
		if err == nil {
			t.Errorf("Expected error due to invalid, but got nil.")
		} else if err.Error() != "invalid SCT data: insufficient length" {
			t.Errorf("Expected invalid SCT data: insufficient length, but got: %v", err)
		} else {
			t.Logf("Extract SCT Data failed as expected: %v", err)
		}
	})

	// Test case 2: Invalid SCT data (insufficient length for signature)
	t.Run("InvalidSCTDataInsufficientLengthSignature", func(t *testing.T) {
		// Add an SCT extension to the certificate template
		sctExtension, err := generateSCTExtension(
			server.SCTData{
				Version:    1,
				LogID:      []byte("insufficient_signature"),
				Timestamp:  asn1.RawValue{FullBytes: timestamp},
				Extensions: []byte("insufficient_signature"),
				Signature:  []byte("insufficient_signature"),
			},
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Generate a new RSA private key
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Generate a self-signed certificate with the SCT extension
		template := x509.Certificate{
			SerialNumber:    serialNumber,
			Subject:         subject,
			Issuer:          subject,
			NotBefore:       notBefore,
			NotAfter:        notAfter,
			KeyUsage:        x509.KeyUsageDigitalSignature,
			ExtKeyUsage:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			ExtraExtensions: []pkix.Extension{sctExtension},
		}

		// Create a self-signed certificate
		certx, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Parse the certificate
		x509Cert, err := x509.ParseCertificate(certx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		_, err = ctVerifier.ExtractSCTsFromCertificate(x509Cert)
		if err == nil {
			t.Errorf("Expected error due to invalid, but got nil.")
		} else if err.Error() != "invalid SCT data: insufficient length for signature" {
			t.Errorf("Expected invalid SCT data: insufficient length for signature, but got: %v", err)
		} else {
			t.Logf("Extract SCT Data failed as expected: %v", err)
		}
	})
}
