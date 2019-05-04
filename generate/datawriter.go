package generate

import (
	"rtiddsgo/parse"
	"io"
	"text/template"
)

func DataWriterFile(sd parse.StructDef, packageName, rtiInstallDir, rtiLibDir, cFileName string, w io.Writer) error {
	return template.Must(template.New("dataWriterFileTmpl").Parse(dataWriterFileTmpl)).Execute(w, struct {
		PackageName   string
		RtiInstallDir string
		RtiLibDir     string
		CFileName     string
		GoName        string
		CName         string
	}{
		PackageName:   packageName,
		RtiInstallDir: rtiInstallDir,
		RtiLibDir:     rtiLibDir,
		CFileName:     cFileName,
		GoName:        sd.GoName,
		CName:         sd.CName,
	})
}

var dataWriterFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}

import (
	"fmt"
	"rtiddsgo"
)

` + flags + `
// #include <stdlib.h>
// #include <ndds/ndds_c.h>
// #include "{{.CFileName}}.h"
// #include "{{.CFileName}}Support.h"
// #include "{{.CFileName}}Plugin.h"
import "C"

type {{.GoName}}DataWriter struct {
	dw  rtiddsgo.DataWriter
	cdw *C.{{.CName}}DataWriter
}

func New{{.GoName}}DataWriter(pub rtiddsgo.Publisher, t rtiddsgo.Topic, qosLibraryName, qosProfileName string) ({{.GoName}}DataWriter, error) {
	messageDW := {{.GoName}}DataWriter{}
	var err error
	messageDW.dw, err = rtiddsgo.CreateDataWriter(pub, t, qosLibraryName, qosProfileName)
	if err != nil {
		return messageDW, err
	}
	messageDW.cdw = C.{{.CName}}DataWriter_narrow((*C.DDS_DataWriter)(messageDW.dw.GetUnsafe()))
	return messageDW, nil
}

func (dw {{.GoName}}DataWriter) Free() {
	dw.dw.Free()
}

func (dw {{.GoName}}DataWriter) Write(m {{.GoName}}) error {
	instance := C.{{.CName}}TypeSupport_create_data()
	m.Store(instance)

	rc := C.{{.CName}}DataWriter_write(
		dw.cdw,
		instance,
		&C.DDS_HANDLE_NIL)
	defer func() {
		C.{{.CName}}TypeSupport_delete_data(instance)
	}()
	if rc != C.DDS_RETCODE_OK {
		return fmt.Errorf("Failed to write. Return code was %s", rtiddsgo.ReturnCodeToString(int(rc)))
	}
	return nil
}
`
