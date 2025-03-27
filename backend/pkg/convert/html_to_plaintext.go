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

// HTMLToPlainText is an [AST]-based interpreter that converts HTML content to plain text.
// It parses the HTML and extracts text nodes, concatenating them into a single string.
// If parsing fails, it returns the original HTML content as a fallback.
//
// TODO: Improving this will require additional filtering, possibly using regex.
//
// [AST]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
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
		nl:           getNewline(),
		inTable:      false,
		headerParsed: false,
		headerSizes:  []int{},
	}

	state.extractText(doc)
	return state.builder.String()
}

// textState maintains the state during text extraction
type textState struct {
	builder      *strings.Builder
	needSpace    bool
	inList       bool
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

// extractText processes HTML nodes and extracts text content.
//
// Note: If you are familiar with [AST] interpreters and you are bored with other programming languages,
// it is possible to create a new language where you write HTML and convert it into Go code or any other language.
// In this case, Go can act as a compiler for your new language.
//
// TODO: This implementation is still unfinished because HTML is complex. The table extraction functionality also needs improvement.
//
// [AST]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
func (s *textState) extractText(n *html.Node) {
	if shouldSkipNode(n) {
		return
	}

	switch n.Type {
	case html.TextNode:
		s.processTextNode(n)
	case html.ElementNode:
		if startFunc, exists := elementStartHandlers[n.Data]; exists {
			startFunc(s, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			s.extractText(c)
		}
		if endFunc, exists := elementEndHandlers[n.Data]; exists {
			endFunc(s, n)
		}
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			s.extractText(c)
		}
	}
}

// processTextNode handles text node content
func (s *textState) processTextNode(n *html.Node) {
	text := strings.TrimSpace(n.Data)
	if text != "" {
		if s.needSpace {
			s.builder.WriteString(" ")
		}
		s.builder.WriteString(text)
		// Set needSpace after writing text
		s.needSpace = true

		if s.inTable && !s.headerParsed {
			s.headerSizes = append(s.headerSizes, len(text))
		}
	} else {
		// Reset needSpace if text is empty
		s.needSpace = false
	}
}

// addNewline adds the specified number of newline characters to the builder.
// It also resets the needSpace flag to ensure proper spacing in the output.
func (s *textState) addNewline(count int) {
	for range count {
		s.builder.WriteString(s.nl)
	}
	s.needSpace = false
}

// elementStartHandlers processes the opening of HTML elements
//
// Note: This reduces cyclomatic complexity by avoiding numerous "if-else" statements, switch cases, and for loops.
// The performance might differ from the previous implementation that used switch cases due to the [trade-off].
// However, using switch cases can become harder to maintain when there are many cases (not good 👎).
// Additionally, this pattern (using a function map) is designed to handle different behaviors based on the type of HTML element encountered,
// providing a clean and flexible way to manage element-specific logic.
//
// TODO: Add support for additional elements here later. This helper function is unfinished, and it is designed for easy extension to handle new elements.
//
// [trade-off]: https://en.wikipedia.org/wiki/Trade-off
var elementStartHandlers = map[string]func(*textState, *html.Node){
	"br":    oneNewLine,
	"p":     twoNewLines,
	"div":   twoNewLines,
	"h1":    oneNewLine,
	"h2":    oneNewLine,
	"h3":    oneNewLine,
	"h4":    oneNewLine,
	"h5":    oneNewLine,
	"h6":    oneNewLine,
	"ul":    inList,
	"ol":    inList,
	"li":    listItem,
	"a":     func(s *textState, n *html.Node) { s.processAnchorStart(n) },
	"img":   func(s *textState, n *html.Node) { s.processImage(n) },
	"table": func(s *textState, n *html.Node) { s.inTable = true; s.headerParsed = false; s.addNewline(1) },
	// Note: The tr, td, and th elements should now be correct and will only display formatting if they are inside a table.
	"tr": tableRowStart,
	"td": tableCellStart,
	"th": tableCellStart,
}

// elementEndHandlers processes the closing of HTML elements
//
// Note: This reduces cyclomatic complexity by avoiding numerous "if-else" statements, switch cases, and for loops.
// The performance might differ from the previous implementation that used switch cases due to the [trade-off].
// However, using switch cases can become harder to maintain when there are many cases (not good 👎).
// Additionally, this pattern (using a function map) is designed to handle different behaviors based on the type of HTML element encountered,
// providing a clean and flexible way to manage element-specific logic.
//
// TODO: Add support for additional elements here later. This helper function is unfinished, and it is designed for easy extension to handle new elements.
//
// [trade-off]: https://en.wikipedia.org/wiki/Trade-off
var elementEndHandlers = map[string]func(*textState, *html.Node){
	"p":     twoNewLines,
	"div":   twoNewLines,
	"h1":    oneNewLine,
	"h2":    oneNewLine,
	"h3":    oneNewLine,
	"h4":    oneNewLine,
	"h5":    oneNewLine,
	"h6":    oneNewLine,
	"li":    listItemEnd,
	"ul":    listEnd,
	"ol":    listEnd,
	"a":     func(s *textState, n *html.Node) { s.processAnchorEnd(n) },
	"table": tableEnd,
	// Note: The tr, td, and th elements should now be correct and will only display formatting if they are inside a table.
	"tr": tableRowEnd,
	"td": tableCellEnd,
	"th": tableCellEnd,
}

// addHeaderSeparator adds a markdown separator line for table headers.
// It uses the lengths of the header text to ensure the separator aligns properly.
func (s *textState) addHeaderSeparator() {
	s.builder.WriteString("|")
	for _, size := range s.headerSizes {
		// this should be fine, even with 1 billion tables; it won't overflow like Unix time.
		s.builder.WriteString(strings.Repeat("-", size+2) + "|")
	}
	s.builder.WriteString(s.nl)
	s.headerSizes = nil
}

// processAnchorStart handles the start of anchor tags
func (s *textState) processAnchorStart(n *html.Node) {
	var href string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}
	if href != "" {
		if s.needSpace {
			// Add space before the link if needed
			s.builder.WriteString(" ")
			s.needSpace = false
		}
		s.builder.WriteString("[")
		// Store href for later use
		n.Attr = append(n.Attr, html.Attribute{Key: "_stored_href", Val: href})
	}
}

// processAnchorEnd handles the end of anchor tags
func (s *textState) processAnchorEnd(n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "_stored_href" {
			s.builder.WriteString("](")
			s.builder.WriteString(attr.Val)
			s.builder.WriteString(")")
			// Set needSpace after the link
			s.needSpace = true
			break
		}
	}
}

// processImage handles image tags
func (s *textState) processImage(n *html.Node) {
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
		s.builder.WriteString("![")
		s.builder.WriteString(alt)
		s.builder.WriteString("](")
		s.builder.WriteString(src)
		s.builder.WriteString(")")
	}
}

// HTMLToPlainTextStreams is an [AST]-based interpreter that converts HTML content from an input stream to plain text
// and writes it to an output stream (a.k.a Hybrid Streaming).
//
// TODO: Improving this will require additional filtering, possibly using regex.
//
// [AST]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
func HTMLToPlainTextStreams(i io.Reader, o io.Writer) error {
	builder := getBuilder()
	defer putBuilder(builder)

	// It's better to handle this directly without bufio, allowing the caller to manage buffering.
	doc, err := html.Parse(i)
	if err != nil {
		return err // Return error if parsing fails
	}

	state := &textState{
		builder:      builder,
		needSpace:    false,
		inList:       false,
		nl:           getNewline(),
		inTable:      false,
		headerParsed: false,
		headerSizes:  []int{},
	}

	state.extractText(doc)
	_, err = o.Write([]byte(state.builder.String()))
	return err
}

// HTMLToPlainTextConcurrent is an [AST]-based interpreter that converts multiple HTML strings to plain text concurrently.
// It returns a slice of plain text results corresponding to each HTML input.
//
// Note: This is designed for high-performance scenarios. It also depends on the number of available CPU cores,
// unlike [HTMLToPlainTextStreamsConcurrent], which depends on the input reader.
//
// [AST]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
func HTMLToPlainTextConcurrent(htmlContents []string) []string {
	results := make([]string, len(htmlContents))
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup

	chunkSize := (len(htmlContents) + numWorkers - 1) / numWorkers
	// Launch a goroutine for each worker to process a chunk of the input concurrently.
	for i := range numWorkers {
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

// HTMLToPlainTextStreamsConcurrent is an [AST]-based interpreter that processes multiple readers concurrently
// and writes the plain text to a single writer, returning any errors encountered
//
// Note: This is designed for high-performance scenarios, like non-stop 24/7 streaming hahaha.
// It's where your machine really earns its keep—no coffee breaks here!
//
// [AST]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
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
