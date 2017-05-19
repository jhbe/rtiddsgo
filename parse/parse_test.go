package parse

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	s := `
module A {
	module B {
    const string C = "Foo";
    const long D = 5;
    const long E = 6;

    enum F {
        F_One,
        F_Two
    };

		struct G {
			long id;
		};

		struct H {
			long one, two;
			sequence<G, 2> seqG;
		};

		struct I {
			A::B::G g;
			H h;
      sequence<F, D> seqF;
      sequence<H, E> seqH;
		};

		union J switch (boolean) {
		  case FALSE: bool theFalse;
		};
	};
};
`
	expected := StructsDef{
		"A_B_G": StructDef{
			[]StructMemberDef{
				{"int32", "id", "", false, ""},
			},
			"", false,
		},
		"A_B_H": StructDef{
			[]StructMemberDef{
				StructMemberDef{"int32", "one", "", false, ""},
				StructMemberDef{"int32", "two", "", false, ""},
				StructMemberDef{"A_B_G", "seqG", "2", false, ""},
			},
			"", false,
		},
		"A_B_I": StructDef{
			[]StructMemberDef{
				StructMemberDef{"A_B_G", "g", "", false, ""},
				StructMemberDef{"A_B_H", "h", "", false, ""},
				StructMemberDef{"A_B_F", "seqF", "A_B_D", true, ""},
				StructMemberDef{"A_B_H", "seqH", "A_B_E", false, ""},
			},
			"", false,
		},
		"A_B_J": StructDef{
			[]StructMemberDef{
				StructMemberDef{"bool", "theFalse", "", false, "false"},
			},
			"bool", false,
		},
	}

	structs, _, _, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Error("Did not expect an error parsing this IDL.")
	}
	if !reflect.DeepEqual(expected, structs) {
		t.Errorf("The structs...\n%v\n...did not match the expected value...\n%v\n", structs, expected)
	}
}

func TestFixEnumNames(t *testing.T) {
	// module M {
	//   enum E {    E => M_E
	//      One, Two
	//   };
	// };
	m := &Node{Name: "M", Kind: KindModule}
	e := &Node{Name: "E", Kind: KindEnum}
	one := &Node{Name: "One"}
	two := &Node{Name: "Two"}
	e.Add(one, two)
	m.Add(e)

	fixEnumNames(m)
	if e.Name != "M_E" {
		t.Error("Expected the name of the enum to be fully qualified M_E, but found", e.Name)
	}
	if e.Child("One") == nil {
		t.Error("Expected one enum value to be One, but could not find it.")
	}
	if e.Child("Two") == nil {
		t.Error("Expected one enum value to be Two, but could not find it.")
	}
}

func TestFixConstNames(t *testing.T) {
	// module A {
	//   const long B = 5;       B => A_B
	//   const long A_C = 6;     Unchanged
	// };
	a := &Node{Name: "A", Kind: KindModule}
	b := &Node{Name: "B", Kind: KindConst, Value: "5"}
	c := &Node{Name: "A_C", Kind: KindConst, Value: "6"}
	a.Add(b, c)

	fixConstNames(a)
	if b.Name != "A_B" {
		t.Error("Expected the name of B to be A_B, but got", b.Name)
	}
	if c.Name != "A_C" {
		t.Error("Expected the name of C to remain A_C, but got", c.Name)
	}
}

func TestFixConstUsage(t *testing.T) {
	// module M {
	//   const string LEN = 5;
	//   struct S {
	//     sequence<T, LEN> seqOne;  // LEN needs to be expanded to M_LEN
	//     sequence<T, 5> seqTwo;    // "5" should be unchanged
	//   };
	// };
	m := &Node{Name: "M", Kind: KindModule}
	c := &Node{Name: "LEN", Kind: KindConst, Value: "5"}
	s := &Node{Name: "S", Kind: KindType}
	seqOne := &Node{Name: "seqOne", Kind: KindMember, TypeName: "T", Length: "LEN"}
	seqTwo := &Node{Name: "seqTwo", Kind: KindMember, TypeName: "T", Length: "5"}
	s.Add(seqOne, seqTwo)
	m.Add(c, s)

	fixConstUsage(m)
	if seqOne.Length != "M_LEN" {
		t.Error("Expected the sequence length to be M_LEN, but got", seqOne.Length)
	}
	if seqTwo.Length != "5" {
		t.Error("Expected the sequence length to remain 5, but got", seqTwo.Length)
	}
}

func TestFixTypeNames(t *testing.T) {
	// module A {
	//   struct B {
	//     double C;
	//   };
	//   struct D {
	//      B E;
	//      A_B F;
	//      sequence<B, 5> G
	//   };
	// };
	//
	// Member "E" is defined as "B", which is legal as "B" is part of the same
	// module. But we want fully qualified type names for all types that are not
	// Go built-in types. So in effect any type defined in the IDL.
	//
	// Member "F" is already qualified and should not change.
	//
	// Member "G" should also have its type name fully qualified.
	//
	a := &Node{Name: "A", Kind: KindModule}
	b := &Node{Name: "B", Kind: KindType}
	c := &Node{Name: "C", Kind: KindBaseMember, TypeName: "double"}
	d := &Node{Name: "D", Kind: KindType}
	e := &Node{Name: "E", Kind: KindMember, TypeName: "B"}
	f := &Node{Name: "F", Kind: KindMember, TypeName: "A_B"}
	g := &Node{Name: "G", Kind: KindMember, TypeName: "B", Length: "5"}
	a.Add(b)
	b.Add(c)
	a.Add(d)
	d.Add(e, f, g)

	fixTypeNames(a)
	if e.TypeName != "A_B" {
		t.Errorf("Expected the type name for E to be fully qualified, but found it to be \"%s\"", e.TypeName)
	}
	if f.TypeName != "A_B" {
		t.Errorf("Expected the type name for F to remain unchanged, but found it to be \"%s\"", f.TypeName)
	}
	if g.TypeName != "A_B" {
		t.Errorf("Expected the type name for G to be fully qualified, but found it to be \"%s\"", g.TypeName)
	}
}
