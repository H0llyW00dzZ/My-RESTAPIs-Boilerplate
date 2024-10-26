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

// Mirror (My K8s for Sandbox Development very secure): https://git.b0zal.io/H0llyW00dzZ.gpg
const testPublicKeyfromGitea = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
=Rn5P
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
	// Sample public key
	publicKeys := []string{
		testPublicKey,
		// This key will be invalid due to GPG Proton's implementation
		// (it does not support a single armored block with multiple keys, so each must be processed individually).
		testPublicKeyfromGitea,
	}

	// Create a buffer to simulate the input file
	inputData := []byte("Hello GPG/OpenPGP From H0llyW00dzZ.")
	inputBuffer := bytes.NewReader(inputData)

	// Create a buffer to simulate the output file
	outputBuffer := &bytes.Buffer{}

	gpg, err := gpg.NewEncryptor(publicKeys)
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
		// This key will be invalid due to GPG Proton's implementation
		// (it does not support a single armored block with multiple keys, so each must be processed individually).
		testPublicKeyfromGitea,
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

	// Create a temporary file to simulate the input
	inputFile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary input file: %v", err)
	}

	// Keep this to clean up the input file
	defer os.Remove(inputFile.Name())
	defer inputFile.Close()

	// Write some data to the file
	_, err = inputFile.WriteString("Hello GPG/OpenPGP From H0llyW00dzZ.")
	if err != nil {
		t.Fatalf("Failed to write to input file: %v", err)
	}

	// Re-open the file for reading
	inputFile, err = os.Open(inputFile.Name())
	if err != nil {
		t.Fatalf("Failed to open input file for reading: %v", err)
	}

	// Define the output file
	//
	// Note: The output file should include the original file extension before ".gpg".
	// For example, if the file is a text file, use ".txt.gpg" (e.g., filename.txt.gpg).
	// This way, when a GPG frontend or other OpenPGP mechanism decrypts it, the result will be filename.txt.
	outputFile, err := os.CreateTemp("", "test_output_*.txt.gpg")
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
	if err = encryptor.EncryptStream(inputFile, outputFile); err != nil {
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

// Use this for local development, such as testing with different GPG keys.
func TestEncryptStreamFromBytesToFile(t *testing.T) {
	// Sample public key
	publicKeys := []string{
		testPublicKey,
	}

	// Create a buffer to simulate the input data
	inputData := []byte("Hello GPG/OpenPGP From H0llyW00dzZ.")
	inputBuffer := bytes.NewReader(inputData)

	// Define the output file
	//
	// Note: The output file should include the original file extension before ".gpg".
	// For example, if the file is a text file, use ".txt.gpg" (e.g., filename.txt.gpg).
	// This way, when a GPG frontend or other OpenPGP mechanism decrypts it, the result will be filename.txt.
	outputFile, err := os.CreateTemp("", "test_output_*.txt.gpg")
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
