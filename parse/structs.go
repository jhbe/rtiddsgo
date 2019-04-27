package parse

// StructDef represents a struct.
type StructDef struct {
	GoName               string // Fully qualified go name ("Module1_Module2_Name") of the struct.
	CName                string // Fully qualified C/DDS name of struct ("module1_module2_name"). Casing is unchanged.
	BaseType             string
	Members              []StructMember
}

type StructMember struct {
	GoName               string // Name of the member.
	CName                string // Name of the member as defined in C.
	GoType               string
	CType                string
	SeqLen               string
}

// GetStructsDef traverses the result of the parsing of the XML file ("me")
// and extracts all struct definitions. All struct names will be fully qualified.
func (me ModuleElements) GetStructsDef() []StructDef {
	structs := make([]StructDef, 0)

	me.TraverseStructs("", "", func(cPath, goPath string, s StructDecl) {
		structDef := StructDef{
			GoName:   goNameOf(goPath, s.Name),
			CName:    cNameOf(cPath, s.Name),
			BaseType: goTypeOf("nonBasic", s.BaseType),
		}

		for _, m := range s.Members {
			member := StructMember{
				GoName: goNameOf("", m.Name),
				CName:  cNameOf("", m.Name),
				GoType: goTypeOf(m.Type, m.NonBasicTypeName),
				CType:  ddsTypeOf(m.Type, m.NonBasicTypeName),
				SeqLen: cSeqLenOf(m.SequenceMaxLength),
			}
			structDef.Members = append(structDef.Members, member)
		}

		structs = append(structs, structDef)
	})
	return structs
}
