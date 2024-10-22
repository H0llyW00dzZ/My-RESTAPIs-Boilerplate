// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg_test

import (
	"bytes"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/gpg"
	"os"
	"testing"
)

// Sample PGP/GPG keys for testing (RFC 9580) Sections 12.7, 5.2.3.4, and 11.5 Latest strong mechanisms for GPG/OpenPGP.
//
// KEY:
//
// - https://keys.openpgp.org/search?q=95F9A1D43F57344AB88BFFFEA0F9424A7002343A
//
// REST APIs GPG Proton Lookup (created by H0llyW00dzZ):
//
//	curl -X POST https://api.b0zal.io/v1/gpg/proton/lookup \
//	-H "Content-Type: application/json" \
//	-d '{"email":"H0llyW00dzZ@pm.me"}'
//
// Note: If you attempt to look up the GPG Proton Public Key using the REST API and receive a 403 Forbidden response,
// it means your IP network has been blocked due to suspicious activity (e.g., your network might be compromised, such as by a botnet).
// My firewall mechanism is precise in identifying normal users, bots, or infected devices.
const testPublicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEZhww9xYJKwYBBAHaRw8BAQdAA9nmVRaTTKJe7EDCQ8OhshfDim+9kjCpbUU6
dSsYkfi0JWgwbGx5dzAwZHp6QHBtLm1lIDxoMGxseXcwMGR6ekBwbS5tZT6IjAQQ
FgoAPgWCZhww9wQLCQcICZCg+UJKcAI0OgMVCAoEFgACAQIZAQKbAwIeARYhBJX5
odQ/VzRKuIv//qD5QkpwAjQ6AACUggD+Pm+exMl9WgD7ignm/nW4HXYCyaGe7ZBF
pILgsOh96twA/122jRFkH5bzcbRjIGuL+9+Nr+69cnuBBtAJNfNFelYPuDgEZhww
9xIKKwYBBAGXVQEFAQEHQI55aMA1TdV6P/DNh+/TMb3bb1jN7bAlha3HRs5BB9dD
AwEIB4h4BBgWCgAqBYJmHDD3CZCg+UJKcAI0OgKbDBYhBJX5odQ/VzRKuIv//qD5
QkpwAjQ6AABELAD/YG153FordpFJMJTI8OEzAvZwRxAvszdvPAMzqI+BSlYBAIBj
zAozXAC69DgM8AOJzEnsiA55ic1D56y64baz31cD
=m5PK
-----END PGP PUBLIC KEY BLOCK-----
`

func TestEncryptFile(t *testing.T) {
	// Create a temporary file to encrypt
	inputFile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Write some data to the input file
	_, err = inputFile.WriteString("Hello GPG/OpenPGP From H0llyW00dzZ.")
	if err != nil {
		t.Fatalf("Failed to write to input file: %v", err)
	}
	inputFile.Close()

	// Define the output file
	outputFile := inputFile.Name() + ".gpg"
	defer os.Remove(outputFile)

	// Encrypt the backup file
	gpg, err := gpg.NewEncryptor(testPublicKey)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	// Call the EncryptFile function
	if err = gpg.EncryptFile(inputFile.Name(), outputFile); err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	// Check if the output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}
}

func TestEncryptStream(t *testing.T) {
	// Create a buffer to simulate the input file
	inputData := []byte("Hello GPG/OpenPGP From H0llyW00dzZ.")
	inputBuffer := bytes.NewReader(inputData)

	// Create a buffer to simulate the output file
	outputBuffer := &bytes.Buffer{}

	gpg, err := gpg.NewEncryptor(testPublicKey)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	// Call the EncryptStream function
	if err = gpg.EncryptStream(inputBuffer, outputBuffer); err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Check if the output buffer has data
	if outputBuffer.Len() == 0 {
		t.Fatalf("Output buffer is empty")
	}

	// Compare original and encrypted data
	if bytes.Equal(inputData, outputBuffer.Bytes()) {
		t.Fatalf("Encrypted data is the same as original data")
	}

	// Optionally, you can add more checks to see if the data is encrypted
	// This would typically involve decrypting with a private key and verifying the content
}
