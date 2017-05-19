package parse

import (
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

// Lexer consumes text and emits tokens representing characters or words from
// the text.
//
// Create a new Lexer with NewLexer. Call Lex() to retrieve tokens. A token
// of zero means EOF. Lex satisfies the interface goyacc.yyLexer, which
// mandates the Lex and Error functions.
//
type Lexer struct {
	in        io.Reader
	tokenChan chan token
}

// NewLexer returns a IdlLexer from which tokens can be retrieved with Lex().
func NewLexer(in io.Reader) IdlLexer {
	l := Lexer{in: in, tokenChan: make(chan token)}
	go func() {
		s := bufio.NewScanner(l.in)
		s.Split(split)
		for s.Scan() {
			l.tokenChan <- toToken(s.Text())
		}
		close(l.tokenChan)
	}()
	return l
}

// Call to retrieve a token. IdlSymType will be updated with additional information
// associated with the token. Lex returns zero when there are no more tokens to get.
// Part of the goyacc.yyLexer interface.
func (l Lexer) Lex(v *IdlSymType) int {
	t, ok := <-l.tokenChan
	if !ok {
		return 0
	}
	v.value = t.v
	return t.t
}

// Error is called by the parser to emit errors.
// Part of the goyacc.yyLexer interface.
func (l Lexer) Error(s string) {
	parsingError = s
	TheSpecification = nil
}

// split is an implementation of bufio.SplitFunc that tokenizes IDL files.
func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 && atEOF {
		// No data in the buffer and no more data to read.
		return 0, nil, nil
	}

	// Keep track of where in data we are.
	start := 0

	// Skip whitespace.
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}

	// Did we reach the end of the data?
	if start >= len(data) {
		// We did. Nothing but whitespace, so nothing to return. Ask for more data.
		return 0, nil, nil
	}

	// Are we looking at a token character?
	if r, width := utf8.DecodeRune(data[start:]); isToken(r) {
		// If this is the last rune in the buffer and there's more then we can't tell what the next rune is. Ask for
		// more data.
		if !atEOF && start+width >= len(data) {
			return 0, nil, nil
		}
		// Are we looking at a scope (::), shift right (>>) or shift left (<<)?
		if r2, width2 := utf8.DecodeRune(data[start+width:]); (r == ':' && r2 == ':') || (r == '<' && r2 == '<') || (r == '>' && r2 == '>') {
			return start + width + width2, data[start : start+width+width2], nil
		}
		return start + width, data[start : start+width], nil
	}

	// Are we looking at a quoted string?
	if r, width := utf8.DecodeRune(data[start:]); r == '"' {
		// Find the end quote. Keep the quotes in the returned token.
		for i := start + width; start+i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if r == '"' {
				return i + width, data[start : i+width], nil
			}
		}
		// Could not find the end quote. Need more data.
		return 0, nil, nil
	}

	// Must be a multi character token. Scan until we hit something other than a valid name rune.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		// Go on until we find a rune that is not a letter, digit or period.
		if !isName(r) && r != '.' {
			return i, data[start:i], nil
		}
	}

	// If we're at EOF, then we have another non-empty word in data. Return it.
	if atEOF {
		return len(data), data[start:], nil
	}

	// There's not enough data, ask for more.
	return 0, nil, nil
}

// isSpace returns true if the rune is a space.
func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

// isName returns true if the rune in the string is a letter, digit or underscore.
func isName(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isToken returns true if the rune is a token by itself.
func isToken(r rune) bool {
	switch r {
	case '{', '}', '[', ']', '(', ')', ':', ',', ';', '=', '+', '-', '*', '/', '%', '~', '|', '^', '&', '<', '>':
		return true
	}
	return false
}
