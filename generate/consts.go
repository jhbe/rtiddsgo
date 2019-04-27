package generate

import (
	"io"
	"rtiddsgo/parse"
	"text/template"
)

// Generates a golang file into "w" with the constants from "cd" with "packageName"
func ConstsFile(cd []parse.ConstDef, packageName string, w io.Writer) error {
	allConsts := tmplConsts{
		PackageName: packageName,
		Consts:      make([]parse.ConstDef, len(cd)),
	}

	for ix, c := range cd {
		allConsts.Consts[ix].Name = c.Name
		allConsts.Consts[ix].Type = c.Type
		allConsts.Consts[ix].Value = c.Value
	}

	return template.Must(template.New("constsFileTmpl").Parse(constsFileTmpl)).Execute(w, allConsts)
}

type tmplConsts struct {
	PackageName string
	Consts      []parse.ConstDef
}

var constsFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}
{{range $const := .Consts}}
const {{$const.Name}} {{$const.Type}} = {{$const.Value}}{{end}}
`
