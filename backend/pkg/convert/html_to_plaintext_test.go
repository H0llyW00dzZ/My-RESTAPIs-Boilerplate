// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert_test

import (
	"bytes"
	"h0llyw00dz-template/backend/pkg/convert"
	"io"
	"runtime"
	"strings"
	"testing"
)

const (
	simpleInput = `<div><h1>Hello</h1> <span>HTML</span> <p>Frontend,</p> <strong>from Go</strong></div>`
	largeInput  = `
<!DOCTYPE html>
<html>
<head>
    <title>Go Programming Language</title>
    <style>
        .content { font-family: Arial; }
    </style>
</head>
<body>
    <div class="content">
        <h1>Why Go is Great for Systems Programming ? 🤔</h1>
		<img src="https://go.dev/images/gophers/biplane.svg" alt="Gopher Biplane Ready To Fly">
        <p>Go, also known as Golang, is designed for simplicity and efficiency.</p>
        <p>Here are some reasons why Go excels:</p>
        <ul>
            <li>Concurrency support with goroutines</li>
            <li>Fast compilation times</li>
            <li>Robust standard library</li>
        </ul>
        <p>Discover more about Go at the <a href="https://go.dev">official site</a>.</p>
    </div>
</body>
</html>`
)

func TestHTMLToPlainText(t *testing.T) {
	// Determine the newline character based on the operating system.
	crlf := getNewline()

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
			input:    "<p>Hello,\r\nWorld!</p>",
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
			expected: crlf + crlf + crlf + "Hello" + crlf + "HTML" + crlf + crlf + "Frontend," + crlf + crlf + "from Go" + crlf + crlf,
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
			expected: crlf + crlf + "Hello HTML Frontend, from Go." + crlf + crlf,
		},
		{
			name: "Style with Class",
			input: `<style class="example">
						.example { color: red; }
					</style>
					<p>Hello HTML Frontend, from Go.</p>`,
			expected: crlf + crlf + "Hello HTML Frontend, from Go." + crlf + crlf,
		},
		{
			name:     "Another Links",
			input:    "Visit <a href=\"https://go.dev/dl/\">Go Dev</a> to download Go.",
			expected: "Visit [Go Dev](https://go.dev/dl/) to download Go.",
		},
		{
			name:     "Large Input",
			input:    largeInput,
			expected: "Go Programming Language" + crlf + crlf + crlf + "Why Go is Great for Systems Programming ? 🤔" + crlf + "![Gopher Biplane Ready To Fly](https://go.dev/images/gophers/biplane.svg)" + crlf + crlf + "Go, also known as Golang, is designed for simplicity and efficiency." + crlf + crlf + crlf + crlf + "Here are some reasons why Go excels:" + crlf + crlf + crlf + "- Concurrency support with goroutines" + crlf + "- Fast compilation times" + crlf + "- Robust standard library" + crlf + crlf + crlf + crlf + "Discover more about Go at the [official site](https://go.dev) ." + crlf + crlf + crlf + crlf,
		},
		{
			name: "Table Structure",
			input: `<table>
						<tr><th>Header 1</th><th>Header 2</th></tr>
						<tr><td>Row 1 Col 1</td><td>Row 1 Col 2</td></tr>
						<tr><td>Row 2 Col 1</td><td>Row 2 Col 2</td></tr>
					</table>`,
			expected: crlf + crlf + " | Header 1 | Header 2" + crlf + " | Row 1 Col 1 | Row 1 Col 2" + crlf + " | Row 2 Col 1 | Row 2 Col 2" + crlf + crlf,
		},
		{
			name:     "Noscript Element",
			input:    `<noscript>This is a noscript content.</noscript><p>Visible content.</p>`,
			expected: crlf + crlf + "Visible content." + crlf + crlf,
		},
		{
			name: "Script Element",
			input: `<script>
						console.log('This is a script.');
					</script>
					<p>Visible content.</p>`,
			expected: crlf + crlf + "Visible content." + crlf + crlf,
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
	crlf := getNewline()

	// Expected plain text output.
	expected := "Go Programming Language" + crlf + crlf + crlf + "Why Go is Great for Systems Programming ? 🤔" + crlf + "![Gopher Biplane Ready To Fly](https://go.dev/images/gophers/biplane.svg)" + crlf + crlf + "Go, also known as Golang, is designed for simplicity and efficiency." + crlf + crlf + crlf + crlf + "Here are some reasons why Go excels:" + crlf + crlf + crlf + "- Concurrency support with goroutines" + crlf + "- Fast compilation times" + crlf + "- Robust standard library" + crlf + crlf + crlf + crlf + "Discover more about Go at the [official site](https://go.dev) ." + crlf + crlf + crlf + crlf

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

// getNewline returns the appropriate newline characters based on the operating system.
//
// Note: Currently supports only Linux/Unix and Windows (MS-DOS). Other OS support is marked as TODO.
func getNewline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func TestHTMLToPlainTextConcurrent(t *testing.T) {
	// Determine the newline character based on the operating system.
	crlf := getNewline()

	htmlContents := []string{
		"<p>Hello, World!</p>",
		"<h1>Title</h1><p>Paragraph</p>",
		"<div><a href=\"https://example.com\">Link</a></div>",
	}

	expectedResults := []string{
		crlf + crlf + "Hello, World!" + crlf + crlf,
		crlf + "Title" + crlf + crlf + crlf + "Paragraph" + crlf + crlf,
		crlf + crlf + "[Link](https://example.com)" + crlf + crlf,
	}

	results := convert.HTMLToPlainTextConcurrent(htmlContents)

	for i, result := range results {
		t.Log("Expected:", expectedResults[i])
		t.Log("Result:", result)
		if result != expectedResults[i] {
			t.Errorf("Test %d failed: expected %q, got %q", i, expectedResults[i], result)
		}
	}
}

func TestHTMLToPlainTextStreamsConcurrent(t *testing.T) {
	crlf := getNewline()

	htmlInputs := []string{
		"<p>Hello, World!</p>",
		"<h1>Title</h1><p>Paragraph</p>",
		"<div><a href=\"https://example.com\">Link</a></div>",
	}

	expectedOutputs := []string{
		crlf + crlf + "Hello, World!" + crlf + crlf,
		crlf + "Title" + crlf + crlf + crlf + "Paragraph" + crlf + crlf,
		crlf + crlf + "[Link](https://example.com)" + crlf + crlf,
	}

	var readers []io.Reader
	for _, input := range htmlInputs {
		readers = append(readers, strings.NewReader(input))
	}

	var output bytes.Buffer
	errs := convert.HTMLToPlainTextStreamsConcurrent(readers, &output)

	if len(errs) > 0 {
		t.Fatalf("Encountered errors: %v", errs)
	}

	result := output.String()
	for i, expected := range expectedOutputs {
		t.Logf("Test %d - Expected Output: %q", i, expected)
		t.Logf("Test %d - Actual Result: %q", i, result)
		if !strings.Contains(result, expected) {
			t.Errorf("Test %d failed: expected to find %q in result, but it was missing", i, expected)
		} else {
			t.Logf("Test %d passed: found expected output.", i)
		}
	}
}
