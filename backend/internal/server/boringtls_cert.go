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
	if err := verifySCT(response, hash, cert); err != nil {
		return err
	}

	return nil
}

// verifySCT verifies the signed certificate timestamp (SCT).
func verifySCT(response SCTResponse, hash [32]byte, cert *x509.Certificate) error {
	if response.SCTVersion != 0 {
		return fmt.Errorf("unsupported SCT version: %d", response.SCTVersion)
	}

	// Decode the base64-encoded signature
	signature, err := base64.StdEncoding.DecodeString(response.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	// Verify the timestamp
	if !verifyTimestamp(response.Timestamp) {
		return fmt.Errorf("invalid timestamp: %d", response.Timestamp)
	}

	// Verify the signature
	data := append(hash[:], []byte(fmt.Sprintf("%d", response.Timestamp))...)
	if err := cert.CheckSignature(x509.SHA256WithRSA, data, signature); err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	return nil
}

// verifyTimestamp checks if the given timestamp is within a valid range.
func verifyTimestamp(timestamp uint64) bool {
	now := time.Now().Unix()
	return timestamp >= uint64(now-24*60*60) && timestamp <= uint64(now+24*60*60)
}
