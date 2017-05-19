package parse

import (
	"strings"
	"testing"
)

func TestTokens(t *testing.T) {
	cases := []struct {
		in  string // Input to the lexer
		out []int  // The expected sequences of token integers.
	}{
		// An empty string should return the end token.
		{"", []int{0}},

		// Test leading and trailing whitespace.
		{"  module  ", []int{T_MODULE, 0}},
		{"module", []int{T_MODULE, 0}},
		{"  module", []int{T_MODULE, 0}},
		{"module  ", []int{T_MODULE, 0}},

		// Test identifiers versus integers and floats.
		{"12", []int{T_INTEGER_LITERAL, 0}},
		{"12.3", []int{T_FLOATING_PT_LITERAL, 0}},
		{"blah", []int{T_IDENTIFIER, 0}},
		{" blah 12 12.3 ", []int{T_IDENTIFIER, T_INTEGER_LITERAL, T_FLOATING_PT_LITERAL, 0}},

		// Test strings.
		{"struct \"gaah blah\" module", []int{T_STRUCT, T_STRING_LITERAL, T_MODULE, 0}},

		// Brackets. With and without newlines.
		{"module Gaah {};", []int{T_MODULE, T_IDENTIFIER, T_LEFT_CURLY_BRACKET, T_RIGHT_CURLY_BRACKET, T_SEMICOLON, 0}},
		{`
module Gaah {};
`, []int{T_MODULE, T_IDENTIFIER, T_LEFT_CURLY_BRACKET, T_RIGHT_CURLY_BRACKET, T_SEMICOLON, 0}},

		// Semicolons at the end immediately after an identifier.
		{"long gaah;", []int{T_LONG, T_IDENTIFIER, T_SEMICOLON, 0}},

		// Double colons, a.k.a scope.
		{"com::this;", []int{T_IDENTIFIER, T_SCOPE, T_IDENTIFIER, T_SEMICOLON, 0}},

		// Shift left and right.
		{"<<>>;", []int{T_SHIFTLEFT, T_SHIFTRIGHT, T_SEMICOLON, 0}},
	}

	for index, c := range cases {
		l := NewLexer(strings.NewReader(c.in))
		var st IdlSymType
		for i, expectedToken := range c.out {
			tok := l.Lex(&st)
			if tok != expectedToken {
				t.Errorf("Case %d, iteration %d: Expected %s, got %s", index, i, token{expectedToken, ""}, token{tok, ""})
			}
		}
	}
}

func TestStringValue(t *testing.T) {
	s := " \"gaah blah\" "
	l := NewLexer(strings.NewReader(s))

	var st IdlSymType
	tok := l.Lex(&st)
	if tok != T_STRING_LITERAL {
		t.Error("Expected a string literal.")
	}
	// Verify that the token.value holds the string.
	if st.value != "gaah blah" {
		t.Errorf("Expected the symtype value to be \"gaah blah\", but got \"%s\"", st.value)
	}
}
