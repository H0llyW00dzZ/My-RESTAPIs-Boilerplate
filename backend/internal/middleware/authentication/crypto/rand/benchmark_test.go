// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Test Benchmark on My Laptop
// Result:
//
//	goos: windows
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz
//	BenchmarkFixedSize32Bytes-16    	 			 5556951	       224.1 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-256-16         	 5560414	       222.3 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-384-16         	 4748517	       257.8 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-521-16         	 4839511	       259.7 ns/op	       0 B/op	       0 allocs/op
//
// Note: These benchmarks were conducted without overclocking. If overclocked, the performance may reach goes crazy up to 1 billion operations.
//
// Additionally, the results shown here are from a semi-overclocked state (3 GHz):
//
//	goos: windows
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz
//	BenchmarkFixedSize32Bytes-16    	 6492642	       191.1 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-256-16         	 6795412	       184.9 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-384-16         	 5626008	       214.1 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-521-16         	 5705864	       211.9 ns/op	       0 B/op	       0 allocs/op
//
// On an old PC with a broken motherboard due to Windows (physical memory corruption), switched to Linux for freedom and privacy:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkFixedSize32Bytes-24    	 2482843	       505.0 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-256-24         	 2539428	       480.1 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-384-24         	 1914794	       628.2 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-521-24         	 1856566	       624.2 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-224-24           	 2471674	       501.2 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-256-24           	 2398852	       502.6 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-384-24           	 1923326	       648.6 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-521-24           	 1549710	       824.4 ns/op	       0 B/op	       0 allocs/op
//
// Note that the old PC runs on Ubuntu Pro with 16 GiB of RAM and performs smoothly without issues on Linux (it literally fixes what Windows can't), unlike Windows.
// However, there is no chance of going back to Windows due to physical memory corruption issues that also damage the memory slots, which is why
// the broken motherboard can't support more than 32 GiB++ of RAM (because Windows caused the damage).

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

// Results from the default Go benchmark test (around 10s):
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
//
// Results in a semi-overclocked state (3 GHz) for a 30s benchmark:
//
//	goos: windows
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz
//	BenchmarkFixedSizeECC/P-224-16          140615835              251.1 ns/op
//	BenchmarkFixedSizeECC/P-256-16          139366987              246.5 ns/op
//	BenchmarkFixedSizeECC/P-384-16          129598723              273.7 ns/op
//	BenchmarkFixedSizeECC/P-521-16          100000000              325.4 ns/op
//
// Note: The "semi-overclocked state" is not fully overclocked to 5 GHz because it requires detaching the laptop battery
// to use direct power and implementing additional cooling mechanisms for overclocking.
// It is also worth noting that "Intel(R) Core(TM) i9-10980HK CPU @ 2.40GHz" refers to the CPU name, not its running speed.
// Furthermore, the "Intel(R) Core(TM) i9-10980HK" is excellent for concurrency and parallelism in Go for laptop.
// When compared to other CPUs, such as those from Apple, it still performs exceptionally well (win).
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

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith10Length/Lowercase-24         	  152373	      7647 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Uppercase-24         	  152648	      7674 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Mixed-24             	  176691	      6617 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Special-24           	  162062	      7199 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/MixedSpecial-24      	  136326	      8669 ns/op	     512 B/op	      32 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	6.256s
func BenchmarkGenerateTextWith10Length(b *testing.B) {
	tests := []struct {
		length   int
		textCase rand.TextCase
	}{
		{10, rand.Lowercase},
		{10, rand.Uppercase},
		{10, rand.Mixed},
		{10, rand.Special},
		{10, rand.MixedSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith25Length/Lowercase-24         	   62334	     19020 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Uppercase-24         	   62594	     18980 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Mixed-24             	   72399	     16385 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Special-24           	   66516	     17864 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/MixedSpecial-24      	   55254	     21464 ns/op	    1264 B/op	      77 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	6.902s
func BenchmarkGenerateTextWith25Length(b *testing.B) {
	tests := []struct {
		length   int
		textCase rand.TextCase
	}{
		{25, rand.Lowercase},
		{25, rand.Uppercase},
		{25, rand.Mixed},
		{25, rand.Special},
		{25, rand.MixedSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith50Length/Lowercase-24         	   31390	     38097 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Uppercase-24         	   31441	     37985 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Mixed-24             	   36450	     32775 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Special-24           	   33487	     35665 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/MixedSpecial-24      	   27771	     43050 ns/op	    2528 B/op	     152 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	7.888s
func BenchmarkGenerateTextWith50Length(b *testing.B) {
	tests := []struct {
		length   int
		textCase rand.TextCase
	}{
		{50, rand.Lowercase},
		{50, rand.Uppercase},
		{50, rand.Mixed},
		{50, rand.Special},
		{50, rand.MixedSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith100Length/Lowercase-24         	   14958	     78638 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Uppercase-24         	   15226	     78503 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Mixed-24             	   17640	     67999 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Special-24           	   16220	     73888 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/MixedSpecial-24      	   13552	     88412 ns/op	    5024 B/op	     302 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	9.914s
func BenchmarkGenerateTextWith100Length(b *testing.B) {
	tests := []struct {
		length   int
		textCase rand.TextCase
	}{
		{100, rand.Lowercase},
		{100, rand.Uppercase},
		{100, rand.Mixed},
		{100, rand.Special},
		{100, rand.MixedSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith500Length/Lowercase-24         	    3056	    378141 ns/op	   25027 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Uppercase-24         	    3061	    378357 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Mixed-24             	    3512	    326301 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Special-24           	    3286	    355469 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/MixedSpecial-24      	    2696	    428179 ns/op	   25024 B/op	    1502 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	5.987s
func BenchmarkGenerateTextWith500Length(b *testing.B) {
	tests := []struct {
		length   int
		textCase rand.TextCase
	}{
		{500, rand.Lowercase},
		{500, rand.Uppercase},
		{500, rand.Mixed},
		{500, rand.Special},
		{500, rand.MixedSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
