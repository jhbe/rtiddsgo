package parse

import (
	"fmt"
	"io"
)

// StructsDef represents a collection of structs with their members and types.
type StructsDef map[string]StructDef

// StructDef represents a struct as a collection of member names (variables) and
// their types.
type StructDef struct {
	Members              []StructMemberDef
	DiscriminantType     string // Empty string if this struct is not a union.
	DiscriminantIsAnEnum bool
}

// StructMemberDef represents a single member in a struct with its type.
type StructMemberDef struct {
	TypeName, MemberName string
	SequenceLength       string // "": Not a sequence
	IsAnEnum             bool

	// If this member is part of a union (rather than struct), then this field
	// indicates the value of the discriminant for which this member should be
	// used.
	UnionValue string
}

// getStructsDef traverses the result of the IDL file parsing and returns all
// structs (types) except built-in types (int, float, bool, etc). Essentially, getStructsDef
// returns all struct and union definitions.
func getStructsDef(spec *Node) StructsDef {
	structs := make(StructsDef)

	// Find all structs and unions in the specification and add them to the Types map.
	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindType || n.Kind == KindUnionType
	}, func(n *Node) {
		var structDef StructDef

		// If this is a union, add record the discriminator type and whether it is an enum or not.
		if n.Kind == KindUnionType {
			typeName := getFullName(n.TypeName, n)
			structDef.DiscriminantType = typeName
			structDef.DiscriminantIsAnEnum = isAnEnumType(spec, typeName)
		}

		// Iterate of all members and add them to a StructDef
		for _, c := range n.Children() {
			structDef.Members = append(structDef.Members, StructMemberDef{
				c.TypeName,
				c.Name,
				c.Length,
				isAnEnumType(spec, c.TypeName),
				c.Value,
			})
		}

		// Add the StructDef to the structs map.
		structs[getFullName(n.Name, n)] = structDef
	})

	return structs
}

// Dump prints a representation of the structs to the io.Writer
func (structs StructsDef) Dump(out io.Writer) {
	for structName, structDef := range structs {
		fmt.Fprintln(out, structName, structDef.DiscriminantType)
		for _, member := range structDef.Members {
			switch member.SequenceLength {
			case "-1":
				fmt.Fprintf(out, "  %s []%s %v %s\n", member.TypeName, member.MemberName, member.IsAnEnum, member.UnionValue)
			case "":
				fmt.Fprintf(out, "  %s %s %v %s\n", member.TypeName, member.MemberName, member.IsAnEnum, member.UnionValue)
			default:
				fmt.Fprintf(out, "  %s [%s]%s %v %s\n", member.TypeName, member.SequenceLength, member.MemberName, member.IsAnEnum, member.UnionValue)
			}
		}
	}
}
