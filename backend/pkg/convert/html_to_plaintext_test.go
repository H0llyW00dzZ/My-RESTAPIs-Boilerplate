// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert_test

import (
	"h0llyw00dz-template/backend/pkg/convert"
	"runtime"
	"strings"
	"testing"
)

func TestHTMLToPlainText(t *testing.T) {
	// Note: This test might depend on the OS. Linux/Unix uses "\n", while Windows uses "\r\n" due to MS-DOS conventions.
	// When testing on Windows, "\r" might be required.
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
			input:    "<p>Hello,\nWorld!</p>",
			expected: crlf + crlf + "Hello,\nWorld!" + crlf + crlf,
		},
		{
			name:     "HTML with CRLF",
			input:    "<p>Hello,\r\nWorld!</p>",
			expected: strings.ReplaceAll(crlf+crlf+"Hello,\nWorld!"+crlf+crlf, "\n", crlf),
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
