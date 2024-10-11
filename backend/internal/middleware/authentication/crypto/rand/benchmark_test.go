// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Test Benchmark on My Laptop
// Result:
// goos: windows
// goarch: amd64
// pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
// cpu: Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz
// BenchmarkFixedSize32Bytes-16    	 			 5556951	       224.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkFixedSizeECDSA/P-256-16         	 5560414	       222.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkFixedSizeECDSA/P-384-16         	 4748517	       257.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkFixedSizeECDSA/P-521-16         	 4839511	       259.7 ns/op	       0 B/op	       0 allocs/op
//
// Note: These benchmarks were conducted without overclocking. If overclocked, the performance may reach goes crazy up to 1 billion operations.
//
// Additionally, the results shown here are from a semi-overclocked state (3 GHz):
//
// goos: windows
// goarch: amd64
// pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
// cpu: Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz
// BenchmarkFixedSize32Bytes-16    	 6492642	       191.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkFixedSizeECDSA/P-256-16         	 6795412	       184.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkFixedSizeECDSA/P-384-16         	 5626008	       214.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkFixedSizeECDSA/P-521-16         	 5705864	       211.9 ns/op	       0 B/op	       0 allocs/op

package rand_test

import (
	"crypto/elliptic"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	"testing"
)

func BenchmarkFixedSize32Bytes(b *testing.B) {
	reader := rand.FixedSize32Bytes()
	buf := make([]byte, 32)

	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		if _, err := reader.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFixedSizeECDSA(b *testing.B) {
	curves := []elliptic.Curve{
		elliptic.P256(),
		elliptic.P384(),
		elliptic.P521(),
	}

	for _, curve := range curves {
		b.Run(curve.Params().Name, func(b *testing.B) {
			reader := rand.FixedSizeECDSA(curve)
			expectedSize := 32
			if curve.Params().BitSize > 256 {
				expectedSize = 48
			}
			buf := make([]byte, expectedSize)

			b.ResetTimer() // Reset the timer to exclude setup time
			for i := 0; i < b.N; i++ {
				if _, err := reader.Read(buf); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Result:
//
//	goos: windows
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz
//	BenchmarkFixedSizeECC/P-224-16         	 5035641	       244.4 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-256-16         	 4856640	       241.6 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-384-16         	 4361732	       288.9 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-521-16         	 3676401	       346.8 ns/op	       0 B/op	       0 allocs/op
//
// Note: These results are without overclocking. If overclocked (e.g., fully overclocked), performance may increase significantly (e.g., the op).
func BenchmarkFixedSizeECC(b *testing.B) {
	curves := []elliptic.Curve{
		elliptic.P224(),
		elliptic.P256(),
		elliptic.P384(),
		elliptic.P521(),
	}

	for _, curve := range curves {
		b.Run(curve.Params().Name, func(b *testing.B) {
			reader := rand.FixedSizeECC(curve)
			byteSize := (curve.Params().BitSize + 7) / 8
			buf := make([]byte, byteSize)

			b.ResetTimer() // Reset the timer to exclude setup time
			for i := 0; i < b.N; i++ {
				if _, err := reader.Read(buf); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
