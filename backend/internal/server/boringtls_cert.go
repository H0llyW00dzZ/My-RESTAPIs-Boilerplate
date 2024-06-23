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
)

// CTLog represents a Certificate Transparency log.
// It contains the URL of the log server.
type CTLog struct {
	URL string
}

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
	req, err := http.NewRequest("POST", ctLog.URL+"/ct/v1/add-chain", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit certificate to CT log: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response body
	var response struct {
		SCTVersion uint8  `json:"sct_version"`
		ID         string `json:"id"`
		Timestamp  uint64 `json:"timestamp"`
		Extensions string `json:"extensions"`
		Signature  string `json:"signature"`
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Note: Reusable, instead of multiple calls to json encoder/decoder, following DRY (Don't Repeat Yourself)
	if err := s.app.Config().JSONDecoder(responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Verify the signed certificate timestamp (SCT)
	if response.SCTVersion != 0 {
		return fmt.Errorf("unsupported SCT version: %d", response.SCTVersion)
	}

	// Decode the base64-encoded signature
	signature, err := base64.StdEncoding.DecodeString(response.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	// Verify the timestamp
	if response.Timestamp < uint64(time.Now().Add(-24*time.Hour).Unix()) ||
		response.Timestamp > uint64(time.Now().Add(24*time.Hour).Unix()) {
		return fmt.Errorf("invalid timestamp: %d", response.Timestamp)
	}

	// Verify the signature
	data := append(hash[:], []byte(fmt.Sprintf("%d", response.Timestamp))...)
	if err := cert.CheckSignature(x509.SHA256WithRSA, data, signature); err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	return nil
}
