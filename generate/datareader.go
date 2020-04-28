package generate

import (
	"rtiddsgo/parse"
	"io"
	"text/template"
)

func DataReaderFile(sd parse.StructDef, packageName, rtiInstallDir, rtiLibDir, cFileName string, w io.Writer) error {
	return template.Must(template.New("dataReaderTmpl").Parse(dataReaderTmpl)).Execute(w, struct {
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

var dataReaderTmpl = `
type {{.GoName}}DataReader struct {
	dr  *rtiddsgo.DataReader
	cdr *C.{{.CName}}DataReader
	rx  func(alive bool, data {{.GoName}})
}

func New{{.GoName}}DataReader(sub rtiddsgo.Subscriber, t rtiddsgo.Topic, qosLibraryName, qosProfileName string, rx func(alive bool, data {{.GoName}})) (*{{.GoName}}DataReader, error) {
	messageDR := &{{.GoName}}DataReader{
		rx: rx,
	}

	var err error
	messageDR.dr, err = rtiddsgo.CreateDataReader(sub, t, qosLibraryName, qosProfileName, messageDR)
	if err != nil {
		return messageDR, err
	}

	messageDR.cdr = C.{{.CName}}DataReader_narrow((*C.DDS_DataReader)(messageDR.dr.GetUnsafe()))
	return messageDR, nil
}

func (dr *{{.GoName}}DataReader) Free() {
	dr.dr.Free()
}

func (dr *{{.GoName}}DataReader) DataAvailable() error {
	var rc C.DDS_ReturnCode_t
	for rc = C.DDS_RETCODE_OK; rc != C.DDS_RETCODE_NO_DATA; {
		var dataSeq C.struct_{{.CName}}Seq
		C.{{.CName}}Seq_initialize(&dataSeq)

		var sampleInfoSeq C.struct_DDS_SampleInfoSeq
		C.DDS_SampleInfoSeq_initialize(&sampleInfoSeq)

		rc = C.{{.CName}}DataReader_take(
			dr.cdr,
			&dataSeq,
			&sampleInfoSeq,
			C.DDS_LENGTH_UNLIMITED,
			C.DDS_ANY_SAMPLE_STATE,
			C.DDS_ANY_VIEW_STATE,
			C.DDS_ANY_INSTANCE_STATE)
		if rc != C.DDS_RETCODE_NO_DATA && rc != C.DDS_RETCODE_OK {
			return errors.New("_take failed")
		} else if rc != C.DDS_RETCODE_NO_DATA {
			for i := C.DDS_Long(0); i < C.{{.CName}}Seq_get_length(&dataSeq); i++ {
				si := C.DDS_SampleInfoSeq_get_reference(&sampleInfoSeq, i)
				var goData {{.GoName}}
				if si.instance_state != C.DDS_ALIVE_INSTANCE_STATE {
					dr.rx(false, goData)
				} else if si.valid_data == C.DDS_BOOLEAN_TRUE {
					rxData := C.{{.CName}}Seq_get(&dataSeq, i)

					goData.Retrieve(rxData)
					dr.rx(true, goData)
				}
			}
		}

		C.{{.CName}}DataReader_return_loan(dr.cdr, &dataSeq, &sampleInfoSeq)
	}
	return nil
}
`
