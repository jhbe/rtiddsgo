package generate

import (
	"testing"
	"bytes"
	"rtiddsgo/parse"
	"strings"
)

func TestConsts(t *testing.T) {
	members := parse.ConstsDef{
		{"MyConst", "string", "\"Foo\""},
		{"TheName", "int", "5"},
	}
	outBuffer := bytes.Buffer{}
	err := CreateConstsFile("mypackagename", members, &outBuffer)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(outBuffer.String(), "const MyConst string = \"Foo\"") {
		t.Error("Expected to find the const, but could not.", outBuffer.String())
	}
	if !strings.Contains(outBuffer.String(), "const TheName int = 5") {
		t.Error("Expected to find the const, but could not.", outBuffer.String())
	}
}