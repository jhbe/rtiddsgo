package generate

import (
	"io"
	"rtiddsgo/parse"
	"text/template"
)

func TypeDefsFile(td []parse.TypeDef, packageName, rtiInstallDir, rtiLibDir, cFileName string, w io.Writer) error {
	typeDefs := struct {
		PackageName   string
		RtiInstallDir string
		RtiLibDir     string
		CFileName     string
		Unsafe        bool
		TypeDefs      []parse.TypeDef
	}{
		PackageName:   packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
		TypeDefs:      td,
	}

	for _, t := range td {
		if t.GoType == "string" {
			typeDefs.Unsafe = true
		}
	}

	return template.Must(template.New("typedefsFileTmpl").Funcs(template.FuncMap{
		"Type":     Type,
		"Store":    Store,
		"Retrieve": Retrieve,
	}).Parse(typedefsFileTmpl)).Execute(w, typeDefs)
}

var typedefsFileTmpl = `
{{range $typedef := .TypeDefs}}
type {{.GoName}} {{Type .GoType .SeqLen .ArrayDims}}

func (from {{.GoName}}) Store(to *C.{{.CName}}) {
	{{Store .GoType .CType "from" "(*to)" "to" .SeqLen .ArrayDims}}
}

func (to *{{.GoName}}) Retrieve(from C.{{.CName}}) {
	{{Retrieve .GoName .GoType .CType "(*to)" "from" .SeqLen .ArrayDims true}}
}

{{end}}
`
