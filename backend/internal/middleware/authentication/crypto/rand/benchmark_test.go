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
//	BenchmarkFixedSize32Bytes-16		    	 6492642	       191.1 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-256-16         	 6795412	       184.9 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-384-16         	 5626008	       214.1 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-521-16         	 5705864	       211.9 ns/op	       0 B/op	       0 allocs/op
//
// On an old PC with a broken motherboard due to Windows (physical memory corruption), switched to Linux (Kernel 6.11) for freedom and privacy:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkFixedSize32Bytes-24    			11736326	       101.4 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-256-24         	11798857	       101.7 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-384-24         	 8896884	       134.1 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECDSA/P-521-24         	 8884756	       134.3 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-224-24         		12967860	        92.56 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-256-24         		11932174	       100.4 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-384-24	         	 9008350	       132.9 ns/op	       0 B/op	       0 allocs/op
//	BenchmarkFixedSizeECC/P-521-24	         	 6845301	       173.9 ns/op	       0 B/op	       0 allocs/op
//
// Note that the old PC runs on Ubuntu Pro with 16 GiB of RAM and performs smoothly without issues on Linux (Kernel 6.11) (it literally fixes what Windows can't), unlike Windows.
// However, there is no chance of going back to Windows due to physical memory corruption issues that also damage the memory slots, which is why
// the broken motherboard can't support more than 32 GiB++ of RAM (because Windows caused the damage).
//
// Additionally, these benchmarks have been updated. The power of AMD has returned, and it remains powerful even on a broken PC. So, get good, get AMD.

package rand_test

import (
	"crypto/elliptic"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand"
	"testing"
)

func BenchmarkFixedSize32Bytes(b *testing.B) {
	reader := rand.FixedSize32Bytes()
	buf := make([]byte, 32)

	for b.Loop() {
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

			for b.Loop() {
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

			for b.Loop() {
				if _, err := reader.Read(buf); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux (Kernel 6.11) with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith10Length/Lowercase-24         	  641248	      1639 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Uppercase-24         	  726758	      1634 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Mixed-24             	  778693	      1518 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Special-24           	  709845	      1592 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/MixedSpecial-24      	  620727	      1747 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/Number-24            	  604453	      1859 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/NumSpecial-24        	  596522	      1939 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/LowercaseSpecial-24  	  719533	      1637 ns/op	     512 B/op	      32 allocs/op
//	BenchmarkGenerateTextWith10Length/UppercaseSpecial-24  	  703129	      1637 ns/op	     512 B/op	      32 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	7.550s
//
// Note: These benchmarks are updated. The power of AMD has returned, and it remains powerful even on a broken PC. So, get good, get AMD.
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
		{10, rand.NumSpecial},
		{10, rand.LowercaseSpecial},
		{10, rand.UppercaseSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for b.Loop() {
				if _, err := rand.GenerateText(tt.length, tt.textCase); err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux (Kernel 6.11) with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith25Length/Lowercase-24         	  275292	      4071 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Uppercase-24         	  297151	      4043 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Mixed-24             	  319545	      3739 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Special-24           	  302397	      3921 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/MixedSpecial-24      	  272220	      4332 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/Number-24            	  256512	      4630 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/NumSpecial-24        	  238528	      4987 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/LowercaseSpecial-24  	  297484	      3985 ns/op	    1264 B/op	      77 allocs/op
//	BenchmarkGenerateTextWith25Length/UppercaseSpecial-24  	  295318	      3987 ns/op	    1264 B/op	      77 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	8.347s
//
// Note: These benchmarks are updated. The power of AMD has returned, and it remains powerful even on a broken PC. So, get good, get AMD.
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
		{25, rand.NumSpecial},
		{25, rand.LowercaseSpecial},
		{25, rand.UppercaseSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for b.Loop() {
				if _, err := rand.GenerateText(tt.length, tt.textCase); err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux (Kernel 6.11) with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith50Length/Lowercase-24         	  145754	      8023 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Uppercase-24         	  147306	      8034 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Mixed-24             	  159758	      7436 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Special-24           	  152212	      7763 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/MixedSpecial-24      	  139192	      8579 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/Number-24            	  131802	      9084 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/NumSpecial-24        	  125325	      9435 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/LowercaseSpecial-24  	  149457	      7966 ns/op	    2528 B/op	     152 allocs/op
//	BenchmarkGenerateTextWith50Length/UppercaseSpecial-24  	  149586	      7952 ns/op	    2528 B/op	     152 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	9.586s
//
// Note: These benchmarks are updated. The power of AMD has returned, and it remains powerful even on a broken PC. So, get good, get AMD.
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
		{50, rand.NumSpecial},
		{50, rand.LowercaseSpecial},
		{50, rand.UppercaseSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for b.Loop() {
				if _, err := rand.GenerateText(tt.length, tt.textCase); err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux (Kernel 6.11) with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith100Length/Lowercase-24         	   73147	     16147 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Uppercase-24         	   74023	     16136 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Mixed-24             	   80100	     14843 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Special-24           	   76915	     15583 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/MixedSpecial-24      	   68668	     17290 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/Number-24            	   64477	     18476 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/NumSpecial-24        	   62823	     18901 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/LowercaseSpecial-24  	   74540	     15975 ns/op	    5024 B/op	     302 allocs/op
//	BenchmarkGenerateTextWith100Length/UppercaseSpecial-24  	   75105	     15830 ns/op	    5024 B/op	     302 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	11.974s
//
// Note: These benchmarks are updated. The power of AMD has returned, and it remains powerful even on a broken PC. So, get good, get AMD.
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
		{100, rand.NumSpecial},
		{100, rand.LowercaseSpecial},
		{100, rand.UppercaseSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for b.Loop() {
				if _, err := rand.GenerateText(tt.length, tt.textCase); err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on an old PC with a broken motherboard. It still works fine & stable on Linux (Kernel 6.11) with only 16 GiB of RAM:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateTextWith500Length/Lowercase-24         	   15006	     79826 ns/op	   25025 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Uppercase-24         	   15014	     79825 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Mixed-24             	   16274	     73732 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Special-24           	   15488	     77354 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/MixedSpecial-24      	   14090	     85175 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/Number-24            	   13257	     90534 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/NumSpecial-24        	   12939	     92687 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/LowercaseSpecial-24  	   15104	     79776 ns/op	   25024 B/op	    1502 allocs/op
//	BenchmarkGenerateTextWith500Length/UppercaseSpecial-24  	   15072	     79801 ns/op	   25024 B/op	    1502 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	7.232s
//
// Note: These benchmarks are updated. The power of AMD has returned, and it remains powerful even on a broken PC. So, get good, get AMD.
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
		{500, rand.NumSpecial},
		{500, rand.LowercaseSpecial},
		{500, rand.UppercaseSpecial},
	}

	for _, tt := range tests {
		b.Run(textCaseToString(tt.textCase), func(b *testing.B) {
			for b.Loop() {
				if _, err := rand.GenerateText(tt.length, tt.textCase); err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Results on a broken PC:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkGenerateFixedUUID-24    	 2163597	       555.5 ns/op	     184 B/op	       7 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/internal/middleware/authentication/crypto/rand	1.207s
func BenchmarkGenerateFixedUUID(b *testing.B) {
	for b.Loop() {
		if _, err := rand.GenerateFixedUUID(); err != nil {
			b.Fatal(err)
		}
	}
}
