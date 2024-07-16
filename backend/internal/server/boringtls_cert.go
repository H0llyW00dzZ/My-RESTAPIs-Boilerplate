// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package server

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
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
	// Note: The [RandTLS()] function provides sufficient randomness for the purposes of [x509.CreateCertificate],
	// including generating serial numbers and signing certificates with any type of key,
	// instead of multiple calls to [io.Reader], following DRY (Don't Repeat Yourself).
	certDER, err := x509.CreateCertificate(RandTLS(), cert, cert, publicKey(privateKey), privateKey)
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %v", err)
	}

	// Calculate the SHA-256 hash of the certificate
	// TODO: Do we really need to improve this to make it more flexible (e.g., if the certificate does not use SHA-256)?
	h := sha256.Sum256(certDER)

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
	jsonConfig := json{
		Marshal: s.App.Config().JSONEncoder,
	}
	// Verify the signed certificate timestamp (SCT)
	sctVerifier := &SCTVerifier{
		Response: response,
		Hash:     h,
		Cert:     cert,
		json:     jsonConfig,
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
	json
}

// json is a struct that holds the JSON encoding/decoding configuration.
// It provides a way to customize the JSON encoding/decoding behavior by specifying
// a custom Marshal/Unmarshal function.
type json struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

// VerifySCT verifies the signed certificate timestamp (SCT).
func (v *SCTVerifier) VerifySCT() error {
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
	switch v.Response.SCTVersion {
	case CTVersion1:
		data = v.constructTransmissionItemV1()
	case CTVersion2:
		data, err = v.constructTransmissionItemV2()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported SCT version: %d", v.Response.SCTVersion)
	}

	if err := v.verifySignature(data, signature); err != nil {
		return err
	}

	return nil
}

// VerifyTimestamp checks if the timestamp in the SCT response is within a valid range.
func (v *SCTVerifier) VerifyTimestamp() bool {
	now := time.Now().Unix()
	sctTime := int64(v.Response.Timestamp)

	// Check if the SCT timestamp is within ±24 hours from the current time
	if sctTime < now-24*60*60 || sctTime > now+24*60*60 {
		return false
	}

	// Check if the SCT timestamp is within the validity period of the certificate
	if sctTime < v.Cert.NotBefore.Unix() || sctTime > v.Cert.NotAfter.Unix() {
		return false
	}

	return true
}

// publicKey returns the public key associated with the given private key.
func publicKey(priv crypto.PrivateKey) crypto.PublicKey {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

// constructTransmissionItemV1 constructs the transmission Item structure for SCT version 1.
func (v *SCTVerifier) constructTransmissionItemV1() []byte {
	var data []byte
	data = append(data, v.Cert.Raw...)
	data = append(data, byte(v.Response.SCTVersion))
	data = append(data, []byte(v.Response.ID)...)
	data = append(data, byte(v.Response.Timestamp>>56),
		byte(v.Response.Timestamp>>48),
		byte(v.Response.Timestamp>>40),
		byte(v.Response.Timestamp>>32),
		byte(v.Response.Timestamp>>24),
		byte(v.Response.Timestamp>>16),
		byte(v.Response.Timestamp>>8),
		byte(v.Response.Timestamp))
	data = append(data, []byte(v.Response.Extensions)...)

	return data
}

// constructTransmissionItemV2 constructs the transmission Item structure for SCT version 2.
func (v *SCTVerifier) constructTransmissionItemV2() ([]byte, error) {
	transMissionItem := struct {
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

	transmissionItemBytes, err := v.json.Marshal(transMissionItem)
	if err != nil {
		return nil, fmt.Errorf("failed to encode TransMissionItem: %v", err)
	}

	return transmissionItemBytes, nil
}

// verifySignature verifies the signature based on the public key type.
func (v *SCTVerifier) verifySignature(data, signature []byte) error {
	switch publicKey := v.Cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		if !ecdsa.VerifyASN1(publicKey, data, signature) {
			return errors.New("failed to verify ECDSA signature")
		}
	case *rsa.PublicKey:
		if err := v.verifyRSASignature(publicKey, data, signature); err != nil {
			return err
		}
	case ed25519.PublicKey:
		if !ed25519.Verify(publicKey, data, signature) {
			return errors.New("failed to verify Ed25519 signature")
		}
	default:
		return fmt.Errorf("unsupported public key type: %T", v.Cert.PublicKey)
	}

	return nil
}

// verifyRSASignature verifies the RSA signature.
func (v *SCTVerifier) verifyRSASignature(publicKey *rsa.PublicKey, data, signature []byte) error {
	h := sha256.New()
	h.Write(data)
	hashedData := h.Sum(nil)

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedData, signature); err != nil {
		return fmt.Errorf("failed to verify RSA signature: %v", err)
	}

	return nil
}
