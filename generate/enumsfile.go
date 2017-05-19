package generate

import (
	"io"
	"rtiddsgo/parse"
	"log"
	"text/template"
)

func CreateEnumsFile(packageName string, enumsDef parse.EnumsDef, outWriter io.Writer) error {
	cf := enumFile{
		PackageName: packageName,
		Enums: enumsDef,
	}

	tmpl := template.Must(template.New("enumFileTmpl").Funcs(template.FuncMap{
		"ToGoName": toGoName,
	}).Parse(enumFileTmpl))
	if err := tmpl.Execute(outWriter, cf); err != nil {
		log.Fatal(err)
	}
	return nil
}

type enumFile struct {
	PackageName string
	Enums     parse.EnumsDef
}

var enumFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}
{{range $enum := .Enums}}
type {{$enum.Name | ToGoName}} int
const (
{{- range $index, $value := $enum.Values}}
	{{$value | ToGoName}}{{if eq $index 0}} {{$enum.Name | ToGoName}} = iota{{end}}{{end}}
)
{{end}}
`
