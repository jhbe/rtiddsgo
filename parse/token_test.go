package parse

import "testing"

func TestToString(t *testing.T) {
	cases := []struct {
		tok      token
		expected string
	}{
		{token{0, ""}, "EOF"},
		{token{T_MODULE, ""}, "T_MODULE"},
		{token{T_STRING_LITERAL, "FOO"}, "T_STRING_LITERAL \"FOO\""},
		{token{T_INTEGER_LITERAL, "12"}, "T_INTEGER_LITERAL \"12\""},
		{token{T_FLOATING_PT_LITERAL, "23.9"}, "T_FLOATING_PT_LITERAL \"23.9\""},
	}

	for index, c := range cases {
		s := c.tok.String()
		if s != c.expected {
			t.Error("Case ", index, ": Expected ", c.expected, ", but got ", s)
		}
	}
}

func TestToToken(t *testing.T) {
	cases := []struct {
		in  string
		tok token
	}{
		{"module", token{T_MODULE, ""}},
		{"Module", token{T_MODULE, ""}},
		{"mODULE", token{T_MODULE, ""}},
		{"MODULE", token{T_MODULE, ""}},

		{"\"12\"", token{T_STRING_LITERAL, "12"}},
		{"12", token{T_INTEGER_LITERAL, "12"}},
		{"12.3", token{T_FLOATING_PT_LITERAL, "12.3"}},

		{";", token{T_SEMICOLON, ""}},
	}

	for index, c := range cases {
		tok := toToken(c.in)
		if tok != c.tok {
			t.Error("Case ", index, ": Expected ", c.tok, ", but got ", tok)
		}
	}

}
