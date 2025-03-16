// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"fmt"
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
		} else if n.Type == html.ElementNode {
			switch n.Data {
			case "br":
				textContent.WriteString("\n")
			case "p", "div":
				textContent.WriteString("\n\n")
			case "h1", "h2", "h3", "h4", "h5", "h6":
				textContent.WriteString("\n")
			case "ul", "ol":
				textContent.WriteString("\n")
			case "li":
				textContent.WriteString("- ")
			case "a":
				href := getAttrValue(n, "href")

				// Note: This is what it will look like in markdown format "Visit [Example](https://example.com) website."
				textContent.WriteString(fmt.Sprintf("[%s](%s)", getTextContent(n), href))
				return // Skip processing child nodes of the <a> tag
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)
	return textContent.String()
}

// getAttrValue retrieves the value of the specified attribute from a node
func getAttrValue(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

// getTextContent retrieves the text content of a node
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.FirstChild != nil {
		return getTextContent(n.FirstChild)
	}
	return ""
}
