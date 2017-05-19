package parse

import "testing"

func TestGetConstsDef(t *testing.T) {
	// module M {
	//   const string theString = "Foo";
	// };
	m := &Node{Name: "M", Kind: KindModule}
	c := &Node{Name: "theString", Kind: KindConst, TypeName: "string", Value: "\"Foo\""}
	m.Add(c)

	consts := getConstsDef(m)
	if len(consts) != 1 {
		t.Error("Expected one const")
	} else {
		if consts[0].Name != "theString" {
			t.Error("Expected the name to be theString, but was", consts[0].Name)
		}
		if consts[0].Type != "string" {
			t.Error("Expected the type to be string, but was", consts[0].Type)
		}
		if consts[0].Value != "\"Foo\"" {
			t.Error("Expected the value to be Foo, but was", consts[0].Value)
		}
	}
}
