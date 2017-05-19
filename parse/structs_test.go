package parse

import (
	"reflect"
	"testing"
)

func TestGetStructDefs(t *testing.T) {
	// module M {
	//   const string LEN = 5;
	//   enum MyEnum {};
	//   struct T {
	//     MyEnum theEnum;
	//     long dummy;
	//   };
	//   struct S {
	//     double V;
	//     sequence<T, LEN> seq;  // T needs to be expanded to M_T. Same for LEN.
	//   };
	//   union U switch(boolean) {
	//     case false: T ut;
	//     case true: S us;
	//   };
	// };
	m := &Node{Name: "M", Kind: KindModule}
	c := &Node{Name: "LEN", Kind: KindConst, Value: "5"}
	e := &Node{Name: "M_MyEnum", Kind: KindEnum}

	_t := &Node{Name: "M_T", Kind: KindType}
	theEnum := &Node{Name: "theEnum", Kind: KindMember, TypeName: "M_MyEnum"}
	dummy := &Node{Name: "dummy", Kind: KindBaseMember, TypeName: "int32"}

	s := &Node{Name: "M_S", Kind: KindType}
	v := &Node{Name: "V", Kind: KindBaseMember, TypeName: "float64"}
	seq := &Node{Name: "seq", Kind: KindMember, TypeName: "M_T", Length: "LEN"}

	u1 := &Node{Name: "U1", Kind: KindUnionType, TypeName: "bool"}
	u1t := &Node{Name: "ut", Kind: KindMember, TypeName: "M_T", Value: "false"}
	u1s := &Node{Name: "us", Kind: KindMember, TypeName: "M_S", Value: "true"}

	u2 := &Node{Name: "U2", Kind: KindUnionType, TypeName: "M_MyEnum"}
	u2t := &Node{Name: "ut", Kind: KindMember, TypeName: "M_T", Value: "MyEnum_One"}
	u2s := &Node{Name: "us", Kind: KindMember, TypeName: "M_S", Value: "MyEnum_Two"}

	u1.Add(u1t, u1s)
	u2.Add(u2t, u2s)
	_t.Add(theEnum, dummy)
	s.Add(v, seq)
	m.Add(c, e, _t, s, u1, u2)

	structs := getStructsDef(m)

	if mt, exist := structs["M_T"]; !exist {
		t.Error("Expected a struct named M_T, but could not find it.")
	} else {
		expected := StructDef{
			[]StructMemberDef{
				{"M_MyEnum", "theEnum", "", true, ""},
				{"int32", "dummy", "", false, ""},
			},
			"",
			false,
		}
		if !reflect.DeepEqual(expected, mt) {
			t.Errorf("Got:\n\n%-v\n\nbut expected:\n\n%-v\n", mt, expected)
		}
	}

	if st, exist := structs["M_S"]; !exist {
		t.Error("Expected a struct named M_S, but could not find it.")
	} else {
		expected := StructDef{
			[]StructMemberDef{
				{"float64", "V", "", false, ""},
				{"M_T", "seq", "LEN", false, ""},
			},
			"",
			false,
		}
		if !reflect.DeepEqual(expected, st) {
			t.Errorf("Got:\n\n%-v\n\nbut expected:\n\n%-v\n", st, expected)
		}
	}

	if u, exist := structs["M_U1"]; !exist {
		t.Error("Expected a union named M_U1, but could not find it. Got", structs)
	} else {
		expected := StructDef{
			[]StructMemberDef{
				{"M_T", "ut", "", false, "false"},
				{"M_S", "us", "", false, "true"},
			},
			"bool",
			false,
		}
		if !reflect.DeepEqual(expected, u) {
			t.Errorf("Got:\n\n%-v\n\nbut expected:\n\n%-v\n", u, expected)
		}
	}

	if u, exist := structs["M_U2"]; !exist {
		t.Error("Expected a union named M_U2, but could not find it. Got", structs)
	} else {
		expected := StructDef{
			[]StructMemberDef{
				{"M_T", "ut", "", false, "MyEnum_One"},
				{"M_S", "us", "", false, "MyEnum_Two"},
			},
			"M_MyEnum",
			true,
		}
		if !reflect.DeepEqual(expected, u) {
			t.Errorf("Got:\n\n%-v\n\nbut expected:\n\n%-v\n", u, expected)
		}
	}
}
