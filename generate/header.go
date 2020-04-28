package generate

import (
	"io"
	"text/template"
)

func HeaderFile(packageName, rtiInstallDir, rtiLibDir, cFileName string, unsafe bool, w io.Writer) error {
	header := tmplHeader{
		PackageName:   packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
		Unsafe:        unsafe,
	}

	return template.Must(template.New("headerTmpl").Parse(headerTmpl)).Execute(w, header)
}

type tmplHeader struct {
	PackageName   string
	RtiInstallDir string
	RtiLibDir     string
	CFileName     string
	Unsafe        bool
}

var headerTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}

import (
	"errors"
	"fmt"
	"log"
	"rtiddsgo"
{{if .Unsafe}}    "unsafe"
{{end}})

` + flags + `
// #include <stdlib.h>
// #include <ndds/ndds_c.h>
// #include "{{.CFileName}}.h"
// #include "{{.CFileName}}Support.h"
// #include "{{.CFileName}}Plugin.h"
import "C"
`
