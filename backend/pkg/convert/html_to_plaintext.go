// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"strings"

	"golang.org/x/net/html"
)

// HTMLToPlainText converts HTML content to plain text.
// It parses the HTML and extracts text nodes, concatenating them into a single string.
// If parsing fails, it returns the original HTML content as a fallback.
//
// Note: This function does not fully handle elements like "<style>" or other non-text content.
//
// TODO: Improving this will require additional filtering, possibly using regex.
func HTMLToPlainText(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent // Fallback to original HTML if parsing fails
	}

	var textContent strings.Builder
	var extractText func(*html.Node)

	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			textContent.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)
	return textContent.String()
}
