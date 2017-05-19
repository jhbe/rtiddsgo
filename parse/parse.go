//go:generate /home/johan/go/bin/goyacc -o parser.go -p "Idl" parser.y
package parse

import (
	"io"
	"strings"
	"errors"
	"unicode"
)

// Used from within the goyacc-generated code, so has to be a global variable.
var parsingError string

// Parse parses the IDL file "idlFile" and extracts all structs/unions,
// enumerations and constants.
func Parse(idlFile io.Reader) (StructsDef, ConstsDef, EnumsDef, error) {
	IdlErrorVerbose = true

	// Create a lexer and start the goyacc parser. The result will be in the
	// variable TheSpecification.
	l := NewLexer(idlFile)
	IdlParse(l)
	if TheSpecification == nil {
		return StructsDef{}, ConstsDef{}, EnumsDef{}, errors.New(parsingError)
	}

	// Types and name in a module may be referenced as-is, but as Golang does
	// not have the concept of scope we will force all types to have a name
	// that is fully qualified, i.e. includes the module name(s).
	//
	// First change all enum specifications to be fully qualified (leave the actual
	// enum values alone):
	//
	// module M {
	//   enum E {    E => M_E
	//      One, Two
	//   };
	// };
	//
	fixEnumNames(TheSpecification)
	// ...then change all const names to be fully qualified and use fully
	// qualified type names...
	fixConstNames(TheSpecification)
	// ...and change all names of structs to be fully qualified and
	// make sure the types used within the structs are fully qualified.
	fixTypeNames(TheSpecification)
	// ...and finally change all uses of consts to use fully qualified names
	fixConstUsage(TheSpecification)

	return getStructsDef(TheSpecification), getConstsDef(TheSpecification), getEnumsDef(TheSpecification), nil
}

// fixEnumNames replaces all enum names with a fully qualified variant.
//
// module A {
//   enum B {    E => A_B
//      One, Two
//   };
// };
//
// The name of "B" will be replaced with "A_B".
//
func fixEnumNames(spec *Node) {
	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindEnum
	}, func(n *Node) {
		n.Name = n.FullPathName("_")
	})
}

// fixConstNames replaces all const names with a fully qualified variant.
//
// module A {
//   const long B = 5;  B => A_B
// };
//
func fixConstNames(spec *Node) {
	firstModuleName := getTopmostModuleName(spec)
	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindConst && !strings.HasPrefix(n.Name, firstModuleName)
	}, func(n *Node) {
		n.Name = n.FullPathName("_")
	})

}

// fixConstUsage replaces references to consts that
// and not fully qualified (because they reside in the same module as where
// they are defined), with fully qualified const name.
//
// module A {
//   const long B = 5;
//   struct D {
//      sequence<long, B> E;    B => A_B
//   };
// };
//
// The sequence is defined as "B" long, which is legal as "B" is part of the same
// module. But we want fully qualified names for all consts. So the sequence
// should be defined with A_B as the length.
//
func fixConstUsage(spec *Node) {
	firstModuleName := getTopmostModuleName(spec)

	// Fix all sequences where the length is a const. Make sure only references are
	// updated, not hardcoded numbers. Leave it unchecked if the const reference
	// already is qualified.
	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindMember && n.Length != "" && !unicode.IsDigit(rune(n.Length[0])) && !strings.HasPrefix(n.Length, firstModuleName)
	}, func(n *Node) {
		if modulePath := modulePath(n); modulePath != "" {
			n.Length = modulePath + "_" + n.Length
		}
	})
}

// fixTypesNames replaces types that are not built-in types (int, float, etc)
// and not fully qualified (because they reside in the same module as where
// they are defined), with fully qualified type names.
//
// module A {
//   struct B {
//     double C;
//   };
//   struct D {
//      B E;
//   };
// };
//
// Member "E" is defined as "B", which is legal as "B" is part of the same
// module. But we want fully qualified type names for all types that are not
// Go built-in types. So in effect any type defined in the IDL. So the type
// for "E" will be "B_E", not "B".
//
func fixTypeNames(spec *Node) {
	firstModuleName := getTopmostModuleName(spec)

	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindMember && !isBaseType(n.TypeName) && !strings.HasPrefix(n.TypeName, firstModuleName)
	}, func(n *Node) {
		if modulePath := modulePath(n); modulePath != "" {
			n.TypeName = modulePath + "_" + n.TypeName
		}
	})
}
