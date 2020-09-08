package generate

import (
	"io"
	"rtiddsgo/parse"
	"text/template"
)

func HeaderFile(me parse.ModuleElements, packageName, rtiInstallDir, rtiLibDir, cFileName string, w io.Writer) error {
	header := tmplHeader{
		PackageName:   packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
	}

	// Work out which import we need in this file.
	for _, sd := range me.GetStructsDef() {
		if !sd.Nested {
			header.RtiDdsGo = true
			header.Errors = true
			header.Fmt = true
			header.Log = true
		}
		for _, m := range sd.Members {
			if m.GoType == "string" {
				header.Unsafe = true
			}
		}
	}
	for _, ud := range me.GetUnionsDef() {
		if !ud.Nested {
			header.Errors = true
			header.RtiDdsGo = true
		}
		for _, m := range ud.Members {
			if m.GoType == "string" {
				header.Unsafe = true
			}
		}
	}
	for _, t := range me.GetTypeDefs() {
		if t.GoType == "string" {
			header.Unsafe = true
		}
	}
	header.SomeImports = header.Errors || header.Fmt || header.Log || header.RtiDdsGo || header.Unsafe

	return template.Must(template.New("headerTmpl").Parse(headerTmpl)).Execute(w, header)
}

type tmplHeader struct {
	PackageName                                     string
	RtiInstallDir                                   string
	RtiLibDir                                       string
	CFileName                                       string
	SomeImports, Errors, Fmt, Log, RtiDdsGo, Unsafe bool
}

var headerTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}

{{if .SomeImports}}import (
{{if .Errors}}	"errors"
{{end}}{{if .Fmt}}	"fmt"
{{end}}{{if .Log}}	"log"
{{end}}{{if .RtiDdsGo}}	"rtiddsgo"
{{end}}{{if .Unsafe}}    "unsafe"
{{end}}){{end}}

` + flags + `
// #include <stdlib.h>
// #include <ndds/ndds_c.h>
// #include "{{.CFileName}}.h"
// #include "{{.CFileName}}Support.h"
// #include "{{.CFileName}}Plugin.h"
import "C"
`
