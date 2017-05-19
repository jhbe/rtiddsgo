package generate

import (
"testing"
"bytes"
"rtiddsgo/parse"
	"strings"
)

func TestEnums(t *testing.T) {
	members := parse.EnumsDef{
		{"M_MyFirstEnum", []string{"MyFirst", "MySecond"}},
		// Check that names and types and Go-ified (capital leading letter at a minimum)
		{"m_mySecondEnum", []string{"second"}},
	}
	outBuffer := bytes.Buffer{}
	err := CreateEnumsFile("mypackagename", members, &outBuffer)
	if err != nil {
		t.Error(err)
	}

	expected := `package mypackagename

type M_MyFirstEnum int
const (
	MyFirst M_MyFirstEnum = iota
	MySecond
)

type M_MySecondEnum int
const (
	Second M_MySecondEnum = iota
)`
	if !strings.Contains(outBuffer.String(), expected) {
		t.Error("Expected this enum file:", expected, "but got", outBuffer.String())
	}
}