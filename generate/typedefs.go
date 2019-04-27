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
		Unsafe bool
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
		"Store": Store,
		"Retrieve": Retrieve,
	}).Parse(typedefsFileTmpl)).Execute(w, typeDefs)
}

var typedefsFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}

{{if .Unsafe}}import "unsafe"{{end}}

` + flags + `
// #include <stdlib.h>
// #include <ndds/ndds_c.h>
// #include "{{.CFileName}}.h"
// #include "{{.CFileName}}Support.h"
import "C"

{{range $typedef := .TypeDefs}}
type {{.GoName}} {{if .Len}}[]{{end}}{{.GoType}}

func (from {{.GoName}}) Store(to *C.{{.CName}}) {
{{if .Len}}
    C.{{.CType}}Seq_set_maximum(to, C.DDS_Long({{.Len}}))
	C.{{.CType}}Seq_set_length(to, C.DDS_Long(len(from)))
	for index, _ := range from {
		value := C.{{.CType}}Seq_get_reference(to, C.DDS_Long(index))
		{{Store .GoType .CType "from[index]" "*value" "value"}}
	}
{{else}}
	{{Store .GoType .CType "from" "(*to)" "to"}}
{{end}}
}

func (to *{{.GoName}}) Retrieve(from C.{{.CName}}) {
{{if .Len}}
	*to = make([]{{.GoType}}, C.{{.CType}}Seq_get_length(&from))
	for index, _ := range *to {
		value := C.{{.CType}}Seq_get(&from, C.DDS_Long(index))
		{{Retrieve .GoName .GoType "(*to)[index]" "value" false}}
	}
{{else}}
	{{Retrieve .GoName .GoType "(*to)" "from" true}}
{{end}}
}

{{end}}
`
