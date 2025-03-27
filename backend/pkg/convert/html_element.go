// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import "golang.org/x/net/html"

// oneNewLine adds a single newline to the text state.
func oneNewLine(s *textState, n *html.Node) { s.addNewline(1) }

// twoNewLines adds two newlines to the text state.
func twoNewLines(s *textState, n *html.Node) { s.addNewline(2) }

// inList marks the state as being inside a list and adds a newline.
func inList(s *textState, n *html.Node) { s.inList = true; s.addNewline(1) }

// listItem adds a list item prefix if inside a list and resets spacing.
func listItem(s *textState, n *html.Node) {
	if s.inList {
		s.builder.WriteString("- ")
	}
	s.needSpace = false
}

// listItemEnd adds a newline after a list item if inside a list.
func listItemEnd(s *textState, n *html.Node) {
	if s.inList {
		s.addNewline(1)
	}
}

// listEnd marks the end of a list and adds a newline.
func listEnd(s *textState, n *html.Node) { s.inList = false; s.addNewline(1) }

// tableRowStart begins a table row with a pipe character if inside a table.
func tableRowStart(s *textState, n *html.Node) {
	if s.inTable {
		s.builder.WriteString("| ")
	}
	s.needSpace = false
}

// tableCellStart begins a table cell with a space if inside a table.
func tableCellStart(s *textState, n *html.Node) {
	if s.inTable {
		s.builder.WriteString(" ")
	}
	s.needSpace = false
}

// tableRowEnd ends a table row, adds a newline, and handles header parsing.
func tableRowEnd(s *textState, n *html.Node) {
	if s.inTable {
		s.addNewline(1)
		if !s.headerParsed {
			s.headerParsed = true
			s.addHeaderSeparator()
		}
	}
}

// tableCellEnd ends a table cell with a pipe character if inside a table.
func tableCellEnd(s *textState, n *html.Node) {
	if s.inTable {
		s.builder.WriteString(" |")
		s.needSpace = true
	}
}

// tableEnd marks the end of a table and adds two newlines.
func tableEnd(s *textState, n *html.Node) { s.inTable = false; s.addNewline(2) }
