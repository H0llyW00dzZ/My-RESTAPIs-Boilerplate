// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CTLog represents a Certificate Transparency log.
// It contains the URL of the log server.
type CTLog struct {
	URL string
}

// SCTResponse represents the response from the CT log server.
type SCTResponse struct {
	SCTVersion uint8  `json:"sct_version"`
	ID         string `json:"id"`
	Timestamp  uint64 `json:"timestamp"`
	Extensions string `json:"extensions"`
	Signature  string `json:"signature"`
}

const (
	// CTPath is the path for submitting certificates to the CT log server.
	// It represents the API endpoint for adding a certificate chain to the log.
	// The path is typically in the format "/ct/v1/add-chain".
	CTPath = "/ct/v1/add-chain"

	// ContentTypeJSON represents the content type for JSON data.
	// It is used in the HTTP request header to specify that the request body contains JSON data.
	// The value is set to "application/json" using the [fiber.MIMEApplicationJSON] constant from the Fiber framework.
	ContentTypeJSON = fiber.MIMEApplicationJSON

	// ContentType represents the key for the Content-Type header in an HTTP request or response.
	// It is used to specify the media type of the resource being sent or received.
	// The value is set to "Content-Type" using the [fiber.HeaderContentType] constant from the Fiber framework.
	ContentType = fiber.HeaderContentType
)

// SubmitToCTLog submits the given certificate to the specified Certificate Transparency log.
//
// Example Usage:
//
//	// Load the certificate from a PEM-encoded file
//	certPEM, err := os.ReadFile("path/to/certificate.pem")
//	if err != nil {
//		log.Fatalf("failed to read certificate file: %v", err)
//	}
//
//	// Decode the PEM-encoded certificate
//	block, _ := pem.Decode(certPEM)
//	if block == nil || block.Type != "CERTIFICATE" {
//		log.Fatal("failed to decode PEM-encoded certificate")
//	}
//
//	// Parse the X.509 certificate
//	cert, err := x509.ParseCertificate(block.Bytes)
//	if err != nil {
//		log.Fatalf("failed to parse certificate: %v", err)
//	}
//
//	// Define the CT log server URL
//	ctLog := server.CTLog{
//		URL: "https://ct.example.com",
//	}
//
//	// Submit the certificate to the CT log
//	err = server.SubmitToCTLog(cert, ctLog)
//	if err != nil {
//		log.Fatalf("failed to submit certificate to CT log: %v", err)
//	}
//	 fmt.Println("Certificate submitted to CT log successfully")
//
// Note: Currently unused and marked as TODO.
func (s *FiberServer) SubmitToCTLog(cert *x509.Certificate, ctLog CTLog) error {
	// Encode the certificate in DER format
	certDER, err := x509.CreateCertificate(nil, cert, cert, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %v", err)
	}

	// Calculate the SHA-256 hash of the certificate
	hash := sha256.Sum256(certDER)

	// Create the JSON payload for submitting the certificate to the CT log
	payload := struct {
		Chain []string `json:"chain"`
	}{
		Chain: []string{base64.StdEncoding.EncodeToString(certDER)},
	}

	// Marshal the JSON payload
	// Note: Reusable, instead of multiple calls to json encoder/decoder, following DRY (Don't Repeat Yourself)
	jsonPayload, err := s.app.Config().JSONEncoder(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Create the HTTP request to submit the certificate to the CT log
	req, err := http.NewRequest(http.MethodPost, ctLog.URL+CTPath, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set(ContentType, ContentTypeJSON)

	// Send the HTTP request using the helper function
	resp, err := s.makeHTTPRequest(req)
	if err != nil {
		return fmt.Errorf("failed to submit certificate to CT log: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response body
	var response SCTResponse

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Note: Reusable, instead of multiple calls to json encoder/decoder, following DRY (Don't Repeat Yourself)
	if err := s.app.Config().JSONDecoder(responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Verify the signed certificate timestamp (SCT)
	sctVerifier := &SCTVerifier{
		Response: response,
		Hash:     hash,
		Cert:     cert,
	}
	if err := sctVerifier.VerifySCT(); err != nil {
		return err
	}

	return nil
}

// SCTVerifier represents a verifier for signed certificate timestamps (SCTs).
type SCTVerifier struct {
	Response SCTResponse
	Hash     [32]byte
	Cert     *x509.Certificate
}

// VerifySCT verifies the signed certificate timestamp (SCT).
func (v *SCTVerifier) VerifySCT() error {
	if v.Response.SCTVersion != 0 {
		return fmt.Errorf("unsupported SCT version: %d", v.Response.SCTVersion)
	}

	// Decode the base64-encoded signature
	signature, err := base64.StdEncoding.DecodeString(v.Response.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	// Verify the timestamp
	if !v.VerifyTimestamp() {
		return fmt.Errorf("invalid timestamp: %d", v.Response.Timestamp)
	}

	// Verify the signature
	data := append(v.Hash[:], []byte(fmt.Sprintf("%d", v.Response.Timestamp))...)
	if err := v.Cert.CheckSignature(v.Cert.SignatureAlgorithm, data, signature); err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	return nil
}

// VerifyTimestamp checks if the timestamp in the SCT response is within a valid range.
func (v *SCTVerifier) VerifyTimestamp() bool {
	now := time.Now().Unix()
	return v.Response.Timestamp >= uint64(now-24*60*60) && v.Response.Timestamp <= uint64(now+24*60*60)
}
