package parse

import "testing"

func TestGetFullName(t *testing.T) {
	// Nil node.
	if name := getFullName("foo", nil); name != "" {
		t.Error("Expected a nil node to return an empty string.")
	}

	a := &Node{Name: "A"}
	b := &Node{Name: "B", Kind: KindModule}
	c := &Node{Name: "C"}
	b.Add(c)
	a.Add(b)

	// Base name
	if name := getFullName("bool", c); name != "bool" {
		t.Error("Expected bool to return the same.")
	}

	// Normal name
	if name := getFullName("MyType", c); name != "A_B_MyType" {
		t.Error("Expected A_B_MyType, got", name)
	}

	a.Kind = KindModule
	// Name is already fully qualified.
	if name := getFullName("A_B_C", c); name != "A_B_C" {
		t.Error("Expected A_B_C, got", name)
	}
}

func TestModulePath(t *testing.T) {
	if p := modulePath(nil); p != "" {
		t.Error("Expected the path of nil to be an empty string, but got", p)
	}

	a := &Node{Name: "A"}
	b := &Node{Name: "B", Kind: KindModule}
	c := &Node{Name: "C"}
	b.Add(c)
	a.Add(b)

	if p := modulePath(c); p != "A_B" {
		t.Error("Expected the path to be A_B, but got", p)
	}
}

func TestGetTopmostModuleName(t *testing.T) {
	a := &Node{Name: "A"}
	b := &Node{Name: "B", Kind: KindModule}
	c := &Node{Name: "C"}
	b.Add(c)
	a.Add(b)

	if m := getTopmostModuleName(a); m != b.Name {
		t.Error("Expected b to be the topmost module, but got", m)
	}
}

func TestIsBaseType(t *testing.T) {
	cases := []struct{
		in string
		out bool
	} {
		{"", false},
		{"bool", true},
		{"float", false},
		{"float32",true },
		{"string", true},
		{"String", false},
		{"int", true},
		{"int32", true},
	}

	for index, c := range cases {
		if out := isBaseType(c.in); out != c.out {
			t.Errorf("Case %d (%s): got %v, expected %v.", index, c.in, out, c.out)
		}
	}
}

func TestIsAnEnum(t *testing.T) {
	a:= &Node{Name: "A"}
	b:= &Node{Name: "B"}
	c:= &Node{Name: "C"}
	d:= &Node{Name: "D", Kind: KindEnum}
	c.Add(d)
	b.Add(c)
	a.Add(b)

	if isAnEnumType(a, "B") {
		t.Error("B is not an enum.")
	}
	if !isAnEnumType(a, "D") {
		t.Error("D is an enum.")
	}
}
