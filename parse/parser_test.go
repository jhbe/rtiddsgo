package parse

import (
	"bytes"
	"strings"
	"testing"
)

func TestBasicType(t *testing.T) {
	s := `
struct A {
	long B;
};
`
	IdlErrorVerbose = true
	l := NewLexer(strings.NewReader(s))
	IdlParse(l)
	if TheSpecification == nil {
		t.Error(parsingError)
	} else {
		if b := TheSpecification.Child("A").Child("B"); b == nil {
			t.Error("Expected to find \"fg\".")
		} else if b.Name != "B" {
			t.Error("Expected the name of B to be B, but got", b.Name)
		} else if b.TypeName != "int32" {
			t.Error("Expected the typename of B to be int32, but got", b.TypeName, " B was", b)
		}
	}
}

func TestConst(t *testing.T) {
	s := `
module com {
	module MyModule {
    const string THE_CONST_STRING = "Blah";
    const long MAX_LENGTH = 5;

		struct My {
    	sequence<string, MAX_LENGTH> theString;
    };
	};
};
`

	IdlErrorVerbose = true
	l := NewLexer(strings.NewReader(s))
	IdlParse(l)
	if TheSpecification == nil {
		t.Error(parsingError)
	} else if c := TheSpecification.Child("com").Child("MyModule").Child("THE_CONST_STRING"); c == nil {
		t.Error("Expected to find com_MyModule_THE_CONST_STRING")
	} else {
		if c.Name != "THE_CONST_STRING" {
			t.Error("Expected the name of MyModule to be THE_CONST_STRING, but it was", c.Name)
		}
		if c.Kind != KindConst {
			t.Error("Expected the kind of MyModule to be KindConst, but it was", c.Kind)
		}
		if c.Value != "\"Blah\"" {
			t.Error("Expected the value to be \"Blah\", but it was", c.Value)
		}
	}

	if c := TheSpecification.Child("com").Child("MyModule").Child("MAX_LENGTH"); c == nil {
		t.Error("Expected to find com_MyModule_MAX_LENGTH")
	} else {
		if c.Name != "MAX_LENGTH" {
			t.Error("Expected the name of MyModule to be MAX_LENGTH, but it was", c.Name)
		}
		if c.Kind != KindConst {
			t.Error("Expected the kind of MyModule to be KindConst, but it was", c.Kind)
		}
		if c.Value != "5" {
			t.Error("Expected the kind of MyModule to be 5, but it was", c.Value)
		}
	}

	if c := TheSpecification.Child("com").Child("MyModule").Child("My").Child("theString"); c == nil {
		t.Error("Expected to find com_MyModule_My_theString")
	} else {
		if c.Name != "theString" {
			t.Error("Expected the name to be theString, but it was", c.Name)
		}
		if c.Kind != KindBaseMember {
			t.Error("Expected the kind to be KindSequence, but it was", c.Kind)
		}
		if c.TypeName != "string" {
			t.Error("Expected the type to be string, but it was", c.Value)
		}
		if c.Length != "MAX_LENGTH" {
			t.Error("Expected the length to be MAX_LENGTH, but it was", c.Length)
		}
	}
}

func TestEnum(t *testing.T) {
	s := `
module com {
	module MyModule {
    enum MyEnum {
      MyEnum_One,
      MyEnum_Two
    };
	};
};
`

	IdlErrorVerbose = true
	l := NewLexer(strings.NewReader(s))
	IdlParse(l)
	if TheSpecification == nil {
		t.Error(parsingError)
	} else {
		if c := TheSpecification.Child("com").Child("MyModule").Child("MyEnum"); c == nil {
			t.Error("Expected to find com_MyModule_MyEnum")
		} else {
			if c.Name != "MyEnum" {
				t.Error("Expected the name to be MyEnum, but it was", c.Name)
			}
			if c.Kind != KindEnum {
				t.Error("Expected the kind to be KindEnum, but it was", c.Kind)
			}
			if len(c.Children()) != 2 {
				t.Error("Expected two children.")
			}
			if one := c.Child("MyEnum_One"); one == nil {
				t.Error("Expected to find a child named MyEnum_One")
			} else {
				if one.Name != "MyEnum_One" {
					t.Error("Expected the name of the first child to be MyEnum_One.")
				}
			}
			if two := c.Child("MyEnum_Two"); two == nil {
				t.Error("Expected to find a child named MyEnum_Two")
			} else {
				if two.Name != "MyEnum_Two" {
					t.Error("Expected the name of the second child to be MyEnum_Two.")
				}
			}
		}
	}
}

func TestStructs(t *testing.T) {
	s := `
module com {
	module MyModule {
		struct MyError {
			long id;
		};
		struct MyType {
			long id_one, id_two;
			string text;
			sequence<com::MyModule::MyError, 2> theSeq;
		};
		struct MyMessage {
			com::MyModule::MyType theType;
			MyError theError;
		};
	};
};
`

	expected := &Node{}
	com := &Node{Name: "com", Kind: KindModule}
	myModule := &Node{Name: "MyModule", Kind: KindModule}

	myError := &Node{Name: "MyError", Kind: KindType}
	id := &Node{Name: "id", Kind: KindBaseMember, TypeName: "int32"}

	myType := &Node{Name: "MyType", Kind: KindType}
	id_one := &Node{Name: "id_one", Kind: KindBaseMember, TypeName: "int32"}
	id_two := &Node{Name: "id_two", Kind: KindBaseMember, TypeName: "int32"}
	text := &Node{Name: "text", Kind: KindBaseMember, TypeName: "string"}
	theSeq := &Node{Name: "theSeq", Kind: KindMember, TypeName: "com_MyModule_MyError", Length: "2"}

	myMessage := &Node{Name: "MyMessage", Kind: KindType}
	theType := &Node{Name: "theType", Kind: KindMember, TypeName: "com_MyModule_MyType"}
	theError := &Node{Name: "theError", Kind: KindMember, TypeName: "MyError"}

	myMessage.Add(theType, theError)
	myType.Add(id_one, id_two, text, theSeq)
	myError.Add(id)
	myModule.Add(myError, myType, myMessage)
	com.Add(myModule)
	expected.Add(com)

	IdlErrorVerbose = true
	l := NewLexer(strings.NewReader(s))
	IdlParse(l)
	if TheSpecification == nil {
		t.Error(parsingError)
	}
	if !expected.Equal(TheSpecification) {
		actualBuffer := bytes.Buffer{}
		TheSpecification.Dump(&actualBuffer)

		expectedBuffer := bytes.Buffer{}
		expected.Dump(&expectedBuffer)

		t.Errorf("The specification:\n%s\n...  did not match the expected:\n%s\n", actualBuffer.String(), expectedBuffer.String())
	}
}

func TestSequenceMember(t *testing.T) {
	s := `
const long C = 45;
enum E {
	E1
};
struct A {
	sequence<E, 22> B;
	sequence<E, C> F;
};
`
	IdlErrorVerbose = true
	l := NewLexer(strings.NewReader(s))
	IdlParse(l)
	if TheSpecification == nil {
		t.Error(parsingError)
	} else {
		if b := TheSpecification.Child("A").Child("B"); b == nil {
			t.Error("Expected to find \"b\".")
		} else if b.Name != "B" {
			t.Error("Expected the Name of B to be B but got", b.Name)
		} else if b.Kind != KindMember {
			t.Error("Expected the Kind of B to be KindMember, but got", b.Kind)
		} else if b.Length != "22" {
			t.Error("Expected the Length of B to be 22, but got ", b.Length)
		} else if b.TypeName != "E" {
			t.Error("Expected the TypeName of B to be E, but got ", b.TypeName)
		}

		if f := TheSpecification.Child("A").Child("F"); f == nil {
			t.Error("Expected to find \"f\".")
		} else if f.Name != "F" {
			t.Error("Expected the Name of F to be F, but got", f.Name)
		} else if f.Kind != KindMember {
			t.Error("Expected the Kind of F to be KindMember, but got", f.Kind)
		} else if f.Length != "C" {
			t.Error("Expected the Length of F to be C, but got ", f.Length)
		} else if f.TypeName != "E" {
			t.Error("Expected the TypeName of F to be E, but got ", f.TypeName)
		}
	}
}

func TestUnion(t *testing.T) {
	s := `
union A switch (MyEnum) {
	case MyEnum_One: long C;
	case MyEnum_Two: double D;
};
`
	expected := &Node{}
	u := &Node{Name: "A", Kind:KindUnionType, TypeName: "MyEnum"}
	c := &Node{Name: "C", Kind:KindBaseMember, TypeName: "int32", Value: "MyEnum_One"}
	d := &Node{Name: "D", Kind:KindBaseMember, TypeName: "float64", Value: "MyEnum_Two"}
	u.Add(c, d)
	expected.Add(u)

	IdlErrorVerbose = true
	l := NewLexer(strings.NewReader(s))
	IdlParse(l)
	if TheSpecification == nil {
		t.Error(parsingError)
	} else if !expected.Equal(TheSpecification){
		actualBuffer := bytes.Buffer{}
		TheSpecification.Dump(&actualBuffer)

		expectedBuffer := bytes.Buffer{}
		expected.Dump(&expectedBuffer)

		t.Errorf("The specification:\n%s\n...  did not match the expected:\n%s\n", actualBuffer.String(), expectedBuffer.String())
	}
}