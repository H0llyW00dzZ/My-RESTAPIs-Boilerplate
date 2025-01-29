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
//	BenchmarkGenerateTextWith10Length/Lowercase-24         	  152569	      7684 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Uppercase-24         	  153183	      7663 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Mixed-24             	  176103	      6646 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Special-24           	  163130	      7227 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/MixedSpecial-24      	  135039	      8666 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Number-24            	  122524	      9646 ns/op	     512 B/op	      32 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	7.550s
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
		{10, rand.Number},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
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
//	BenchmarkGenerateTextWith25Length/Lowercase-24         	   62409	     19039 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Uppercase-24         	   62414	     19094 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Mixed-24             	   71863	     16410 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Special-24           	   66769	     17901 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/MixedSpecial-24      	   54975	     21632 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Number-24            	   49626	     23905 ns/op	    1264 B/op	      77 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	8.347s
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
		{25, rand.Number},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
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
//	BenchmarkGenerateTextWith50Length/Lowercase-24         	   31603	     38002 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Uppercase-24         	   31491	     38100 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Mixed-24             	   36438	     32906 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Special-24           	   33471	     35800 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/MixedSpecial-24      	   27567	     42972 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Number-24            	   25027	     48102 ns/op	    2528 B/op	     152 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	9.586s
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
		{50, rand.Number},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
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
//	BenchmarkGenerateTextWith100Length/Lowercase-24         	   15736	     76019 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Uppercase-24         	   15748	     76029 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Mixed-24             	   18262	     65867 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Special-24           	   16801	     71570 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/MixedSpecial-24      	   13963	     86205 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Number-24            	   12474	     95632 ns/op	    5024 B/op	     302 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	11.974s
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
		{100, rand.Number},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
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
//	BenchmarkGenerateTextWith500Length/Lowercase-24         	    3074	    380179 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Uppercase-24         	    3069	    381277 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Mixed-24             	    3506	    327678 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Special-24           	    3254	    356761 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/MixedSpecial-24      	    2700	    428901 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Number-24            	    2452	    476670 ns/op	   25024 B/op	    1502 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	7.232s
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
		{500, rand.Number},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rand.GenerateText(tt.length, tt.textCase)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
