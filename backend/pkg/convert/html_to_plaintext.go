// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"golang.org/x/net/html"
)

// HTMLToPlainText converts HTML content to plain text.
// It parses the HTML and extracts text nodes, concatenating them into a single string.
// If parsing fails, it returns the original HTML content as a fallback.
//
// Note: This function does not fully handle elements like "<script>" or other non-text content.
//
// TODO: Improving this will require additional filtering, possibly using regex.
func HTMLToPlainText(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent // Fallback to original HTML if parsing fails
	}

	var textContent strings.Builder
	extractText(doc, &textContent, false)
	return textContent.String()
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

// handleElementNode processes HTML element nodes and appends corresponding
// plain text representations to the textContent based on the tag type.
func handleElementNode(n *html.Node, textContent *strings.Builder, inList *bool) {
	newline := getNewline()
	switch n.Data {
	case "br":
		textContent.WriteString(newline)
	case "p", "div":
		textContent.WriteString(newline + newline)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		textContent.WriteString(newline)
	case "ul", "ol":
		*inList = true
		textContent.WriteString(newline)
	case "li":
		if *inList {
			textContent.WriteString("- ")
		}
		// TODO: This case for "a" might be unnecessary; will remove it later.
	case "a":
		handleAnchorTag(n, textContent)
	case "img":
		handleImageTag(n, textContent)
	}
}

// extractText recursively traverses the HTML node tree, converting nodes
// to plain text and appending them to textContent.
func extractText(n *html.Node, textContent *strings.Builder, inList bool) {
	if n.Type == html.TextNode {
		textContent.WriteString(n.Data)
	} else if n.Type == html.ElementNode {
		if n.Data == "style" {
			return // Skip <style> tags entirely
		}
		if n.Data == "a" {
			handleAnchorTag(n, textContent)
			return // Skip further processing for child nodes of <a>
		}
		handleElementNode(n, textContent, &inList)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, textContent, inList)
	}

	if n.Type == html.ElementNode {
		handleClosingTags(n, textContent, &inList)
	}
}

// handleAnchorTag processes <a> tags, extracting the href attribute and
// text content, then appending a markdown formatted link to textContent.
func handleAnchorTag(n *html.Node, textContent *strings.Builder) {
	href := ""
	linkText := ""

	// Extract href attribute
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	// Extract the text content of the link
	if n.FirstChild != nil {
		linkText = n.FirstChild.Data
	}

	// Append markdown formatted link
	//
	// Note: This is what it will look like in markdown format "Visit [Example](https://example.com) website."
	textContent.WriteString(fmt.Sprintf("[%s](%s)", linkText, href))
}

// handleClosingTags appends appropriate plain text representations
// for closing tags, managing list states and formatting.
func handleClosingTags(n *html.Node, textContent *strings.Builder, inList *bool) {
	newline := getNewline()
	switch n.Data {
	case "li":
		if *inList {
			textContent.WriteString(newline)
		}
	case "ul", "ol":
		*inList = false
		textContent.WriteString(newline)
	case "p", "div":
		textContent.WriteString(newline + newline)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		textContent.WriteString(newline)
	}
}

// HTMLToPlainTextStreams converts HTML content from an input stream to plain text
// and writes it to an output stream.
//
// Note: This function does not fully handle elements like "<script>" or other non-text content.
//
// TODO: Improving this will require additional filtering, possibly using regex.
func HTMLToPlainTextStreams(i io.Reader, o io.Writer) error {
	doc, err := html.Parse(i)
	if err != nil {
		return err // Return error if parsing fails
	}

	var textContent strings.Builder
	extractText(doc, &textContent, false)

	_, err = o.Write([]byte(textContent.String()))
	return err
}

// handleImageTag processes <img> tags.
//
// TODO: Automatically handle the size of the image as well.
func handleImageTag(n *html.Node, textContent *strings.Builder) {
	src := ""
	alt := ""

	for _, attr := range n.Attr {
		if attr.Key == "src" {
			src = attr.Val
		}
		if attr.Key == "alt" {
			alt = attr.Val
		}
	}

	textContent.WriteString(fmt.Sprintf("![%s](%s)", alt, src))
}
