// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert_test

import (
	"h0llyw00dz-template/backend/pkg/convert"
	"testing"
)

// Results on a broken PC:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/pkg/convert
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkHTMLToPlainText/Simple_HTML-24         	  263047	      4461 ns/op	    6552 B/op	      32 allocs/op
//	BenchmarkHTMLToPlainText/Large_HTML-24          	   85545	     13971 ns/op	   11395 B/op	      94 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/pkg/convert	2.336s
func BenchmarkHTMLToPlainText(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{"Simple HTML", simpleInput},
		{"Large HTML", largeInput},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for b.Loop() {
				convert.HTMLToPlainText(bm.input)
			}
		})
	}
}
