// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
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
	SCTVersion   uint8  `json:"sct_version"`
	ID           string `json:"id"`
	Timestamp    uint64 `json:"timestamp"`
	Extensions   string `json:"extensions"`
	Signature    string `json:"signature"`
	STHExtension string `json:"sth_extension,omitempty"`
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

const (
	// CTVersion1 represents version 1 of the Certificate Transparency (CT) protocol.
	// It is defined using an iota constant, allowing for easy extensibility and readability.
	CTVersion1 uint8 = iota + 1 // RFC 6962

	// CTVersion2 represents version 2 of the Certificate Transparency (CT) protocol.
	// It is automatically assigned the next value in the iota sequence.
	CTVersion2 // RFC 9162

	// LatestCTVersion represents the latest version of the Certificate Transparency (CT) protocol.
	// It should be updated whenever a new version is added.
	LatestCTVersion = CTVersion2
)

// SubmitToCTLog submits the given certificate to the specified Certificate Transparency log.
//
// The function takes the following parameters:
//   - cert: The X.509 certificate to be submitted to the CT log.
//   - privateKey: The private key associated with the certificate.
//   - ctLog: The CTLog struct representing the CT log server.
//   - httpRequestMaker: An optional HTTPRequestMaker instance for making HTTP requests.
//
// The function performs the following steps:
//  1. Encodes the certificate in DER format.
//  2. Calculates the SHA-256 hash of the certificate.
//  3. Creates a JSON payload containing the base64-encoded certificate.
//  4. Sends an HTTP POST request to the CT log server with the JSON payload.
//  5. Parses the response from the CT log server.
//  6. Verifies the signed certificate timestamp (SCT) received in the response.
//
// If the submission is successful and the SCT is valid, the function returns nil.
// If an error occurs during the submission or verification process, the function returns an error.
//
// Example Usage:
//
//	// Load the certificate and private key
//	certPEM, err := os.ReadFile("path/to/certificate.pem")
//	if err != nil {
//		log.Fatalf("failed to read certificate file: %v", err)
//	}
//	privateKeyPEM, err := os.ReadFile("path/to/private_key.pem")
//	if err != nil {
//		log.Fatalf("failed to read private key file: %v", err)
//	}
//
//	// Decode the PEM-encoded certificate and private key
//	certBlock, _ := pem.Decode(certPEM)
//	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
//		log.Fatal("failed to decode PEM-encoded certificate")
//	}
//	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
//	if privateKeyBlock == nil || privateKeyBlock.Type != "PRIVATE KEY" {
//		log.Fatal("failed to decode PEM-encoded private key")
//	}
//
//	// Parse the X.509 certificate and private key
//	cert, err := x509.ParseCertificate(certBlock.Bytes)
//	if err != nil {
//		log.Fatalf("failed to parse certificate: %v", err)
//	}
//	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
//	if err != nil {
//		log.Fatalf("failed to parse private key: %v", err)
//	}
//
//	// Define the CT log server URL
//	ctLog := server.CTLog{
//		URL: "https://ct.example.com",
//	}
//
//	// Create a Fiber server instance
//	app := fiber.New()
//	fiberServer := &server.FiberServer{App: app}
//
//	// Submit the certificate to the CT log
//	err = fiberServer.SubmitToCTLog(cert, privateKey, ctLog, nil)
//	if err != nil {
//		log.Fatalf("failed to submit certificate to CT log: %v", err)
//	}
//	fmt.Println("Certificate submitted to CT log successfully")
//
// Note: Currently unused because it's boring to submit Certificate Transparency logs, unlike implementing a Cryptographic Protocol.
func (s *FiberServer) SubmitToCTLog(cert *x509.Certificate, privateKey crypto.PrivateKey, ctLog CTLog, httpRequestMaker *HTTPRequestMaker) error {
	// Encode the certificate in DER format
	certDER, err := x509.CreateCertificate(rand.Reader, cert, cert, publicKey(privateKey), privateKey)
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %v", err)
	}

	// Calculate the SHA-256 hash of the certificate
	// TODO: Do we really need to improve this to make it more flexible (e.g., if the certificate does not use SHA-256)?
	hash := sha256.Sum256(certDER)

	// Create the JSON payload for submitting the certificate to the CT log
	payload := struct {
		Chain []string `json:"chain"`
	}{
		Chain: []string{base64.StdEncoding.EncodeToString(certDER)},
	}

	// Marshal the JSON payload
	// Note: Reusable, instead of multiple calls to json encoder/decoder, following DRY (Don't Repeat Yourself)
	jsonPayload, err := s.App.Config().JSONEncoder(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Create the HTTP request to submit the certificate to the CT log
	req, err := http.NewRequest(http.MethodPost, ctLog.URL+CTPath, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set(ContentType, ContentTypeJSON)

	// Send the HTTP request using the helper function or MakeHTTPRequest directly
	var resp *http.Response
	if httpRequestMaker != nil {
		resp, err = httpRequestMaker.MakeHTTPRequest(req)
	} else {
		resp, err = s.MakeHTTPRequest(req)
	}
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
	if err := s.App.Config().JSONDecoder(responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Verify the signed certificate timestamp (SCT)
	sctVerifier := &SCTVerifier{
		Response: response,
		Hash:     hash,
		Cert:     cert,
	}
	if err := sctVerifier.VerifySCT(s.App.Config().JSONEncoder); err != nil {
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
func (v *SCTVerifier) VerifySCT(jsonEncoder func(v any) ([]byte, error)) error {
	// Note: This is a method Go idiom that uses the constant iota sequence.
	// It is particularly useful in cryptographic operations (e.g., implementing custom ciphers, custom protocols, or any cryptography-related tasks).
	if v.Response.SCTVersion < CTVersion1 || v.Response.SCTVersion > LatestCTVersion {
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

	// Calculate the hash of the certificate in DER format
	// TODO: Do we really need to improve this to make it more flexible (e.g., if the certificate does not use SHA-256)?
	var data []byte
	hash := sha256.Sum256(v.Cert.Raw)
	switch v.Response.SCTVersion {
	case CTVersion1:
		data = append(hash[:], []byte(fmt.Sprintf("%d", v.Response.Timestamp))...)
		// Note: When there is another version (e.g., Version 3), this Version 2 logic should be extracted
		// into a separate function to keep the code simple and maintainable.
	case CTVersion2:
		// Construct the TransItem structure for signature verification
		transItem := struct {
			SCTVersion   uint8
			Timestamp    uint64
			Extensions   []byte
			STHExtension []byte
		}{
			SCTVersion:   v.Response.SCTVersion,
			Timestamp:    v.Response.Timestamp,
			Extensions:   []byte(v.Response.Extensions),
			STHExtension: []byte(v.Response.STHExtension),
		}

		// Encode the TransItem structure using Fiber's JSON encoding configuration
		transItemBytes, err := jsonEncoder(transItem)
		if err != nil {
			return fmt.Errorf("failed to encode TransItem: %v", err)
		}
		data = transItemBytes
	default:
		return fmt.Errorf("unsupported SCT version: %d", v.Response.SCTVersion)
	}

	// Verify the signature based on the public key type
	switch publicKey := v.Cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		if !ecdsa.VerifyASN1(publicKey, data, signature) {
			return fmt.Errorf("failed to verify ECDSA signature")
		}
	case *rsa.PublicKey:
		// Hash the data before verifying the RSA signature
		// TODO: Do we really need to improve this to make it more flexible (e.g., if the certificate does not use SHA-256)?
		hasher := sha256.New()
		hasher.Write(data)
		hashedData := hasher.Sum(nil)

		if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedData, signature); err != nil {
			return fmt.Errorf("failed to verify RSA signature: %v", err)
		}
	default:
		return fmt.Errorf("unsupported public key type: %T", v.Cert.PublicKey)
	}

	return nil
}

// VerifyTimestamp checks if the timestamp in the SCT response is within a valid range.
func (v *SCTVerifier) VerifyTimestamp() bool {
	now := time.Now().Unix()
	return v.Response.Timestamp >= uint64(now-24*60*60) && v.Response.Timestamp <= uint64(now+24*60*60)
}

// publicKey returns the public key associated with the given private key.
func publicKey(priv crypto.PrivateKey) crypto.PublicKey {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
