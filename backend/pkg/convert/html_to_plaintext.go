// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"io"
	"runtime"
	"strings"
	"sync"

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
	var href, linkText string

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

	if href != "" && linkText != "" {
		// Append markdown formatted link
		//
		// Note: This is what it will look like in markdown format "Visit [Example](https://example.com) website."
		textContent.WriteString("[")
		textContent.WriteString(linkText)
		textContent.WriteString("](")
		textContent.WriteString(href)
		textContent.WriteString(")")
	}
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
// and writes it to an output stream (a.k.a Hybrid Streaming).
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
	var src, alt string

	for _, attr := range n.Attr {
		switch attr.Key {
		case "src":
			src = attr.Val
		case "alt":
			alt = attr.Val
		}
	}

	if src != "" {
		textContent.WriteString("![")
		textContent.WriteString(alt)
		textContent.WriteString("](")
		textContent.WriteString(src)
		textContent.WriteString(")")
	}
}

// HTMLToPlainTextConcurrent converts multiple HTML strings to plain text concurrently.
// It returns a slice of plain text results corresponding to each HTML input.
//
// Note: This is designed for high-performance scenarios.
func HTMLToPlainTextConcurrent(htmlContents []string) []string {
	results := make([]string, len(htmlContents))
	var wg sync.WaitGroup

	// Iterate over each HTML content and process concurrently
	for i, content := range htmlContents {
		wg.Add(1)
		go func(i int, content string) {
			defer wg.Done()
			// Convert HTML to plain text and store the result
			results[i] = HTMLToPlainText(content)
		}(i, content)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return results
}

// SafeWriter wraps an [io.Writer] with a mutex for thread-safe writing
type SafeWriter struct {
	writer io.Writer
	mu     sync.Mutex
}

// Write safely writes data to the underlying writer using a mutex
func (sw *SafeWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.writer.Write(p)
}

// HTMLToPlainTextStreamsConcurrent processes multiple readers concurrently
// and writes the plain text to a single writer, returning any errors encountered
//
// Note: This is designed for high-performance scenarios, like non-stop 24/7 streaming hahaha.
// It's where your machine really earns its keep—no coffee breaks here!
func HTMLToPlainTextStreamsConcurrent(i []io.Reader, o io.Writer) []error {
	safeWriter := &SafeWriter{writer: o}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	// Iterate over each reader and process concurrently
	for _, reader := range i {
		wg.Add(1)
		go func(r io.Reader) {
			defer wg.Done()
			// Convert HTML to plain text and write to the safe writer
			if e := HTMLToPlainTextStreams(r, safeWriter); e != nil {
				// Capture any errors encountered
				mu.Lock()
				errs = append(errs, e)
				mu.Unlock()
			}
		}(reader)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return errs
}
