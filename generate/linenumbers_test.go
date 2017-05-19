package generate

import "testing"

func TestWithLineNumbers(t *testing.T) {
	s := `A
B
C`
	expected := `1: A
2: B
3: C
`
	if out := withLineNumbers(s); out != expected {
		t.Errorf("Expected\n%s\n, but got: \n%s\n", expected, out)
	}
}