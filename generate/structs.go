package generate

import "C"
import (
	"io"
	"rtiddsgo/parse"
	"text/template"
)

func StructFile(sd parse.StructDef, packageName, rtiInstallDir, rtiLibDir, cFileName string, w io.Writer) error {
	ts := struct {
		parse.StructDef
		PackageName   string
		RtiInstallDir string
		RtiLibDir     string
		CFileName     string
		Unsafe        bool
	}{
		StructDef:     sd,
		PackageName:   packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
	}

	for _, m := range sd.Members {
		if m.GoType == "string" {
			ts.Unsafe = true
		}
	}

	return template.Must(template.New("structsFileTmpl").Funcs(template.FuncMap{
		"Store":    Store,
		"Retrieve": Retrieve,
	}).Parse(structsFileTmpl)).Execute(w, ts)
}

var structsFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

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
// #include "{{.CFileName}}Plugin.h"
import "C"

//=====================================================================
// Go type definition of the IDL type
//=====================================================================

type {{.GoName}} struct {
  {{if .BaseType}}{{.BaseType}}{{end}}
{{- range $member := .Members}}
  {{.GoName}} {{if .SeqLen}}[]{{end}}{{.GoType}}{{end}}
}

//=====================================================================
// Functions for copying data from the Go type to the C type
//=====================================================================

func (from {{.GoName}}) Store(to *C.{{.CName}}) {
{{if .BaseType}}    from.{{.BaseType}}.Store(&to.parent){{end}}
{{- range $member := .Members}}
{{- if .SeqLen}}
    C.{{.CType}}Seq_set_maximum(&to.{{.CName}}, C.DDS_Long({{.SeqLen}}))
	C.{{.CType}}Seq_set_length(&to.{{.CName}}, C.DDS_Long(len(from.{{.GoName}})))
	for index, _ := range from.{{.GoName}} {
		value := C.{{.CType}}Seq_get_reference(&to.{{.CName}}, C.DDS_Long(index))
		{{Store .GoType .CType (printf "from.%s[index]" .GoName) "*value" "value"}}
	}
{{- else}}
    {{Store .GoType .CType (printf "from.%s" .GoName) (printf "(*to).%s" .CName) (printf "&(to.%s)" .CName)}}
{{- end -}}
{{end}}
}

func (from {{.GoName}}) PostStore(to *C.{{.CName}}) {
}

//=====================================================================
// Functions for copying data from the C type to the Go type
//=====================================================================

func (to *{{.GoName}}) Retrieve(from C.{{.CName}}) {
{{if .BaseType}}    to.{{.BaseType}}.Retrieve(from.parent){{end}}
{{- range $member := .Members}}
{{- if .SeqLen}}
	(*to).{{.CName}} = make([]{{.GoType}}, C.{{.CType}}Seq_get_length(&from.{{.GoName}}))
	for index, _ := range (*to).{{.GoName}} {
		value := C.{{.CType}}Seq_get(&from.{{.GoName}}, C.DDS_Long(index))
		{{Retrieve .GoName .GoType (printf "(*to).%s[index]" .GoName) "value" false}}
	}
{{- else}}
	{{Retrieve .GoName .GoType (printf "(*to).%s" .CName) (printf "from.%s" .GoName) false}}
{{- end -}}
{{end}}
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

//=====================================================================
// Test Support
//=====================================================================

func (in {{.GoName}}) StoreAndRetrieve() {{.GoName}} {
	instance := C.{{.CName}}TypeSupport_create_data()
	in.Store(instance)

	var unpacked {{.GoName}}
	unpacked.Retrieve(*instance)

	return unpacked
}
{{end}}
`
