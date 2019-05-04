package parse

type UnionDef struct {
	GoName             string // Fully qualified go name ("Module1_Module2_Name") of the union.
	CName              string // Fully qualified C/DDS name of the union ("module1_module2_name"). Casing is unchanged.
	Nested             bool
	Members            []UnionMember
	GoDiscriminantType string
	CDiscriminantType  string
}

type UnionMember struct {
	GoName               string // Name of the member.
	CName                string // Name of the member as defined in C.
	GoType               string
	CType                string
	SeqLen               string
	ArrayDims            string
	GoDiscriminatorValue string // Fully qualified go value of the discriminator for this member.
}

func (me ModuleElements) GetUnionsDef() []UnionDef {
	var unions []UnionDef

	me.TraverseUnions("", "", func(cPath, goPath string, u UnionDecl) {
		unionDef := UnionDef{
			GoName:             goNameOf(goPath, u.Name),
			CName:              cNameOf(cPath, u.Name),
			Nested:             u.Nested == "true",
			GoDiscriminantType: goTypeOf(u.Discriminator.Type, u.Discriminator.NonBasicTypeName),
			CDiscriminantType:  ddsTypeOf(u.Discriminator.Type, u.Discriminator.NonBasicTypeName),
		}

		for _, cd := range u.CaseDecls {
			member := UnionMember{
				GoName:               goNameOf("", cd.Member.Name),
				CName:                cNameOf("", cd.Member.Name),
				GoType:               goTypeOf(cd.Member.Type, cd.Member.NonBasicTypeName),
				CType:                ddsTypeOf(cd.Member.Type, cd.Member.NonBasicTypeName),
				SeqLen:               cSeqLenOf(cd.Member.SequenceMaxLength),
				ArrayDims:            goArrayDimsOf(cd.Member.ArrayDimensions),
				GoDiscriminatorValue: goNameOf("", cd.CaseDiscriminator.Value),
			}
			unionDef.Members = append(unionDef.Members, member)
		}

		unions = append(unions, unionDef)
	})

	return unions
}
