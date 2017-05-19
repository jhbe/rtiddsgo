package generate

import (
	"bytes"
	"go/parser"
	"go/token"
	"reflect"
	"rtiddsgo/parse"
	"testing"
)

func TestSimpleStruct(t *testing.T) {
	members := parse.StructDef{
		[]parse.StructMemberDef{
			{"float64", "theFloat", "", false, ""},
			{"string", "TheString", "", false, ""},
			{"int32", "theArray", "5", false, ""},
		}, "", false,
	}
	outBuffer := bytes.Buffer{}
	err := CreateStructFile("mypackagename", ".", "mytype", "/opt/rti_stuff", "MyType", members, &outBuffer)
	if err != nil {
		t.Error(err)
	}

	if _, err := parser.ParseFile(token.NewFileSet(), "foo.go", outBuffer.String(), 0); err != nil {
		t.Error(err, "\n\n", withLineNumbers(outBuffer.String()))
	}
}

func TestSequence(t *testing.T) {
	members := parse.StructDef{
		[]parse.StructMemberDef{
			{"MyModule_MyEnum", "E", "5", false, ""},
		}, "", false,
	}

	outBuffer := bytes.Buffer{}
	err := CreateStructFile("mypackagename", "./", "mytype", "", "MyType", members, &outBuffer)
	if err != nil {
		t.Error(err)
	}

	if _, err := parser.ParseFile(token.NewFileSet(), "foo.go", outBuffer.String(), 0); err != nil {
		t.Error(err, "\n\n", withLineNumbers(outBuffer.String()))
	}
}

func TestEnum(t *testing.T) {
	members := parse.StructDef{
		[]parse.StructMemberDef{
			{"MyModule_MyEnum", "E", "", true, ""},
		}, "", false,
	}

	outBuffer := bytes.Buffer{}
	err := CreateStructFile("mypackagename", "./", "mytype", "", "MyType", members, &outBuffer)
	if err != nil {
		t.Error(err)
	}

	if _, err := parser.ParseFile(token.NewFileSet(), "foo.go", outBuffer.String(), 0); err != nil {
		t.Error(err, "\n\n", withLineNumbers(outBuffer.String()))
	}
}

func TestUnion(t *testing.T) {
	members := parse.StructDef{
		[]parse.StructMemberDef{
			{"bool", "theBool", "", false, "MyEnum_One"},
			{"int32", "theInt", "", false, "MyEnum_Two"},
		}, "MyEnum", true,
	}

	outBuffer := bytes.Buffer{}
	err := CreateStructFile("mypackagename", "./", "mytype", "", "MyUnion", members, &outBuffer)
	if err != nil {
		t.Error(err)
	}

	if _, err := parser.ParseFile(token.NewFileSet(), "foo.go", outBuffer.String(), 0); err != nil {
		t.Error(err, "\n\n", withLineNumbers(outBuffer.String()))
	}
}

func TestToGoName(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"", ""},
		{"a", "A"},
		{" a", " A"},
		{"a ", "A "},
		{" a ", " A "},
		{"com", "Com"},
		{"com_this_that", "Com_This_That"},
		{"com this that", "Com This That"},
		{"theThing", "TheThing"},
		{"thething", "Thething"},
		{"int", "int"},
		{"int32", "int32"},
		{"float32", "float32"},
		{"float64", "float64"},
		{"string", "string"},
		{"bool", "bool"},
	}
	for index, c := range cases {
		s := toGoName(c.in)
		if s != c.out {
			t.Error("Case", index, ", expected", c.out, ", but got", s)
		}
	}
}

func TestIsSeparator(t *testing.T) {
	cases := []struct {
		in  rune
		out bool
	}{
		{'a', false},
		{'A', false},
		{' ', true},
		{'_', true},
		{';', false},
		{'0', false},
		{'-', false},
	}
	for index, c := range cases {
		isSep := isSeparator(c.in)
		if isSep && !c.out {
			t.Error("Case", index, ", expected is to not be a separator, but it was.")
		}
		if !isSep && c.out {
			t.Error("Case", index, ", expected is to be a separator, but it was not.")
		}
	}
}

func TestMembersOf(t *testing.T) {
	cases := []struct {
		in  parse.StructDef
		out []memberDef
	}{
		{in: parse.StructDef{
			Members: []parse.StructMemberDef{
				{TypeName: "bool", MemberName: "theBool", SequenceLength: "", IsAnEnum: false, UnionValue: ""},
				{TypeName: "int32", MemberName: "theLongs", SequenceLength: "5", IsAnEnum: false, UnionValue: ""},
				{TypeName: "M_MyEnum", MemberName: "theEnum", SequenceLength: "", IsAnEnum: true, UnionValue: ""},
			}, DiscriminantType: "", DiscriminantIsAnEnum: false,
		},
			out: []memberDef{
				{
					GoComment:      "",
					IsArray:        false,
					ArrayLength:    "",
					IsUnionMember:  false,
					IsEnum:         false,
					CTypeName:      "DDS_Boolean",
					GoTypeName:     "bool",
					GoFullTypeName: "bool",
					CMemberName:    "theBool",
					GoMemberName:   "TheBool",
					CFrom:          "instance.theBool",
					GoFrom:         "data.TheBool",
					CTo:            "instance.theBool",
					GoTo:           "data.TheBool",
					CFromValue:     "instance.theBool",
					GoFromValue:    "data.TheBool",
					CToValue:       "instance.theBool",
					GoToValue:      "data.TheBool",
					UnionValue:     "",
					SeqType:        "",
				},
				{
					GoComment:      "// Max length is 5 ",
					IsArray:        true,
					ArrayLength:    "5",
					IsUnionMember:  false,
					IsEnum:         false,
					CTypeName:      "DDS_Long",
					GoTypeName:     "int32",
					GoFullTypeName: "[]int32",
					CMemberName:    "theLongs",
					GoMemberName:   "TheLongs",
					CFrom:          "instance.theLongs",
					GoFrom:         "data.TheLongs",
					CTo:            "instance.theLongs",
					GoTo:           "data.TheLongs",
					CFromValue:     "*value",
					GoFromValue:    "data.TheLongs[index]",
					CToValue:       "*value",
					GoToValue:      "data.TheLongs[index]",
					UnionValue:     "",
					SeqType:        "DDS_Long",
				},
				{
					GoComment:      "",
					IsArray:        false,
					ArrayLength:    "",
					IsUnionMember:  false,
					IsEnum:         true,
					CTypeName:      "M_MyEnum",
					GoTypeName:     "M_MyEnum",
					GoFullTypeName: "M_MyEnum",
					CMemberName:    "theEnum",
					GoMemberName:   "TheEnum",
					CFrom:          "instance.theEnum",
					GoFrom:         "data.TheEnum",
					CTo:            "instance.theEnum",
					GoTo:           "data.TheEnum",
					CFromValue:     "instance.theEnum",
					GoFromValue:    "data.TheEnum",
					CToValue:       "instance.theEnum",
					GoToValue:      "data.TheEnum",
					UnionValue:     "",
					SeqType:        "",
				},
			},
		},
	}

	for index, c := range cases {
		members := membersOf(c.in)
		if !reflect.DeepEqual(c.out, members) {
			t.Errorf("Case %d: Expected: \n%-v\n, but got:\n%-v\n", index, c.out, members)
		}
	}
}
