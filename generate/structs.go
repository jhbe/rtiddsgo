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
	}{
		StructDef:     sd,
		PackageName:   packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
	}

	return template.Must(template.New("structsFileTmpl").Funcs(template.FuncMap{
		"Type":     Type,
		"Store":    Store,
		"Retrieve": Retrieve,
	}).Parse(structsFileTmpl)).Execute(w, ts)
}

var structsFileTmpl = `

//=====================================================================
// Go type definition of the IDL type
//=====================================================================

type {{.GoName}} struct {
  {{if .BaseType}}{{.BaseType}}{{end}}
{{- range $member := .Members}}
  {{.GoName}} {{Type .GoType .SeqLen .ArrayDims}}{{end}}
}

//=====================================================================
// Functions for copying data from the Go type to the C type
//=====================================================================

func (from {{.GoName}}) Store(to *C.{{.CName}}) {
{{if .BaseType}}    from.{{.BaseType}}.Store(&to.parent){{end}}
{{- range $member := .Members}}
    {{Store .GoType .CType (printf "from.%s" .GoName) (printf "(*to).%s" .CName) (printf "&(to.%s)" .CName) .SeqLen .ArrayDims}}
{{end}}
}

//=====================================================================
// Functions for copying data from the C type to the Go type
//=====================================================================

func (to *{{.GoName}}) Retrieve(from C.{{.CName}}) {
{{if .BaseType}}    to.{{.BaseType}}.Retrieve(from.parent){{end}}
{{- range $member := .Members}}
	{{Retrieve .GoName .GoType .CType (printf "to.%s" .GoName) (printf "from.%s" .CName) .SeqLen .ArrayDims false}}
{{end}}
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
