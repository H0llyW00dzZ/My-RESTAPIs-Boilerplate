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
	var extractText func(*html.Node, bool)

	extractText = func(n *html.Node, inList bool) {
		if n.Type == html.TextNode {
			textContent.WriteString(n.Data)
		} else if n.Type == html.ElementNode {
			switch n.Data {
			case "br":
				textContent.WriteString("\n")
			case "p":
				textContent.WriteString("\n\n")
			case "h1", "h2", "h3", "h4", "h5", "h6":
				textContent.WriteString("\n")
			case "ul", "ol":
				inList = true
				textContent.WriteString("\n")
			case "li":
				if inList {
					textContent.WriteString("- ")
				}
			case "div":
				textContent.WriteString("\n\n")
			case "a":
				href := ""
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}

				// Note: This is what it will look like in markdown format "Visit [Example](https://example.com) website."
				textContent.WriteString(fmt.Sprintf("[%s](%s)", n.FirstChild.Data, href))
				return // Skip processing child nodes of the <a> tag
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c, inList)
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "li":
				if inList {
					textContent.WriteString("\n")
				}
			case "ul", "ol":
				inList = false
				textContent.WriteString("\n")
			case "p":
				textContent.WriteString("\n\n")
			case "div":
				textContent.WriteString("\n\n")
			case "h1", "h2", "h3", "h4", "h5", "h6":
				textContent.WriteString("\n")
			}
		}
	}

	extractText(doc, false)
	return textContent.String()
}
