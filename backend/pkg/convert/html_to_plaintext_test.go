// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert_test

import (
	"bytes"
	"h0llyw00dz-template/backend/pkg/convert"
	"io"
	"os"
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

	tableInput = `
<!DOCTYPE html>
<html>
<head>
    <style>
        table, th, td {
            border: 1px solid black;
            border-collapse: collapse;
            padding: 8px;
            text-align: left;
        }
    </style>
</head>
<body>

<h2>Introduction to Go Programming Language</h2>

<table style="width:100%">
    <tr>
        <th>Feature</th>
        <th>Description</th>
    </tr>
    <tr>
        <td>Concurrency</td>
        <td>Go provides built-in support for concurrent programming with goroutines and channels.</td>
    </tr>
    <tr>
        <td>Static Typing</td>
        <td>Go is statically typed, which helps catch errors at compile time.</td>
    </tr>
    <tr>
        <td>Efficiency</td>
        <td>Go compiles quickly and produces fast executables, making it ideal for performance-critical applications.</td>
    </tr>
</table>

<p>Go, also known as Golang, is a statically typed, compiled programming language designed for simplicity and efficiency.</p>

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
			name:     "Table Structure",
			input:    tableInput,
			expected: crlf + "Introduction to Go Programming Language" + crlf + crlf + "|  Feature | Description |" + crlf + "|---------|-------------|" + crlf + "|  Concurrency | Go provides built-in support for concurrent programming with goroutines and channels. |" + crlf + "|  Static Typing | Go is statically typed, which helps catch errors at compile time. |" + crlf + "|  Efficiency | Go compiles quickly and produces fast executables, making it ideal for performance-critical applications. |" + crlf + crlf + crlf + crlf + crlf + "Go, also known as Golang, is a statically typed, compiled programming language designed for simplicity and efficiency." + crlf + crlf,
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

// Note: This test is a better way
func TestHTMLToPlainTextStreams_LargeInput(t *testing.T) {
	// Determine the newline character based on the operating system.
	crlf := getNewline()

	tests := []struct {
		name     string
		input    func() io.Reader
		output   func() io.Writer
		expected string
	}{
		{
			name:     "Simple HTML",
			input:    func() io.Reader { return strings.NewReader("<p>Hello, World!</p>") },
			output:   func() io.Writer { return &bytes.Buffer{} },
			expected: crlf + crlf + "Hello, World!" + crlf + crlf,
		},
		{
			name:     "Large HTML to Buffer",
			input:    func() io.Reader { return strings.NewReader(largeInput) },
			output:   func() io.Writer { return &bytes.Buffer{} },
			expected: "Go Programming Language" + crlf + crlf + crlf + "Why Go is Great for Systems Programming ? 🤔" + crlf + "![Gopher Biplane Ready To Fly](https://go.dev/images/gophers/biplane.svg)" + crlf + crlf + "Go, also known as Golang, is designed for simplicity and efficiency." + crlf + crlf + crlf + crlf + "Here are some reasons why Go excels:" + crlf + crlf + crlf + "- Concurrency support with goroutines" + crlf + "- Fast compilation times" + crlf + "- Robust standard library" + crlf + crlf + crlf + crlf + "Discover more about Go at the [official site](https://go.dev) ." + crlf + crlf + crlf + crlf,
		},
		{
			name:  "Large HTML to File",
			input: func() io.Reader { return strings.NewReader(largeInput) },
			output: func() io.Writer {
				file, err := os.CreateTemp("", "output.txt")
				if err != nil {
					t.Fatalf("Failed to create temporary file: %v", err)
				}
				t.Cleanup(func() { os.Remove(file.Name()) })
				return file
			},
			expected: "Go Programming Language" + crlf + crlf + crlf + "Why Go is Great for Systems Programming ? 🤔" + crlf + "![Gopher Biplane Ready To Fly](https://go.dev/images/gophers/biplane.svg)" + crlf + crlf + "Go, also known as Golang, is designed for simplicity and efficiency." + crlf + crlf + crlf + crlf + "Here are some reasons why Go excels:" + crlf + crlf + crlf + "- Concurrency support with goroutines" + crlf + "- Fast compilation times" + crlf + "- Robust standard library" + crlf + crlf + crlf + crlf + "Discover more about Go at the [official site](https://go.dev) ." + crlf + crlf + crlf + crlf,
		},
		{
			name: "File Input to Buffer",
			input: func() io.Reader {
				file, err := os.CreateTemp("", "input.html")
				if err != nil {
					t.Fatalf("Failed to create temporary input file: %v", err)
				}
				t.Cleanup(func() { os.Remove(file.Name()) })
				_, err = file.WriteString(largeInput)
				if err != nil {
					t.Fatalf("Failed to write to temporary input file: %v", err)
				}
				file.Seek(0, io.SeekStart)
				return file
			},
			output:   func() io.Writer { return &bytes.Buffer{} },
			expected: "Go Programming Language" + crlf + crlf + crlf + "Why Go is Great for Systems Programming ? 🤔" + crlf + "![Gopher Biplane Ready To Fly](https://go.dev/images/gophers/biplane.svg)" + crlf + crlf + "Go, also known as Golang, is designed for simplicity and efficiency." + crlf + crlf + crlf + crlf + "Here are some reasons why Go excels:" + crlf + crlf + crlf + "- Concurrency support with goroutines" + crlf + "- Fast compilation times" + crlf + "- Robust standard library" + crlf + crlf + crlf + crlf + "Discover more about Go at the [official site](https://go.dev) ." + crlf + crlf + crlf + crlf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input()
			output := tt.output()

			// Run the conversion.
			err := convert.HTMLToPlainTextStreams(input, output)
			if err != nil {
				t.Fatalf("Failed to convert HTML to plain text: %v", err)
			}

			// Read the result from the output.
			var result string
			switch out := output.(type) {
			case *bytes.Buffer:
				result = out.String()
			case *os.File:
				out.Seek(0, io.SeekStart)
				content, err := io.ReadAll(out)
				if err != nil {
					t.Fatalf("Failed to read from file: %v", err)
				}
				result = string(content)
			}

			// Compare the result to the expected output.
			t.Log("Expected:", tt.expected)
			t.Log("Result:", result)
			if result != tt.expected {
				t.Errorf("expected: %q, got %q", tt.expected, result)
			}
		})
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
