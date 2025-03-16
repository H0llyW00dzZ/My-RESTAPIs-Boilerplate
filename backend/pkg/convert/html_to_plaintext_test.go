// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert_test

import (
	"bytes"
	"h0llyw00dz-template/backend/pkg/convert"
	"runtime"
	"strings"
	"testing"
)

func TestHTMLToPlainText(t *testing.T) {
	// Determine the newline character based on the operating system.
	crlf := "\n"
	if runtime.GOOS == "windows" {
		crlf = "\r\n"
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple HTML",
			input:    "<p>Hello, World!</p>",
			expected: crlf + crlf + "Hello, World!" + crlf + crlf,
		},
		{
			name:     "HTML with Newlines",
			input:    "<p>Hello," + crlf + "World!</p>",
			expected: crlf + crlf + "Hello," + crlf + "World!" + crlf + crlf,
		},
		{
			name:     "HTML with CRLF",
			input:    "<p>Hello,\r" + crlf + "World!</p>",
			expected: crlf + crlf + "Hello," + crlf + "World!" + crlf + crlf,
		},
		{
			name:     "HTML with Nested Elements",
			input:    "<div>Hello,<br>World!</div>",
			expected: crlf + crlf + "Hello," + crlf + "World!" + crlf + crlf,
		},
		{
			name:     "Complex HTML Structure",
			input:    "<div><h1>Hello</h1> <span>HTML</span> <p>Frontend,</p> <strong>from Go</strong></div>",
			expected: crlf + crlf + crlf + "Hello" + crlf + " HTML " + crlf + crlf + "Frontend," + crlf + crlf + " from Go" + crlf + crlf,
		},
		{
			name:     "List Items",
			input:    "<ul><li>Get Good</li><li>Get Go</li></ul>",
			expected: crlf + "- Get Good" + crlf + "- Get Go" + crlf + crlf,
		},
		{
			name:     "Ordered List",
			input:    "<ol><li>First</li><li>Second</li></ol>",
			expected: crlf + "- First" + crlf + "- Second" + crlf + crlf,
		},
		{
			name:     "Headings",
			input:    "<h1>Hello HTML Frontend</h1><h2>from</h2><p>Go</p>",
			expected: crlf + "Hello HTML Frontend" + crlf + crlf + "from" + crlf + crlf + crlf + "Go" + crlf + crlf,
		},
		{
			name:     "Links",
			input:    "<p>Visit <a href=\"https://go.dev/dl/\">Go Dev</a> to download Go.</p>",
			expected: crlf + crlf + "Visit [Go Dev](https://go.dev/dl/) to download Go." + crlf + crlf,
		},
		{
			name:     "Multiple Paragraphs",
			input:    "<p>Hello HTML Frontend, from Go.</p><p>Hello HTML Frontend, from Go.</p>",
			expected: crlf + crlf + "Hello HTML Frontend, from Go." + crlf + crlf + crlf + crlf + "Hello HTML Frontend, from Go." + crlf + crlf,
		},
		{
			name: "Complex with Style",
			input: `<style>
						body { font-family: Arial; }
					</style>
					<p>Hello HTML Frontend, from Go.</p>`,
			expected: crlf + "\t\t\t\t\t" + crlf + crlf + "Hello HTML Frontend, from Go." + crlf + crlf,
		},
		{
			name: "Style with Class",
			input: `<style class="example">
						.example { color: red; }
					</style>
					<p>Hello HTML Frontend, from Go.</p>`,
			expected: crlf + "\t\t\t\t\t" + crlf + crlf + "Hello HTML Frontend, from Go." + crlf + crlf,
		},
		{
			name:     "Another Links",
			input:    "Visit <a href=\"https://go.dev/dl/\">Go Dev</a> to download Go.",
			expected: "Visit [Go Dev](https://go.dev/dl/) to download Go.",
		},
		{
			name:     "Large Input",
			input:    largeInput,
			expected: crlf + "    Go Programming Language" + crlf + "    " + crlf + "" + crlf + "" + crlf + "    " + crlf + "" + crlf + "" + crlf + "        " + crlf + "Why Go is Great for Systems Programming" + crlf + "" + crlf + "        " + crlf + "" + crlf + "Go, also known as Golang, is designed for simplicity and efficiency." + crlf + "" + crlf + "" + crlf + "        " + crlf + "" + crlf + "Here are some reasons why Go excels:" + crlf + "" + crlf + "" + crlf + "        " + crlf + "" + crlf + "            - Concurrency support with goroutines" + crlf + "" + crlf + "            - Fast compilation times" + crlf + "" + crlf + "            - Robust standard library" + crlf + "" + crlf + "        " + crlf + "" + crlf + "        " + crlf + "" + crlf + "Discover more about Go at the [official site](https://go.dev)." + crlf + "" + crlf + "" + crlf + "    " + crlf + "" + crlf + "" + crlf + "" + crlf + "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convert.HTMLToPlainText(tt.input)
			t.Log("Expected:", tt.expected)
			t.Log("Result:", result)
			if result != tt.expected {
				t.Errorf("expected: %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestHTMLToPlainTextStreams_LargeInput(t *testing.T) {
	// Determine the newline character based on the operating system.
	crlf := "\n"
	if runtime.GOOS == "windows" {
		crlf = "\r\n"
	}

	// Expected plain text output.
	expected := crlf + "    Go Programming Language" + crlf + "    " + crlf + "" + crlf + "" + crlf + "    " + crlf + "" + crlf + "" + crlf + "        " + crlf + "Why Go is Great for Systems Programming" + crlf + "" + crlf + "        " + crlf + "" + crlf + "Go, also known as Golang, is designed for simplicity and efficiency." + crlf + "" + crlf + "" + crlf + "        " + crlf + "" + crlf + "Here are some reasons why Go excels:" + crlf + "" + crlf + "" + crlf + "        " + crlf + "" + crlf + "            - Concurrency support with goroutines" + crlf + "" + crlf + "            - Fast compilation times" + crlf + "" + crlf + "            - Robust standard library" + crlf + "" + crlf + "        " + crlf + "" + crlf + "        " + crlf + "" + crlf + "Discover more about Go at the [official site](https://go.dev)." + crlf + "" + crlf + "" + crlf + "    " + crlf + "" + crlf + "" + crlf + "" + crlf + ""

	// Create a reader and writer for the test.
	input := strings.NewReader(largeInput)
	var output bytes.Buffer

	// Run the conversion.
	err := convert.HTMLToPlainTextStreams(input, &output)
	if err != nil {
		t.Fatalf("Failed to convert HTML to plain text: %v", err)
	}

	// Get the result and compare it to the expected output.
	result := output.String()
	t.Log("Expected:", expected)
	t.Log("Result:", result)
	if result != expected {
		t.Errorf("expected: %q, got %q", expected, result)
	}
}
