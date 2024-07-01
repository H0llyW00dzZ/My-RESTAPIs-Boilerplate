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
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"

	"h0llyw00dz-template/backend/internal/server"
)

// Note: This test may seem complex when scanned by gocyclo tools due to a single function performing many tasks. However, it is not actually complex because it is safe
// from multiple goroutines that call it simultaneously, even if there are 999999999999999999 goroutines. ¯\_(ツ)_/¯
func TestSubmitToCTLog(t *testing.T) {
	// Generate a self-signed certificate with a valid private key
	cert, privateKey, err := generateSelfSignedCertECDSA()
	if err != nil {
		t.Fatal(err)
	}

	// Create a test CT log
	ctLog := server.CTLog{
		URL: "https://" + testHostName,
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
				// Prepare the mock response with a valid SCT response
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
		} else {
			// Verification & Certificate Transparency submitted successfully
			t.Log("Hello Crypto: Certificate submitted to CT log successfully")
		}
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
				// Prepare the mock response with a valid SCT response
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
		} else {
			// Verification & Certificate Transparency submitted successfully
			t.Log("Hello Crypto: Certificate submitted to CT log successfully")
		}
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
			// Prepare the mock response with a valid SCT response
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
		} else {
			// Verification & Certificate Transparency submitted successfully
			t.Log("Hello Crypto: Certificate submitted to CT log successfully")
		}
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
				transMissionItem := struct {
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
				transMissionItemBytes, _ := sonic.Marshal(transMissionItem)

				signature, err := ecdsa.SignASN1(rand.Reader, privateKey, transMissionItemBytes)
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
		} else {
			// Verification & Certificate Transparency submitted successfully
			t.Log("Hello Crypto: Certificate submitted to CT log successfully")
		}
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
				transMissionItem := struct {
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
				transMissionItemBytes, _ := sonic.Marshal(transMissionItem)

				// Hash the data before signing
				hasher := sha256.New()
				hasher.Write(transMissionItemBytes)
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
		} else {
			// Verification & Certificate Transparency submitted successfully
			t.Log("Hello Crypto: Certificate submitted to CT log successfully")
		}
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
				timestamp := uint64(9999)
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
				// Prepare the mock response with a valid SCT response
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
		} else {
			// Verification & Certificate Transparency submitted successfully
			t.Log("Hello Crypto: Certificate submitted to CT log successfully")
		}
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
				transMissionItem := struct {
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
				transMissionItemBytes, _ := sonic.Marshal(transMissionItem)

				signature := ed25519.Sign(privateKey, transMissionItemBytes)

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
		} else {
			// Verification & Certificate Transparency submitted successfully
			t.Log("Hello Crypto: Certificate submitted to CT log successfully")
		}
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

	// Test case 16: failed verify Ed25519 signature due to expired certificate
	t.Run("FailedVerifyEd25519SignatureExpired", func(t *testing.T) {
		// Generate a self-signed certificate with a valid Ed25519 private key
		cert, privateKey, err := generateSelfSignedCertEd25519WithExpired()
		if err != nil {
			t.Fatal(err)
		}

		var timestamp uint64

		// Create a mock HTTP client
		mockHTTPClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Prepare the mock response with a valid SCT response for CTVersion2
				timestamp = uint64(time.Now().Unix())
				transMissionMissionItem := struct {
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
				transMissionMissionItemBytes, _ := sonic.Marshal(transMissionMissionItem)

				signature := ed25519.Sign(privateKey, transMissionMissionItemBytes)

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

		// Call the SubmitToCTLog method with the HTTPRequestMaker and the mock private key
		err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, httpRequestMaker)

		// Assert the expectations
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		expectedErrorMessage := fmt.Sprintf("invalid timestamp: %d", timestamp)
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

}
