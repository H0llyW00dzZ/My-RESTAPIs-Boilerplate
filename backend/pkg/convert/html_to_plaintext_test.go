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
			expected: "\n\nHello, World!\n\n",
		},
		{
			name:     "HTML with Newlines",
			input:    "<p>Hello,\nWorld!</p>",
			expected: "\n\nHello,\nWorld!\n\n",
		},
		{
			name:     "HTML with CRLF",
			input:    "<p>Hello,\r\nWorld!</p>",
			expected: strings.ReplaceAll("\n\nHello,\nWorld!\n\n", "\n", crlf),
		},
		{
			name:     "HTML with Nested Elements",
			input:    "<div>Hello,<br>World!</div>",
			expected: "\n\nHello,\nWorld!\n\n",
		},
		{
			name:     "Complex HTML Structure",
			input:    "<div><h1>Hello</h1> <span>HTML</span> <p>Frontend,</p> <strong>from Go</strong></div>",
			expected: "\n\n\nHello\n HTML \n\nFrontend,\n\n from Go\n\n",
		},
		{
			name:     "List Items",
			input:    "<ul><li>Get Good</li><li>Get Go</li></ul>",
			expected: "\n- Get Good\n- Get Go\n\n",
		},
		{
			name:     "Ordered List",
			input:    "<ol><li>First</li><li>Second</li></ol>",
			expected: "\n- First\n- Second\n\n",
		},
		{
			name:     "Headings",
			input:    "<h1>Hello HTML Frontend</h1><h2>from</h2><p>Go</p>",
			expected: "\nHello HTML Frontend\n\nfrom\n\n\nGo\n\n",
		},
		{
			name:     "Links",
			input:    "<p>Visit <a href=\"https://go.dev/dl/\">Go Dev</a> to download Go.</p>",
			expected: "\n\nVisit [Go Dev](https://go.dev/dl/) to download Go.\n\n",
		},
		{
			name:     "Multiple Paragraphs",
			input:    "<p>Hello HTML Frontend, from Go.</p><p>Hello HTML Frontend, from Go.</p>",
			expected: "\n\nHello HTML Frontend, from Go.\n\n\n\nHello HTML Frontend, from Go.\n\n",
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
