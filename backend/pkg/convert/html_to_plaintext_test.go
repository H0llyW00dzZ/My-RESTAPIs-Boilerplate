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
	// Note: This depends on the OS. Linux/Unix uses "\n", while Windows uses "\r\n" due to MS-DOS conventions.
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
			expected: "Hello, World!",
		},
		{
			name:     "HTML with Newlines",
			input:    "<p>Hello,\nWorld!</p>",
			expected: "Hello,\nWorld!",
		},
		{
			name:     "HTML with CRLF",
			input:    "<p>Hello,\r\nWorld!</p>",
			expected: strings.ReplaceAll("Hello,\nWorld!", "\n", crlf),
		},
		{
			name:     "HTML with Nested Elements",
			input:    "<div>Hello,<br>World!</div>",
			expected: "Hello,World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convert.HTMLToPlainText(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
