package rtiddsgo

import "testing"

func TestReturnCodeToStringing(t *testing.T) {
	if s := ReturnCodeToString(0); s != "OK" {
		t.Errorf("Expected \"OK\", got \"%s\"", s)
	}
	if s := ReturnCodeToString(1); s != "ERROR" {
		t.Errorf("Expected \"ERROR\", got \"%s\"", s)
	}
	if s := ReturnCodeToString(-1); s != "<unknown>" {
		t.Errorf("Expected \"OK\", got \"%s\"", s)
	}
	if s := ReturnCodeToString(42); s != "<unknown>" {
		t.Errorf("Expected \"OK\", got \"%s\"", s)
	}
}
