package generate

import (
	"io"
	"rtiddsgo/parse"
	"log"
	"text/template"
)

func CreateConstsFile(packageName string, constsDef parse.ConstsDef, outWriter io.Writer) error {
	cf := constFile{
		PackageName: packageName,
		Members: constsDef,
	}

	tmpl := template.Must(template.New("constFileTmpl").Parse(constFileTmpl))
	if err := tmpl.Execute(outWriter, cf); err != nil {
		log.Fatal(err)
	}
	return nil
}

type constFile struct {
	PackageName string
	Members     parse.ConstsDef
}

var constFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}
{{range $member := .Members}}
const {{$member.Name}} {{$member.Type}} = {{$member.Value}}
{{- end}}
`
