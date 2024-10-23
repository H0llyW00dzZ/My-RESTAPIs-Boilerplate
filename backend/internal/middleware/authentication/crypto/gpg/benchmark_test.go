// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gpg_test

import (
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/gpg"
	"os"
	"testing"
)

// Average times on my laptop without overclocking:
//
//	goos: windows
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/gpg
//	cpu: Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz
//	BenchmarkEncryptLargeFile-16    	       1	4775101100 ns/op	 2476872 B/op	    6585 allocs/op
//
// Note that it's arround 4 ~ 5 seconds for 1GiB
func BenchmarkEncryptLargeFile(b *testing.B) {
	// Create a temporary file to encrypt
	inputFile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		b.Fatalf("Failed to create temporary input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Write 1 GiB of data to the input file.
	//
	// Note: During benchmark testing, memory allocation is used. However, in production, memory usage should be minimal,
	// even for large data sizes (e.g., 250 MiB+ backup SQL), with memory usage around 15-16 MiB.
	// This efficiency is achieved by streaming data directly from the file/disk, rather than holding it in memory.
	const size = 1 << 30 // 1 GiB

	// Simulate streaming data to the file in chunks to avoid large memory allocations that might occur on other architectures hahaha
	chunkSize := 4 << 20 // 4 MiB
	data := make([]byte, chunkSize)
	for written := int64(0); written < size; written += int64(chunkSize) {
		if _, err := inputFile.Write(data); err != nil {
			b.Fatalf("Failed to write to input file: %v", err)
		}
	}
	inputFile.Close()

	// Define the output file
	outputFile := inputFile.Name() + ".gpg"
	defer os.Remove(outputFile)

	// Create the encryptor
	gpg, err := gpg.NewEncryptor([]string{testPublicKey})
	if err != nil {
		b.Fatalf("Failed to create encryptor: %v", err)
	}

	// Run the benchmark
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		if err = gpg.EncryptFile(inputFile.Name(), outputFile); err != nil {
			b.Fatalf("EncryptFile failed: %v", err)
		}
	}

	// Check if the output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		b.Fatalf("Output file was not created")
	}
}
