package generate

import "C"
import (
	"io"
	"rtiddsgo/parse"
	"text/template"
)

func UnionFile(ud parse.UnionDef, packageName, rtiInstallDir, rtiLibDir, cFileName string, w io.Writer) error {
	ts := struct {
		parse.UnionDef
		PackageName   string
		RtiInstallDir string
		RtiLibDir     string
		CFileName     string
		Unsafe        bool
	}{
		UnionDef:     ud,
		PackageName:   packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
	}

	for _, m := range ud.Members {
		if m.GoType == "string" {
			ts.Unsafe = true
		}
	}

	return template.Must(template.New("unionFileTmpl").Funcs(template.FuncMap{
		"Store":    Store,
		"Retrieve": Retrieve,
	}).Parse(unionFileTmpl)).Execute(w, ts)
}

var unionFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}

import (
{{if not .Nested}}	"errors"{{end}}
{{if not .Nested}}	"rtiddsgo"{{end}}
{{if .Unsafe}}    "unsafe"
{{end}})

` + flags + `
// #include <stdlib.h>
// #include <ndds/ndds_c.h>
// #include "{{.CFileName}}.h"
// #include "{{.CFileName}}Support.h"
import "C"

//=====================================================================
// Go type definition of the IDL type
//=====================================================================

type {{.GoName}} struct {
  _Discriminant {{.GoDiscriminantType}}
{{- range $member := .Members}}
  {{.GoName}} {{if .SeqLen}}[]{{end}}{{.GoType}} // For case {{.GoDiscriminatorValue}}{{end}}
}

//=====================================================================
// Functions for copying data from the Go type to the C type
//=====================================================================

func (from {{.GoName}}) Store(to *C.{{.CName}}) {
    {{Store .GoDiscriminantType .CDiscriminantType "from._Discriminant" "(*to)._d" "&(to._d)"}}

    switch from._Discriminant {
{{- range $member := .Members}}
    case {{.GoDiscriminatorValue}}:
{{- if .SeqLen}}
        C.{{.CType}}Seq_set_maximum(&to._u.{{.CName}}, C.DDS_Long({{.SeqLen}}))
        C.{{.CType}}Seq_set_length(&to._u.{{.CName}}, C.DDS_Long(len(from.{{.GoName}})))
        for index, _ := range from.{{.GoName}} {
            value := C.{{.CType}}Seq_get_reference(&to._u.{{.CName}}, C.DDS_Long(index))
            {{Store .GoType .CType (printf "from.%s[index]" .GoName) "*value" "value"}}
        }
{{- else}}
        {{Store .GoType .CType (printf "from.%s" .GoName) (printf "(*to)._u.%s" .CName) (printf "&(to._u.%s)" .CName)}}
{{- end -}}
{{end}}
    }
}

func (from {{.GoName}}) PostStore(to *C.{{.CName}}) {
}

//=====================================================================
// Functions for copying data from the C type to the Go type
//=====================================================================

func (to *{{.GoName}}) Retrieve(from C.{{.CName}}) {
	{{Retrieve .GoDiscriminantType .CDiscriminantType "(*to)._Discriminant" "from._d" false}}

    switch to._Discriminant {
{{- range $member := .Members}}
    case {{.GoDiscriminatorValue}}:
{{- if .SeqLen}}
	    (*to).{{.CName}} = make([]{{.GoType}}, C.{{.CType}}Seq_get_length(&from._u.{{.GoName}}))
	    for index, _ := range (*to).{{.GoName}} {
            value := C.{{.CType}}Seq_get(&from._u.{{.GoName}}, C.DDS_Long(index))
            {{Retrieve .GoName .GoType (printf "(*to).%s[index]" .GoName) "value" false}}
        }
{{- else}}
	    {{Retrieve .GoName .GoType (printf "(*to).%s" .CName) (printf "from._u.%s" .GoName) false}}
{{- end -}}
{{end}}
    }
}

func (to *{{.GoName}}) PostRetrieve(from C.{{.CName}}) {

}
{{if not .Nested}} 
//=====================================================================
// Type Support
//=====================================================================

func {{.GoName}}_GetTypeName() string {
  return C.GoString(C.{{.CName}}TypeSupport_get_type_name())
}

func {{.GoName}}_RegisterType(p rtiddsgo.Participant) error {
	rc := C.{{.CName}}TypeSupport_register_type(
		(*C.DDS_DomainParticipant)(p.GetUnsafe()),
		C.{{.CName}}TypeSupport_get_type_name())
	if rc != C.DDS_RETCODE_OK {
		return errors.New("Failed to register the type {{.GoName}}.")
	}
	return nil
}
{{end}}
`
