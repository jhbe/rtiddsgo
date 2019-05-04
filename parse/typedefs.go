package parse

type TypeDef struct {
	GoName string // Fully qualified go name of the typedef type
	CName  string // Fully qualified C name of the typedef type

	GoType string // Fully qualified go name of the type the typedef references
	CType  string // Fully qualified C name of the type the typedef references

	SeqLen    string // Empty string means this is not a sequence.
	ArrayDims string
}

func (me ModuleElements) GetTypeDefs() []TypeDef {
	var typedefs []TypeDef
	me.TraverseTypedefs("", "", func(cPath, goPath string, td TypeDefDecl) {
		typedefs = append(typedefs, TypeDef{
			GoName:    goNameOf(goPath, td.Name),
			CName:     cNameOf(cPath, td.Name),
			GoType:    goTypeOf(td.Type, td.NonBasicTypeName),
			CType:     ddsTypeOf(td.Type, td.NonBasicTypeName),
			SeqLen:    cSeqLenOf(td.SequenceMaxLength),
			ArrayDims: goArrayDimsOf(td.ArrayDimensions),
		})
	})

	return typedefs
}
