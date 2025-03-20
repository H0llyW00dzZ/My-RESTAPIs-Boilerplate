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
		listIndent:   0,
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
	for i := 0; i < count; i++ {
		s.builder.WriteString(s.nl)
	}
	s.needSpace = false
}

// elementStartHandlers processes the opening of HTML elements
//
// Note: This reduces cyclomatic complexity by avoiding numerous "if-else" statements, switch cases, and for loops.
var elementStartHandlers = map[string]func(*textState, *html.Node){
	"br":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"p":   func(s *textState, n *html.Node) { s.addNewline(2) },
	"div": func(s *textState, n *html.Node) { s.addNewline(2) },
	"h1":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h2":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h3":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h4":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h5":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h6":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"ul":  func(s *textState, n *html.Node) { s.inList = true; s.addNewline(1) },
	"ol":  func(s *textState, n *html.Node) { s.inList = true; s.addNewline(1) },
	"li": func(s *textState, n *html.Node) {
		if s.inList {
			s.builder.WriteString("- ")
		}
		s.needSpace = false
	},
	"a":     func(s *textState, n *html.Node) { s.processAnchorStart(n) },
	"img":   func(s *textState, n *html.Node) { s.processImage(n) },
	"table": func(s *textState, n *html.Node) { s.inTable = true; s.headerParsed = false; s.addNewline(1) },
	// Note: The tr, td, and th elements should now be correct and will only display formatting if they are inside a table.
	"tr": func(s *textState, n *html.Node) {
		if s.inTable {
			s.builder.WriteString("| ")
		}
		s.needSpace = false
	},
	"td": func(s *textState, n *html.Node) {
		if s.inTable {
			s.builder.WriteString(" ")
		}
		s.needSpace = false
	},
	"th": func(s *textState, n *html.Node) {
		if s.inTable {
			s.builder.WriteString(" ")
		}
		s.needSpace = false
	},
}

// elementEndHandlers processes the closing of HTML elements
//
// Note: This reduces cyclomatic complexity by avoiding numerous "if-else" statements, switch cases, and for loops.
var elementEndHandlers = map[string]func(*textState, *html.Node){
	"p":   func(s *textState, n *html.Node) { s.addNewline(2) },
	"div": func(s *textState, n *html.Node) { s.addNewline(2) },
	"h1":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h2":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h3":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h4":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h5":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"h6":  func(s *textState, n *html.Node) { s.addNewline(1) },
	"li": func(s *textState, n *html.Node) {
		if s.inList {
			s.addNewline(1)
		}
	},
	"ul":    func(s *textState, n *html.Node) { s.inList = false; s.addNewline(1) },
	"ol":    func(s *textState, n *html.Node) { s.inList = false; s.addNewline(1) },
	"a":     func(s *textState, n *html.Node) { s.processAnchorEnd(n) },
	"table": func(s *textState, n *html.Node) { s.inTable = false; s.addNewline(2) },
	// Note: The tr, td, and th elements should now be correct and will only display formatting if they are inside a table.
	"tr": func(s *textState, n *html.Node) {
		if s.inTable {
			s.addNewline(1)
			if !s.headerParsed {
				s.headerParsed = true
				s.addHeaderSeparator()
			}
		}
	},
	"td": func(s *textState, n *html.Node) {
		if s.inTable {
			s.builder.WriteString(" |")
			s.needSpace = true
		}
	},
	"th": func(s *textState, n *html.Node) {
		if s.inTable {
			s.builder.WriteString(" |")
			s.needSpace = true
		}
	},
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
