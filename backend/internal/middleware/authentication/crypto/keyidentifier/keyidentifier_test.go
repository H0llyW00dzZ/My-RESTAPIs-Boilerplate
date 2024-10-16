// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	std "crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/keyidentifier"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func TestKeyIdentifier(t *testing.T) {
	// Generate an ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), std.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA private key: %v", err)
	}

	// Create a new instance of KeyIdentifier with the desired configuration
	const ecdsaUUID = "ecdsa_authorized:"
	keyIdentifier := keyidentifier.New(keyidentifier.Config{
		Prefix:           ecdsaUUID,
		PrivateKey:       privateKey,
		Digest:           sha256.New,
		SignedContextKey: "signature",
	})

	// Extract the public key from the private key
	publicKey, err := keyIdentifier.GetECDSAPubKey()
	if err != nil {
		t.Fatalf("Failed to extract public key from the private key: %v", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route handler
	app.Get("/test", func(c *fiber.Ctx) error {
		// Get the generated uuid from the context
		uuid := keyIdentifier.GetKeyFunc()(c)
		fmt.Println("Session ID Authorized:", uuid)

		// Get the signature from the context
		signature := c.Locals("signature").([]byte)
		fmt.Println("Signature:", signature)

		// Send the uuid, signature, and public key in the response
		return c.JSON(fiber.Map{
			"uuid":      uuid,
			"signature": hex.EncodeToString(signature),
			"publicKey": hex.EncodeToString(elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)),
		})
	})

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/test", nil)

	// Perform the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Parse the response body using sonic
	var body map[string]string
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Get the uuid, signature, and public key from the response body
	uuid := body["uuid"]
	signatureHex := body["signature"]
	publicKeyHex := body["publicKey"]

	// Decode the signature from hex
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		t.Fatalf("Failed to decode signature: %v", err)
	}

	// Decode the public key from hex
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	// Unmarshal the public key
	x, y := elliptic.Unmarshal(elliptic.P256(), publicKeyBytes)
	if x == nil {
		t.Fatal("Invalid public key")
	}
	publicKey = &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	// Verify the signature
	h := sha256.New()
	h.Write([]byte(uuid[len(ecdsaUUID):]))
	expectedDigest := h.Sum(nil)

	if !ecdsa.VerifyASN1(publicKey, expectedDigest, signature) {
		t.Error("Signature verification failed")
	}
}

// Note: This ECC suitable for GPG, SSH, Signed Commit, HTTPS/TLS, Signed other object, E2E
func TestKeyIdentifierWithFixedRand(t *testing.T) {
	// Generate an ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.FixedSizeECC(elliptic.P256()))
	if err != nil {
		t.Fatalf("Failed to generate ECDSA private key: %v", err)
	}

	// Create a new instance of KeyIdentifier with the desired configuration
	const ecdsaUUID = "ecdsa_fixed_rand_authorized:"
	keyIdentifier := keyidentifier.New(keyidentifier.Config{
		Prefix:           ecdsaUUID,
		PrivateKey:       privateKey,
		Digest:           sha256.New,
		SignedContextKey: "signature",
		Rand:             rand.FixedSizeECC(elliptic.P256()),
	})

	// Extract the public key from the private key
	publicKey, err := keyIdentifier.GetECDSAPubKey()
	if err != nil {
		t.Fatalf("Failed to extract public key from the private key: %v", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route handler
	app.Get("/test", func(c *fiber.Ctx) error {
		// Get the generated uuid from the context
		uuid := keyIdentifier.GetKeyFunc()(c)
		fmt.Println("Session ID Authorized:", uuid)

		// Get the signature from the context
		signature := c.Locals("signature").([]byte)
		fmt.Println("Signature:", signature)

		// Send the uuid, signature, and public key in the response
		return c.JSON(fiber.Map{
			"uuid":      uuid,
			"signature": hex.EncodeToString(signature),
			"publicKey": hex.EncodeToString(elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)),
		})
	})

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/test", nil)

	// Perform the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Parse the response body using sonic
	var body map[string]string
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Get the uuid, signature, and public key from the response body
	uuid := body["uuid"]
	signatureHex := body["signature"]
	publicKeyHex := body["publicKey"]

	// Decode the signature from hex
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		t.Fatalf("Failed to decode signature: %v", err)
	}

	// Decode the public key from hex
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	// Unmarshal the public key
	x, y := elliptic.Unmarshal(elliptic.P256(), publicKeyBytes)
	if x == nil {
		t.Fatal("Invalid public key")
	}
	publicKey = &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	// Verify the signature
	h := sha256.New()
	h.Write([]byte(uuid[len(ecdsaUUID):]))
	expectedDigest := h.Sum(nil)

	if !ecdsa.VerifyASN1(publicKey, expectedDigest, signature) {
		t.Error("Signature verification failed")
	}
}

// Note: If test fail then the fixed size it's incorrect
func TestKeyIdentifierWithFixedRandECDSA(t *testing.T) {
	// Generate an ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.FixedSizeECC(elliptic.P384()))
	if err != nil {
		t.Fatalf("Failed to generate ECDSA private key: %v", err)
	}

	// Create a new instance of KeyIdentifier with the desired configuration
	const ecdsaUUID = "ecdsa_fixed_rand_authorized:"
	keyIdentifier := keyidentifier.New(keyidentifier.Config{
		Prefix:           ecdsaUUID,
		PrivateKey:       privateKey,
		Digest:           sha512.New384,
		SignedContextKey: "signature",
		Rand:             rand.FixedSizeECC(elliptic.P384()),
	})

	// Extract the public key from the private key
	publicKey, err := keyIdentifier.GetECDSAPubKey()
	if err != nil {
		t.Fatalf("Failed to extract public key from the private key: %v", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a test route handler
	app.Get("/test", func(c *fiber.Ctx) error {
		// Get the generated uuid from the context
		uuid := keyIdentifier.GetKeyFunc()(c)
		fmt.Println("Session ID Authorized:", uuid)

		// Get the signature from the context
		signature := c.Locals("signature").([]byte)
		fmt.Println("Signature:", signature)

		// Send the uuid, signature, and public key in the response
		return c.JSON(fiber.Map{
			"uuid":      uuid,
			"signature": hex.EncodeToString(signature),
			"publicKey": hex.EncodeToString(elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)),
		})
	})

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/test", nil)

	// Perform the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Parse the response body using sonic
	var body map[string]string
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Get the uuid, signature, and public key from the response body
	uuid := body["uuid"]
	signatureHex := body["signature"]
	publicKeyHex := body["publicKey"]

	// Decode the signature from hex
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		t.Fatalf("Failed to decode signature: %v", err)
	}

	// Decode the public key from hex
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	// Unmarshal the public key
	x, y := elliptic.Unmarshal(elliptic.P384(), publicKeyBytes)
	if x == nil {
		t.Fatal("Invalid public key")
	}
	publicKey = &ecdsa.PublicKey{Curve: elliptic.P384(), X: x, Y: y}

	// Verify the signature
	h := sha512.New384()
	h.Write([]byte(uuid[len(ecdsaUUID):]))
	expectedDigest := h.Sum(nil)

	if !ecdsa.VerifyASN1(publicKey, expectedDigest, signature) {
		t.Error("Signature verification failed")
	}
}
