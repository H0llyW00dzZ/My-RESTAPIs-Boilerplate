// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert_test

import (
	"bytes"
	"h0llyw00dz-template/backend/pkg/convert"
	"io"
	"strings"
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

// Results on a broken PC:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/pkg/convert
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkHTMLToPlainTextStreams_Hybrid/Simple_HTML-24         	  843254	      1413 ns/op	    5068 B/op	      10 allocs/op
//	BenchmarkHTMLToPlainTextStreams_Hybrid/Large_HTML-24          	  853084	      1374 ns/op	    5068 B/op	      10 allocs/op
//	PASS
//	ok  	h0llyw00dz-template/backend/pkg/convert	2.369s
//
// Note: This is much better because it uses streaming.
func BenchmarkHTMLToPlainTextStreams_Hybrid(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{"Simple HTML", simpleInput},
		{"Large HTML", largeInput},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			input := strings.NewReader(bm.input)
			output := &bytes.Buffer{}
			for b.Loop() {
				if err := convert.HTMLToPlainTextStreams(input, output); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Results on a broken PC:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/pkg/convert
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkHTMLToPlainTextConcurrent/Simple_HTML-24         	   66913	     17689 ns/op	    9176 B/op	      83 allocs/op
//	BenchmarkHTMLToPlainTextConcurrent/Large_HTML-24          	   36975	     32422 ns/op	   14995 B/op	     151 allocs/op
//	BenchmarkHTMLToPlainTextConcurrent/Mixed_HTML-24          	   33387	     35895 ns/op	   21653 B/op	     184 allocs/op
//
//	PASS
//	ok  	h0llyw00dz-template/backend/pkg/convert	3.559s
func BenchmarkHTMLToPlainTextConcurrent(b *testing.B) {
	benchmarks := []struct {
		name   string
		inputs []string
	}{
		{"Simple HTML", []string{simpleInput}},
		{"Large HTML", []string{largeInput}},
		{"Mixed HTML", []string{simpleInput, largeInput}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for b.Loop() {
				convert.HTMLToPlainTextConcurrent(bm.inputs)
			}
		})
	}
}

// Results on a broken PC:
//
//	goos: linux
//	goarch: amd64
//	pkg: h0llyw00dz-template/backend/pkg/convert
//	cpu: AMD Ryzen 9 3900X 12-Core Processor
//	BenchmarkHTMLToPlainTextStreamsConcurrent/Simple_HTML-24         	  376934	      2942 ns/op	    5221 B/op	      16 allocs/op
//	BenchmarkHTMLToPlainTextStreamsConcurrent/Large_HTML-24          	  424011	      2943 ns/op	    5221 B/op	      16 allocs/op
//	BenchmarkHTMLToPlainTextStreamsConcurrent/Mixed_HTML-24          	  196773	      6220 ns/op	   10372 B/op	      28 allocs/op
//
//	PASS
//	ok  	h0llyw00dz-template/backend/pkg/convert	3.586s
func BenchmarkHTMLToPlainTextStreamsConcurrent(b *testing.B) {
	benchmarks := []struct {
		name   string
		inputs []string
	}{
		{"Simple HTML", []string{simpleInput}},
		{"Large HTML", []string{largeInput}},
		{"Mixed HTML", []string{simpleInput, largeInput}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			readers := make([]io.Reader, len(bm.inputs))
			for i, input := range bm.inputs {
				readers[i] = strings.NewReader(input)
			}
			output := &bytes.Buffer{}
			for b.Loop() {
				if errs := convert.HTMLToPlainTextStreamsConcurrent(readers, output); len(errs) > 0 {
					b.Fatal(errs)
				}
			}
		})
	}
}
