package parse

import (
	"strings"
)

// ConstDef represents a single constant.
type ConstDef struct {
	Name  string // Fully qualified go name ("Module1_Module2_Foo") of the constant.
	Type  string // Fully qualified go name ("Module1_Module2_MyEnum") of the type of the constant. Can also be a basic type, such as "int16" or "string".
	Value string // The value assigned to the constant ("Module1_Module2_MyEnumTwo" or "hello" or "34").
}

// GetConstsDef traverses the result of the parsing of the XML file ("me") and extracts all constant definitions. All
// constant names will be fully qualified.
func (me ModuleElements) GetConstsDef() []ConstDef {
	var constants []ConstDef

	me.TraverseConstants("", "", func(cPath, goPath string, cd ConstDecl) {
		constants = append(constants, ConstDef{
			Name:  goNameOf(goPath, cd.Name),
			Type:  goPathTypeOf(goPath, cd.Type, cd.NonBasicTypeName),
			Value: goNameOf("", strings.Trim(cd.Value, "()")),
		})
	})
	return constants
}
