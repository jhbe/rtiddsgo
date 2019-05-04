package generate

import "testing"

func TestType (t *testing.T) {
	cases := []struct{
		goType, seqLen, arrayDims string
		expected string
	} {
		{"", "", "", ""},
		{"bool", "", "", "bool"},
		{"bool", "3", "", "[]bool"},
		{"bool", "", "2", "[2]bool"},
		{"bool", "3", "2", "[2][]bool"},
		{"Com_This_That", "3", "2", "[2][]Com_This_That"},
	}

	for i, c := range cases {
		out := Type(c.goType, c.seqLen, c.arrayDims)
		if out != c.expected {
			t.Errorf("Case %d: Expected %s, got %s", i, c.expected, out)
		}
	}
}