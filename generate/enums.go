package generate

import (
	"rtiddsgo/parse"
	"io"
	"text/template"
)

func EnumsFile(e parse.EnumsDef, packageName, rtiInstallDir, rtiLibDir, cFileName string, w io.Writer) error {
	allEnums := tmplEnums{
		PackageName: packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
		TE: make([]tmplEnum, len(e)),
	}

	for ix, ee := range e {
		allEnums.TE[ix].GoName = ee.GoName
		allEnums.TE[ix].CName = ee.CName
		allEnums.TE[ix].Members = make([]tmplEnumMember, len(ee.Enums))

		for jx, em := range ee.Enums {
			allEnums.TE[ix].Members[jx].Name = em
			allEnums.TE[ix].Members[jx].Value = ee.Values[em]
		}
	}

	return template.Must(template.New("enumFileTmpl").Parse(enumFileTmpl)).Execute(w, allEnums)
}

type tmplEnums struct {
	PackageName string
	RtiInstallDir      string
	RtiLibDir          string
	CFileName          string
	TE          []tmplEnum
}

// Definition of an enum suitable to a template.
type tmplEnum struct {
	GoName    string           // The fully qualified name of the enum. Will not be empty.
	CName string
	Members []tmplEnumMember // An array of the enum members.
}

type tmplEnumMember struct {
	Name  string // The fully qualified name of the member of the enum. Will not be empty.
	Value string // String representation of the integer value. Will not be empty.
}

var enumFileTmpl = `

{{range $enum := .TE}}
type {{$enum.GoName}} int

const (
{{- range $index, $member := $enum.Members}}
  {{$member.Name}} {{$enum.GoName}} = {{$member.Value}}{{end}}
)

func (from {{$enum.GoName}}) Store(to *C.{{.CName}}) {
    *to = C.{{$enum.CName}} (from)
}

func (to *{{$enum.GoName}}) Retrieve(from C.{{.CName}}) {
    *to = {{$enum.GoName}}(from)
}
{{end}}
`
