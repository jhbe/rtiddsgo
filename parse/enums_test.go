package parse

import "testing"

func TestGetEnumsDef(t *testing.T) {
	// module M {
	//   enum MyEnum {
	//     MyEnum_One,
	//     MyEnum_Two,
	//   };
	// };
	m := &Node{Name: "M", Kind: KindModule}
	e := &Node{Name: "MyEnum", Kind: KindEnum}
	one := &Node{Name: "MyEnum_One"}
	two := &Node{Name: "MyEnum_Two"}

	e.Add(one, two)
	m.Add(e)

	enums := getEnumsDef(m)
	if len(enums) != 1 {
		t.Error("Expected one enum.")
	} else {
		if enums[0].Name != "MyEnum" {
			t.Error("Expected the enum to be named MyEnum, but found", enums[0].Name)
		}
		if enums[0].Values[0] != "MyEnum_One" {
			t.Error("Expected the first enum member to be MyEnum_One, but found", enums[0].Values[0])
		}
		if enums[0].Values[1] != "MyEnum_Two" {
			t.Error("Expected the second enum member to be MyEnum_Two, but found", enums[0].Values[1])
		}
	}
}