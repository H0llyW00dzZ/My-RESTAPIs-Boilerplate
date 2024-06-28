// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server_test

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"

	"h0llyw00dz-template/backend/internal/server"
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
			CommonName: "example.com",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

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
			CommonName: "example.com",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

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
			CommonName: "example.com",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

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

func TestSubmitToCTLog(t *testing.T) {
	// Generate a self-signed certificate with a valid private key
	cert, privateKey, err := generateSelfSignedCertECDSA()
	if err != nil {
		t.Fatal(err)
	}

	// Create a test CT log
	ctLog := server.CTLog{
		URL: "https://ct.example.com",
	}

	// Create a test Fiber server
	app := fiber.New(fiber.Config{
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	fiberServer := &server.FiberServer{
		App: app,
	}

	// Create an HTTPRequestMaker with the original MakeHTTPRequest method
	httpRequestMaker := &server.HTTPRequestMaker{
		MakeHTTPRequestFunc: fiberServer.MakeHTTPRequest,
	}

	// Test case 1: Successful submission to CT log with ECDSA key
	t.Run("SuccessfulSubmissionECDSACTVersion1", func(t *testing.T) {
		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Calculate the SHA-256 hash of the certificate
				hash := sha256.Sum256(cert.Raw)

				// Prepare the mock response with a valid SCT response
				timestamp := uint64(time.Now().Unix())
				data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

				signature, err := ecdsa.SignASN1(rand.Reader, privateKey, data)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}

				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and private key
		err := fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verification & Certificate Transparency submitted successfully
		t.Log("Hello Crypto: Certificate submitted to CT log successfully")
	})

	// Test case 2: Failed submission to CT log
	t.Run("FailedSubmission", func(t *testing.T) {
		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock error response
				mockErrorResponse := &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       http.NoBody,
				}
				return mockErrorResponse, errors.New("submission failed")
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker
		err := fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if err != nil && err.Error() != "failed to submit certificate to CT log: submission failed" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 3: Invalid SCT response
	t.Run("InvalidSCTResponse", func(t *testing.T) {
		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with an invalid SCT response
				invalidSCTResponse := server.SCTResponse{
					SCTVersion:   server.CTVersion2 + 1, // Unsupported SCT version
					ID:           "test-ct-log",
					Timestamp:    uint64(time.Now().Unix()),
					Extensions:   "",
					STHExtension: "",
					Signature:    base64.StdEncoding.EncodeToString([]byte("invalid-signature")),
				}
				responseBody, _ := sonic.Marshal(invalidSCTResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker
		err := fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := fmt.Sprintf("unsupported SCT version: %d", server.CTVersion2+1)
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 4: Successful submission to CT log with RSA key
	t.Run("SuccessfulSubmissionRSACTVersion1", func(t *testing.T) {
		// Generate a self-signed certificate with a valid RSA private key
		cert, privateKey, err := generateSelfSignedCertRSA()
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Calculate the SHA-256 hash of the certificate
				hash := sha256.Sum256(cert.Raw)

				// Prepare the mock response with a valid SCT response
				timestamp := uint64(time.Now().Unix())
				data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

				// Hash the data before signing
				hasher := sha256.New()
				hasher.Write(data)
				hashedData := hasher.Sum(nil)

				signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashedData)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}

				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and RSA private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verification & Certificate Transparency submitted successfully
		t.Log("Hello Crypto: Certificate submitted to CT log successfully")
	})

	// Test case 5: Successful submission to CT log using MakeHTTPRequest directly
	t.Run("SuccessfulSubmissionMakeHTTPRequest", func(t *testing.T) {
		// Generate a self-signed certificate with a valid ECDSA private key
		cert, privateKey, err := generateSelfSignedCertECDSA()
		if err != nil {
			t.Fatal(err)
		}

		// Create a test server that mocks the CT log server
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Calculate the SHA-256 hash of the certificate
			hash := sha256.Sum256(cert.Raw)

			// Prepare the mock response with a valid SCT response
			timestamp := uint64(time.Now().Unix())
			data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

			signature, err := ecdsa.SignASN1(rand.Reader, privateKey, data)
			if err != nil {
				t.Fatalf("Failed to generate signature: %v", err)
			}

			sctResponse := server.SCTResponse{
				SCTVersion: server.CTVersion1,
				ID:         "test-ct-log",
				Timestamp:  timestamp,
				Extensions: "",
				Signature:  base64.StdEncoding.EncodeToString(signature),
			}
			responseBody, _ := sonic.Marshal(sctResponse)
			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
		}))
		defer testServer.Close()

		// Create a test CT log with the test server URL
		ctLog := server.CTLog{
			URL: testServer.URL,
		}

		// Create a test Fiber server
		app := fiber.New(fiber.Config{
			JSONEncoder:  sonic.Marshal,
			JSONDecoder:  sonic.Unmarshal,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
		fiberServer := &server.FiberServer{
			App: app,
		}

		// Call the SubmitToCTLog method directly using MakeHTTPRequest
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, nil)

		// Assert the expectations
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verification & Certificate Transparency submitted successfully
		t.Log("Hello Crypto: Certificate submitted to CT log successfully")
	})

	// Test case 6: Successful submission to CT log with CTVersion2
	t.Run("SuccessfulSubmissionECDSACTVersion2", func(t *testing.T) {
		// Generate a self-signed certificate with a valid ECDSA private key
		cert, privateKey, err := generateSelfSignedCertECDSA()
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response for CTVersion2
				timestamp := uint64(time.Now().Unix())
				transItem := struct {
					SCTVersion   uint8
					Timestamp    uint64
					Extensions   []byte
					STHExtension []byte
				}{
					SCTVersion:   server.CTVersion2,
					Timestamp:    timestamp,
					Extensions:   []byte(""),
					STHExtension: []byte(""),
				}
				transItemBytes, _ := sonic.Marshal(transItem)

				signature, err := ecdsa.SignASN1(rand.Reader, privateKey, transItemBytes)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}

				sctResponse := server.SCTResponse{
					SCTVersion:   server.CTVersion2,
					ID:           "test-ct-log",
					Timestamp:    timestamp,
					Extensions:   "",
					STHExtension: "",
					Signature:    base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verification & Certificate Transparency submitted successfully
		t.Log("Hello Crypto: Certificate submitted to CT log successfully")
	})

	// Test case 7: Successful submission to CT log with RSA key and CTVersion2
	t.Run("SuccessfulSubmissionRSACTVersion2", func(t *testing.T) {
		// Generate a self-signed certificate with a valid RSA private key
		cert, privateKey, err := generateSelfSignedCertRSA()
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response for CTVersion2
				timestamp := uint64(time.Now().Unix())
				transItem := struct {
					SCTVersion   uint8
					Timestamp    uint64
					Extensions   []byte
					STHExtension []byte
				}{
					SCTVersion:   server.CTVersion2,
					Timestamp:    timestamp,
					Extensions:   []byte(""),
					STHExtension: []byte(""),
				}
				transItemBytes, _ := sonic.Marshal(transItem)

				// Hash the data before signing
				hasher := sha256.New()
				hasher.Write(transItemBytes)
				hashedData := hasher.Sum(nil)

				signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashedData)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}

				sctResponse := server.SCTResponse{
					SCTVersion:   server.CTVersion2,
					ID:           "test-ct-log",
					Timestamp:    timestamp,
					Extensions:   "",
					STHExtension: "",
					Signature:    base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and RSA private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verification & Certificate Transparency submitted successfully
		t.Log("Hello Crypto: Certificate submitted to CT log successfully")
	})

	// Test case 8: Invalid signature decoding
	t.Run("InvalidSignatureDecoding", func(t *testing.T) {
		// Generate a self-signed certificate with a valid ECDSA private key
		cert, privateKey, err := generateSelfSignedCertECDSA()
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with an invalid base64-encoded signature
				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  uint64(time.Now().Unix()),
					Extensions: "",
					Signature:  "invalid-base64-signature",
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := "failed to decode signature: illegal base64 data at input byte 7"
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 9: failed verify ECDSA signature
	t.Run("FailedVerifyECDSASignature", func(t *testing.T) {
		// Generate a self-signed certificate with a valid ECDSA private key
		cert, privateKey, err := generateSelfSignedCertECDSA()
		if err != nil {
			t.Fatal(err)
		}
		// Generate a valid ECDSA private key
		privateKeyx, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("Failed to generate ECDSA private key: %v", err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response
				// Calculate the SHA-256 hash of the certificate
				hash := sha256.Sum256(cert.Raw)

				// Prepare the mock response with a valid SCT response
				timestamp := uint64(time.Now().Unix())
				data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

				signature, err := ecdsa.SignASN1(rand.Reader, privateKeyx, data)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}
				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and the mock private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := "failed to verify ECDSA signature"
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 10: failed verify RSA signature
	t.Run("FailedVerifyRSASignature", func(t *testing.T) {
		// Generate a self-signed certificate with a valid RSA private key
		cert, privateKey, err := generateSelfSignedCertRSA()
		if err != nil {
			t.Fatal(err)
		}
		// Generate a valid ECDSA private key
		privateKeyx, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("Failed to generate ECDSA private key: %v", err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response
				// Calculate the SHA-256 hash of the certificate
				hash := sha256.Sum256(cert.Raw)

				// Prepare the mock response with a valid SCT response
				timestamp := uint64(time.Now().Unix())
				data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

				signature, err := ecdsa.SignASN1(rand.Reader, privateKeyx, data)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}
				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and the mock private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := "failed to verify RSA signature: crypto/rsa: verification error"
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 11: Failed to Encode certificate
	t.Run("FailedEncodeCertificate", func(t *testing.T) {
		// Generate a self-signed certificate with an Unexpected Key
		cert := &x509.Certificate{
			PublicKey: struct{}{}, // Unexpected Key
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response
				timestamp := uint64(time.Now().Unix())
				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString([]byte("dummy-signature")),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and a dummy private key
		err := fiberServer.SubmitToCTLog(cert, struct{}{}, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := "failed to encode certificate: x509: certificate private key does not implement crypto.Signer"
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 12: Invalid Timestamp
	t.Run("InvalidTimestamp", func(t *testing.T) {
		// Generate a self-signed certificate with a valid ECDSA private key
		cert, privateKey, err := generateSelfSignedCertECDSA()
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response
				// Calculate the SHA-256 hash of the certificate
				hash := sha256.Sum256(cert.Raw)

				// Prepare the mock response with a valid SCT response
				timestamp := uint64(9999)
				data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

				signature, err := ecdsa.SignASN1(rand.Reader, privateKey, data)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}
				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and the mock private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := "invalid timestamp: 9999"
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// Test case 13: Successful submission to CT log with Ed25519 key
	t.Run("SuccessfulSubmissionEd25519CTVersion1", func(t *testing.T) {
		// Generate a self-signed certificate with a valid Ed25519 private key
		cert, privateKey, err := generateSelfSignedCertEd25519()
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Calculate the SHA-256 hash of the certificate
				hash := sha256.Sum256(cert.Raw)

				// Prepare the mock response with a valid SCT response
				timestamp := uint64(time.Now().Unix())
				data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

				signature := ed25519.Sign(privateKey, data)

				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and Ed25519 private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verification & Certificate Transparency submitted successfully
		t.Log("Hello Crypto: Certificate submitted to CT log successfully")
	})

	// Test case 14: Successful submission to CT log with Ed25519 key and CTVersion2
	t.Run("SuccessfulSubmissionEd25519CTVersion2", func(t *testing.T) {
		// Generate a self-signed certificate with a valid Ed25519 private key
		cert, privateKey, err := generateSelfSignedCertEd25519()
		if err != nil {
			t.Fatal(err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response for CTVersion2
				timestamp := uint64(time.Now().Unix())
				transItem := struct {
					SCTVersion   uint8
					Timestamp    uint64
					Extensions   []byte
					STHExtension []byte
				}{
					SCTVersion:   server.CTVersion2,
					Timestamp:    timestamp,
					Extensions:   []byte(""),
					STHExtension: []byte(""),
				}
				transItemBytes, _ := sonic.Marshal(transItem)

				signature := ed25519.Sign(privateKey, transItemBytes)

				sctResponse := server.SCTResponse{
					SCTVersion:   server.CTVersion2,
					ID:           "test-ct-log",
					Timestamp:    timestamp,
					Extensions:   "",
					STHExtension: "",
					Signature:    base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and Ed25519 private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verification & Certificate Transparency submitted successfully
		t.Log("Hello Crypto: Certificate submitted to CT log successfully")
	})

	// Test case 15: failed verify Ed25519 signature
	t.Run("FailedVerifyEd25519Signature", func(t *testing.T) {
		// Generate a self-signed certificate with a valid Ed25519 private key
		cert, privateKey, err := generateSelfSignedCertEd25519()
		if err != nil {
			t.Fatal(err)
		}
		// Generate a valid ECDSA private key
		privateKeyx, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("Failed to generate ECDSA private key: %v", err)
		}

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response
				// Calculate the SHA-256 hash of the certificate
				hash := sha256.Sum256(cert.Raw)

				// Prepare the mock response with a valid SCT response
				timestamp := uint64(time.Now().Unix())
				data := append(hash[:], []byte(fmt.Sprintf("%d", timestamp))...)

				signature, err := ecdsa.SignASN1(rand.Reader, privateKeyx, data)
				if err != nil {
					t.Fatalf("Failed to generate signature: %v", err)
				}
				sctResponse := server.SCTResponse{
					SCTVersion: server.CTVersion1,
					ID:         "test-ct-log",
					Timestamp:  timestamp,
					Extensions: "",
					Signature:  base64.StdEncoding.EncodeToString(signature),
				}
				responseBody, _ := sonic.Marshal(sctResponse)
				mockResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}
				return mockResponse, nil
			},
		}

		// Replace the MakeHTTPRequestFunc with the mock implementation
		httpRequestMaker.MakeHTTPRequestFunc = func(req *http.Request) (*http.Response, error) {
			return mockHTTPClient.Do(req)
		}

		// Call the SubmitToCTLog method with the HTTPRequestMaker and the mock private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := "failed to verify Ed25519 signature"
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

}
