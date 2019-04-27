package parse

import (
	"encoding/xml"
	"fmt"
	"io"
)

func ReadXml(r io.Reader) (Types, error) {
	var t Types
	return t, xml.NewDecoder(r).Decode(&t)
}

type Types struct {
	ModuleElements
}

func (v Types) Dump(w io.Writer, indent int) {
	v.ModuleElements.Dump(w, -2)
}

type ModuleElements struct {
	Includes     []IncludeDecl   `xml:"include"`
	Modules      []ModuleDecl    `xml:"module"`
	Consts       []ConstDecl     `xml:"const"`
	Structs      []StructDecl    `xml:"struct"`
	ValueTypes   []ValueTypeDecl `xml:"valuetype"`
	Unions       []UnionDecl     `xml:"union"`
	TypeDefs     []TypeDefDecl   `xml:"typedef"`
	Enums        []EnumDecl      `xml:"enum"`
	ForwardDecls []ForwardDecl   `xml:"forward_dcl"`
}

func (v ModuleElements) Dump(w io.Writer, indent int) {
	for _, vv := range v.Includes {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.Modules {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.Consts {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.Enums {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.Structs {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.ValueTypes {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.Unions {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.TypeDefs {
		vv.Dump(w, indent+2)
	}
	for _, vv := range v.ForwardDecls {
		vv.Dump(w, indent+2)
	}
}

type IncludeDecl struct {
	FileName string `xml:"file,attr"`
}

func (v IncludeDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sInclude %s\n", indent, "", v.FileName)
}

type ModuleDecl struct {
	Name string `xml:"name,attr"`
	ModuleElements
}

func (v ModuleDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sModule %s\n", indent, "", v.Name)
	v.ModuleElements.Dump(w, indent)
}

type ConstDecl struct {
	Name             string `xml:"name,attr"`
	Type             string `xml:"type,attr"`  // Can be "nonBasic", in which case NonBasicTypeName holds the type name.
	Value            string `xml:"value,attr"` // Strings have double quotes. Non-basic types have brackets.
	NonBasicTypeName string `xml:"nonBasicTypeName"`
	ResolveName      string `xml:"resolveName"`
	StringMaxLength  string `xml:"stringMaxLength"`
}

func (v ConstDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sConst %s %s %s   %s %s %s\n", indent, "", v.Name, v.Type, v.Value, v.NonBasicTypeName, v.ResolveName, v.StringMaxLength)
}

type StructDecl struct {
	Name          string       `xml:"name,attr"`
	BaseType      string       `xml:"baseType,attr"`
	Members       []MemberDecl `xml:"member"`
	TopLevel      string       `xml:"topLevel,attr"`
	Extensibility string       `xml:"extensibility,attr"`
	ResolveName   string       `xml:"resolveName,attr"`
}

func (v StructDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sStruct %s %s\n", indent, "", v.Name, v.BaseType)
	for _, vv := range v.Members {
		vv.Dump(w, indent+2)
	}
}

type MemberDecl struct {
	Name              string `xml:"name,attr"`
	Type              string `xml:"type,attr"`
	NonBasicTypeName  string `xml:"nonBasicTypeName,attr"`
	Key               string `xml:"key,attr"`
	Pointer           string `xml:"pointer,attr"`
	BitField          string `xml:"bitField,attr"`
	StringMaxLength   string `xml:"stringMaxLength,attr"`
	SequenceMaxLength string `xml:"sequenceMaxLength,attr"`
	ArrayDimensions   string `xml:"arrayDimensions,attr"`
	ResolveName       string `xml:"resolveName,attr"`
	Visibility        string `xml:"visibility,attr"`
	Id                string `xml:"id,attr"`
	Optional          string `xml:"optional,attr"`
}

func (v MemberDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sMember %s %s %s\n", indent, "", v.Name, v.Type, v.NonBasicTypeName)
}

type ValueTypeDecl struct {
	Name          string       `xml:"name,attr"`
	TopLevel      string       `xml:"topLevel,attr"`
	BaseClass     string       `xml:"baseClass,attr"`
	TypeModifier  string       `xml:"typeModifier,attr"`
	Extensibility string       `xml:"extensibility,attr"`
	ResolveName   string       `xml:"resolveName,attr"`
	Members       []MemberDecl `xml:"member"`
	Consts        []ConstDecl  `xml:"const"`
	Unions        []UnionDecl  `xml:union`
}

func (vt ValueTypeDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sValueType %s %s %s\n", indent, "", vt.Name, vt.BaseClass, vt.ResolveName)
	for _, cd := range vt.Members {
		cd.Dump(w, indent+2)
	}
	for _, c := range vt.Consts {
		c.Dump(w, indent+2)
	}
	for _, u := range vt.Unions {
		u.Dump(w, indent+2)
	}

}

type UnionDecl struct {
	Discriminator Discriminator `xml:"discriminator"`
	CaseDecls     []CaseDecl    `xml:"case"`
	Name          string        `xml:"name,attr"`
	TopLevel      string        `xml:"topLevel,attr"`
	Extensibility string        `xml:"extensibility,attr"`
	ResolveName   string        `xml:"resolveName,attr"`
}

func (ud UnionDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sUnion %s %s\n", indent, "", ud.Name, ud.Extensibility)
	ud.Discriminator.Dump(w, indent+2)
	for _, cd := range ud.CaseDecls {
		cd.Dump(w, indent+2)
	}
}

type Discriminator struct {
	Type             string `xml:"type,attr"` // Can be "nonBasic", in which case NonBasicTypeName holds the type name.
	NonBasicTypeName string `xml:"nonBasicTypeName,attr"`
}

func (d Discriminator) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sDiscriminator %s %s\n", indent, "", d.Type, d.NonBasicTypeName)
}

type CaseDecl struct {
	CaseDiscriminator CaseDiscriminator `xml:"caseDiscriminator"`
	Member            MemberDecl        `xml:"member"`
}

func (cd CaseDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sCase\n", indent, "")
	cd.CaseDiscriminator.Dump(w, indent+2)
	cd.Member.Dump(w, indent+2)
}

type CaseDiscriminator struct {
	Value string `xml:"value,attr"`
}

func (cd CaseDiscriminator) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sCaseDiscriminator %s\n", indent, "", cd.Value)
}

type TypeDefDecl struct {
	Name              string `xml:"name,attr"`
	Type              string `xml:"type,attr"`
	NonBasicTypeName  string `xml:"nonBasicTypeName,attr"`
	TopLevel          string `xml:"topLevel,attr"`
	StringMaxLength   string `xml:"stringMaxLength,attr"`
	SequenceMaxLength string `xml:"sequenceMaxLength,attr"`
	ArrayDimensions   string `xml:"arrayDimensions,attr"`
	Pointer           string `xml:"pointer,attr"`
	ResolveName       string `xml:"resolveName,attr"`
}

func (td TypeDefDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sTypeDef %s %s %s etc...\n", indent, "", td.Name, td.Type, td.NonBasicTypeName)
}

type EnumDecl struct {
	Name          string           `xml:"name,attr"`
	Extensibility string           `xml:"extensibility,attr"`
	Enumerators   []EnumeratorDecl `xml:"enumerator"`
}

func (ed EnumDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sEnum %s %s\n", indent, "", ed.Name, ed.Extensibility)
	for _, e := range ed.Enumerators {
		e.Dump(w, indent+2)
	}
}

type EnumeratorDecl struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (ed EnumeratorDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sEnumerator %s %s\n", indent, "", ed.Name, ed.Value)
}

type ForwardDecl struct {
	Name string `xml:"name,attr"`
	Kind string `xml:"kind,attr"`
}

func (fd ForwardDecl) Dump(w io.Writer, indent int) {
	fmt.Fprintf(w, "%*sForwardDecl %s %s\n", indent, "", fd.Name, fd.Kind)
}
