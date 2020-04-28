package generate

import (
	"rtiddsgo/parse"
	"io"
	"text/template"
)

func AllInOneFile(sd parse.StructDef, packageName string, w io.Writer) error {
	return template.Must(template.New("allInOneTmpl").Parse(allInOneTmpl)).Execute(w, struct {
		PackageName string
		GoName      string
	}{
		PackageName: packageName,
		GoName:      sd.GoName,
	})
}

var allInOneTmpl = `
//==========================================================================
// Writer
//==========================================================================

func NewWriter{{.GoName}}(domain int, topicName, participantQoSLibraryName, participantQosProfileName, writerQosLibraryName, writerQosProfileName string) {{.GoName}}WriterAllInOne {
	allInOne := {{.GoName}}WriterAllInOne{}
	var err error

	allInOne.P, err = rtiddsgo.New(domain, participantQoSLibraryName, participantQosProfileName)
	if err != nil {
		log.Fatal(err)
	}

	err = {{.GoName}}_RegisterType(allInOne.P)
	if err != nil {
		log.Fatal(err)
	}

	allInOne.T, err = allInOne.P.CreateTopic(topicName, {{.GoName}}_GetTypeName(), "", "")
	if err != nil {
		log.Fatal(err)
	}

	allInOne.Pub, err = allInOne.P.CreatePublisher("", "")
	if err != nil {
		log.Fatal(err)
	}

	allInOne.Dw, err = New{{.GoName}}DataWriter(allInOne.Pub, allInOne.T, writerQosLibraryName, writerQosProfileName)
	if err != nil {
		log.Fatal(err)
	}

	return allInOne
}

func (a *{{.GoName}}WriterAllInOne) Free() {
	a.Dw.Free()
	a.Pub.Free()
	a.T.Free()
	a.P.Free()
}

type {{.GoName}}WriterAllInOne struct {
	P   rtiddsgo.Participant
	T   rtiddsgo.Topic
	Pub rtiddsgo.Publisher
	Dw  {{.GoName}}DataWriter
}

//==========================================================================
// Reader
//==========================================================================

func NewReader{{.GoName}}(domain int, topicName, participantQoSLibraryName, participantQosProfileName, readerQosLibraryName, readerQosProfileName string, rx func(alive bool, data {{.GoName}})) {{.GoName}}ReaderAllInOne {
	allInOne := {{.GoName}}ReaderAllInOne{}
	var err error

	allInOne.P, err = rtiddsgo.New(domain, participantQoSLibraryName, participantQosProfileName)
	if err != nil {
		log.Fatal(err)
	}

	err = {{.GoName}}_RegisterType(allInOne.P)
	if err != nil {
		log.Fatal(err)
	}

	allInOne.T, err = allInOne.P.CreateTopic(topicName, {{.GoName}}_GetTypeName(), "", "")
	if err != nil {
		log.Fatal(err)
	}

	allInOne.Sub, err = allInOne.P.CreateSubscriber("", "")
	if err != nil {
		log.Fatal(err)
	}

	allInOne.Dr, err = New{{.GoName}}DataReader(allInOne.Sub, allInOne.T, "", "", rx)
	if err != nil {
		log.Fatal(err)
	}

	return allInOne
}

func (a *{{.GoName}}ReaderAllInOne) Free() {
	a.Dr.Free()
	a.Sub.Free()
	a.T.Free()
	a.P.Free()
}

type {{.GoName}}ReaderAllInOne struct {
	P   rtiddsgo.Participant
	T   rtiddsgo.Topic
	Sub rtiddsgo.Subscriber
	Dr  *{{.GoName}}DataReader
}
`
