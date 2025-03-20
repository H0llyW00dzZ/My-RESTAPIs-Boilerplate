// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"bufio"
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
// TODO: Improving this will require additional filtering, possibly using regex.
func HTMLToPlainText(htmlContent string) string {
	builder := getBuilder()
	defer putBuilder(builder)

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	state := &textState{
		builder:      builder,
		needSpace:    false,
		inList:       false,
		listIndent:   0,
		nl:           getNewline(),
		inTable:      false,
		headerParsed: false,
		headerSizes:  []int{},
	}

	extractText(doc, state)
	return state.builder.String()
}

// textState maintains the state during text extraction
type textState struct {
	builder      *strings.Builder
	needSpace    bool
	inList       bool
	listIndent   int
	nl           string
	inTable      bool
	headerParsed bool
	headerSizes  []int
}

// shouldSkipNode determines if a node should be skipped during processing
func shouldSkipNode(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	switch n.Data {
	case "script", "style", "noscript":
		return true
	}
	return false
}

// extractText processes HTML nodes and extracts text content
//
// TODO: This is still unfinished because HTML is complex. The table extraction also needs improvement.
func extractText(n *html.Node, state *textState) {
	if shouldSkipNode(n) {
		return
	}

	switch n.Type {
	case html.TextNode:
		processTextNode(n, state)
	case html.ElementNode:
		handleElementStart(n, state)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c, state)
		}
		handleElementEnd(n, state)
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c, state)
		}
	}
}

// processTextNode handles text node content
func processTextNode(n *html.Node, state *textState) {
	text := strings.TrimSpace(n.Data)
	if text != "" {
		if state.needSpace {
			state.builder.WriteString(" ")
		}
		state.builder.WriteString(text)
		// Set needSpace after writing text
		state.needSpace = true

		if state.inTable && !state.headerParsed {
			state.headerSizes = append(state.headerSizes, len(text))
		}
	} else {
		// Reset needSpace if text is empty
		state.needSpace = false
	}
}

// handleElementStart processes the opening of HTML elements
func handleElementStart(n *html.Node, state *textState) {
	switch n.Data {
	case "br":
		state.builder.WriteString(state.nl)
		state.needSpace = false
	case "p", "div":
		state.builder.WriteString(state.nl + state.nl)
		state.needSpace = false
	case "h1", "h2", "h3", "h4", "h5", "h6":
		state.builder.WriteString(state.nl)
		state.needSpace = false
	case "ul", "ol":
		state.inList = true
		state.builder.WriteString(state.nl)
		state.needSpace = false
	case "li":
		if state.inList {
			state.builder.WriteString("- ")
		}
		state.needSpace = false
	case "a":
		processAnchorStart(n, state)
	case "img":
		processImage(n, state)
	case "table":
		state.inTable = true
		state.headerParsed = false
		state.builder.WriteString(state.nl)
		state.needSpace = false
		// Note: The tr, td, and th elements should now be correct and will only display formatting if they are inside a table.
	case "tr":
		if state.inTable {
			state.builder.WriteString("| ")
		}
		state.needSpace = false
	case "td", "th":
		if state.inTable {
			state.builder.WriteString(" ")
		}
		state.needSpace = false
	}
}

// handleElementEnd processes the closing of HTML elements
func handleElementEnd(n *html.Node, state *textState) {
	switch n.Data {
	case "p", "div":
		state.builder.WriteString(state.nl + state.nl)
		state.needSpace = false
	case "h1", "h2", "h3", "h4", "h5", "h6":
		state.builder.WriteString(state.nl)
		state.needSpace = false
	case "li":
		if state.inList {
			state.builder.WriteString(state.nl)
		}
		state.needSpace = false
	case "ul", "ol":
		state.inList = false
		state.builder.WriteString(state.nl)
		state.needSpace = false
	case "a":
		processAnchorEnd(n, state)
	case "table":
		state.inTable = false
		state.builder.WriteString(state.nl + state.nl)
		state.needSpace = false
		// Note: The tr, td, and th elements should now be correct and will only display formatting if they are inside a table.
	case "tr":
		if state.inTable {
			state.builder.WriteString(state.nl)
			if !state.headerParsed {
				state.headerParsed = true
				addHeaderSeparator(state)
			}
		}
		state.needSpace = false
	case "td", "th":
		if state.inTable {
			state.builder.WriteString(" |")
			state.needSpace = true
		}
	}
}

// addHeaderSeparator adds a markdown separator line for table headers.
// It uses the lengths of the header text to ensure the separator aligns properly.
func addHeaderSeparator(state *textState) {
	state.builder.WriteString("|")
	// this should be fine, even with 1 billion tables; it won't overflow like Unix time.
	for _, size := range state.headerSizes {
		state.builder.WriteString(strings.Repeat("-", size+2) + "|")
	}
	state.builder.WriteString(state.nl)
	state.headerSizes = nil
}

// processAnchorStart handles the start of anchor tags
func processAnchorStart(n *html.Node, state *textState) {
	var href string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}
	if href != "" {
		if state.needSpace {
			// Add space before the link if needed
			state.builder.WriteString(" ")
			state.needSpace = false
		}
		state.builder.WriteString("[")
		// Store href for later use
		n.Attr = append(n.Attr, html.Attribute{Key: "_stored_href", Val: href})
	}
}

// processAnchorEnd handles the end of anchor tags
func processAnchorEnd(n *html.Node, state *textState) {
	for _, attr := range n.Attr {
		if attr.Key == "_stored_href" {
			state.builder.WriteString("](")
			state.builder.WriteString(attr.Val)
			state.builder.WriteString(")")
			// Set needSpace after the link
			state.needSpace = true
			break
		}
	}
}

// processImage handles image tags
func processImage(n *html.Node, state *textState) {
	var alt, src string
	for _, attr := range n.Attr {
		switch attr.Key {
		case "alt":
			alt = attr.Val
		case "src":
			src = attr.Val
		}
	}
	if src != "" {
		state.builder.WriteString("![")
		state.builder.WriteString(alt)
		state.builder.WriteString("](")
		state.builder.WriteString(src)
		state.builder.WriteString(")")
	}
}

// HTMLToPlainTextStreams converts HTML content from an input stream to plain text
// and writes it to an output stream (a.k.a Hybrid Streaming).
//
// TODO: Improving this will require additional filtering, possibly using regex.
func HTMLToPlainTextStreams(i io.Reader, o io.Writer) error {
	builder := getBuilder()
	defer putBuilder(builder)

	doc, err := html.Parse(bufio.NewReader(i))
	if err != nil {
		return err // Return error if parsing fails
	}

	state := &textState{
		builder:      builder,
		needSpace:    false,
		inList:       false,
		listIndent:   0,
		nl:           getNewline(),
		inTable:      false,
		headerParsed: false,
		headerSizes:  []int{},
	}

	extractText(doc, state)
	_, err = o.Write([]byte(state.builder.String()))
	return err
}

// HTMLToPlainTextConcurrent converts multiple HTML strings to plain text concurrently.
// It returns a slice of plain text results corresponding to each HTML input.
//
// Note: This is designed for high-performance scenarios. It also depends on the number of available CPU cores,
// unlike [HTMLToPlainTextStreamsConcurrent], which depends on the input reader.
func HTMLToPlainTextConcurrent(htmlContents []string) []string {
	results := make([]string, len(htmlContents))
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup

	chunkSize := (len(htmlContents) + numWorkers - 1) / numWorkers
	// Launch a goroutine for each worker to process a chunk of the input concurrently.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := min(start+chunkSize, len(htmlContents))

			for j := start; j < end; j++ {
				results[j] = HTMLToPlainText(htmlContents[j])
			}
		}(i * chunkSize)
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
