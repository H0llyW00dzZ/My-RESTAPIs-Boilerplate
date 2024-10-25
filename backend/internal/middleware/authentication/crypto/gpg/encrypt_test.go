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

const testPublicECDSACantEncrypt = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mFIEZxh27BMIKoZIzj0DAQcCAwRw2BIEuz/lUbsWB11eKNDzDTS86SU8t5S1+WhL
PnWxuW8ylRjIaLzv6QRs0idiagE9dLVdpm9XwVhojyOCx91mtCRUZXN0IEtleShU
ZXN0IEtleSk8dGVzdEBleGFtcGxlLmNvbT6IkwQTEwgAOxYhBJoZ+uA5zgwcjmzC
3sydPySjCpmPBQJnGHbsAhsDBQsJCAcCAiICBhUKCQgLAgQWAgMBAh4HAheAAAoJ
EMydPySjCpmP548A/3cKzb/YjiNPH5NOQvVeizEuU2Jo8ZBgK52JuVpqxakrAQDP
lQD3Q4dlnY9UeRlO+wvaMYtg/y9UCpdBWG8qrxyMOw==
=zFbO
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
	gpg, err := gpg.NewEncryptor([]string{testPublicKey})
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

	gpg, err := gpg.NewEncryptor([]string{testPublicKey})
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

func TestNewEncryptorWithInvalidKey(t *testing.T) {
	_, err := gpg.NewEncryptor([]string{testPublicECDSACantEncrypt})
	if err == nil {
		t.Fatalf("Expected error when creating encryptor with a key that cannot encrypt, but got none")
	}

	if err != gpg.ErrorCantEncrypt {
		t.Fatalf("Expected ErrorCantEncrypt, but got: %v", err)
	}
}

func TestGetKeyInfos(t *testing.T) {
	// Sample public key
	publicKeys := []string{
		testPublicKey,
		testPublicECDSACantEncrypt,
	}

	// Create an Encryptor instance
	gpg, err := gpg.NewEncryptor(publicKeys)
	if err != nil {
		t.Fatalf("Failed to create Encryptor: %v", err)
	}

	// Get key infos
	keyInfos := gpg.GetKeyInfos()

	// Check that keyInfos is not empty
	if len(keyInfos) == 0 {
		t.Fatal("Expected keyInfos to contain key metadata, but it was empty")
	}

	// Log detailed key information
	for i, info := range keyInfos {
		t.Logf("Key %d:", i+1)
		t.Logf("KeyID: %d", info.KeyID)
		t.Logf("Hex KeyID: %s", info.HexKeyID)
		t.Logf("CanEncrypt: %t", info.CanEncrypt)
		t.Logf("CanVerify: %t", info.CanVerify)
		t.Logf("IsExpired: %t", info.IsExpired)
		t.Logf("IsRevoked: %t", info.IsRevoked)
		t.Logf("Key Fingerprints: %s", info.Fingerprint)
		t.Logf("Digest Fingerprints: %v", info.DigestFingerprint)
	}

	// Example check: Verify the first key's CanEncrypt field
	if !keyInfos[0].CanEncrypt {
		t.Fatal("Expected first key to be capable of encryption")
	}

	// Additional checks can be added based on expected key metadata
}

// Use this for local development, such as testing with different GPG keys.
func TestEncryptStreamToFile(t *testing.T) {
	// Sample public key
	publicKeys := []string{
		testPublicKey,
	}

	// Create a buffer to simulate the input data
	inputData := []byte("Hello GPG/OpenPGP From H0llyW00dzZ.")
	inputBuffer := bytes.NewReader(inputData)

	// Define the output file
	outputFile, err := os.CreateTemp("", "test_output_*.gpg")
	if err != nil {
		t.Fatalf("Failed to create temporary output file: %v", err)
	}
	defer outputFile.Close()
	// Note: Do not defer os.Remove(outputFile.Name()) to keep the file for decryption testing

	// Create an Encryptor instance
	encryptor, err := gpg.NewEncryptor(publicKeys)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	// Call the EncryptStream function
	if err = encryptor.EncryptStream(inputBuffer, outputFile); err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Check if the output file has data
	fileInfo, err := outputFile.Stat()
	if err != nil {
		t.Fatalf("Failed to get output file info: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Fatalf("Output file is empty")
	}

	// Log the name of the output file for reference
	t.Logf("Encrypted data written to file: %s", outputFile.Name())

	// Optionally, add decryption and verification logic here
}
