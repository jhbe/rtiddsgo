package parse

func (me ModuleElements) TraverseEnums(cPath, goPath string, do func(cPath, goPath string, ed EnumDecl)) {
	for _, e := range me.Enums {
		do(cPath, goPath, e)
	}
	for _, m := range me.Modules {
		m.TraverseEnums(cNameOf(cPath, m.Name), goNameOf(goPath, m.Name), do)
	}
}

func (me ModuleElements) TraverseConstants(cPath, goPath string, do func(cPath, goPath string, ed ConstDecl)) {
	for _, e := range me.Consts {
		do(cPath, goPath, e)
	}
	for _, m := range me.Modules {
		m.TraverseConstants(cNameOf(cPath, m.Name), goNameOf(goPath, m.Name), do)
	}
}

func (me ModuleElements) TraverseStructs(cPath, goPath string, do func(cPath, goPath string, s StructDecl)) {
	for _, s := range me.Structs {
		do(cPath, goPath, s)
	}
	for _, m := range me.Modules {
		m.TraverseStructs(cNameOf(cPath, m.Name), goNameOf(goPath, m.Name), do)
	}
}

func (me ModuleElements) TraverseUnions(cPath, goPath string, do func(cPath, goPath string, s UnionDecl)) {
	for _, u := range me.Unions {
		do(cPath, goPath, u)
	}
	for _, m := range me.Modules {
		m.TraverseUnions(cNameOf(cPath, m.Name), goNameOf(goPath, m.Name), do)
	}
}

func (me ModuleElements) TraverseTypedefs(cPath, goPath string, do func(cPath, goPath string, td TypeDefDecl)) {
	for _, td := range me.TypeDefs {
		do(cPath, goPath, td)
	}
	for _, m := range me.Modules {
		m.TraverseTypedefs(cNameOf(cPath, m.Name), goNameOf(goPath, m.Name), do)
	}
}
