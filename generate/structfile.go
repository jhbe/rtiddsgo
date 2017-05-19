package generate

import (
	"fmt"
	"io"
	"log"
	"rtiddsgo/parse"
	"strings"
	"text/template"
	"unicode"
)

// CreateStructFile creates golang source code and sends it to "outWriter"
// implementing the DDS type-specific functionality for the struct/union
// described in "structName" and "StructDef".
func CreateStructFile(packageName, c_files_path string, c_file_name string, rti_path string, structName string, structDef parse.StructDef, outWriter io.Writer) error {
	fd := fileDef{
		CPath:           c_files_path,
		CFileName:       c_file_name,
		CRtiIncludePath: rti_path,

		PackageName:  packageName,
		CStructName:  structName,
		GoStructName: toGoName(structName),

		Members: membersOf(structDef),
	}

	if structDef.DiscriminantType != "" {
		fd.IsUnion = true
		fd.DiscriminantMember = memberOf(parse.StructMemberDef{
			TypeName:   structDef.DiscriminantType,
			MemberName: structName + "_d",
			IsAnEnum:   structDef.DiscriminantIsAnEnum,
		})
		fd.DiscriminantMember.CFrom = "instance._d"
		fd.DiscriminantMember.CFromValue = fd.DiscriminantMember.CFrom
		fd.DiscriminantMember.CTo = "instance._d"
		fd.DiscriminantMember.CToValue = fd.DiscriminantMember.CTo
	}

	tmpl, err := template.New("structFileTmpl").Parse(structFileTmpl)
	if err != nil {
		fmt.Println(withLineNumbers(structFileTmpl))
		log.Fatal(err)
	}
	if err := tmpl.Execute(outWriter, fd); err != nil {
		fmt.Println(withLineNumbers(structFileTmpl))
		log.Fatal(err)
	}

	return nil
}

// fileDef contains all information required to generate a single file for a struct/union.
type fileDef struct {
	CPath           string // Full path to the rtiddsgen-generated file for this struct/union. Must not end with a forward slash.
	CFileName       string // The name of the file excluding postfix (.h, Support.h, Support) for the files in CPath. Case sensitive.
	CRtiIncludePath string // The full path to the RTI DDS installation directory. For example: "/opt/rti_connext_dds-5.2.3".

	PackageName  string // The case sensitive package name. Should be all lowercase, one word.
	CStructName  string // The name of the C struct/union.
	GoStructName string // The corresponding Go name for the struct/union.

	IsUnion            bool
	DiscriminantMember memberDef // If this is a union, then this variable holds information about the distriminant variable.

	Members []memberDef
}

// memberDef defines a variable in a struct/union. Essentially the same as
// parse.StructMemberDef, but easier to use in the template.
type memberDef struct {
	GoComment     string // For union...  Max length 5...
	IsArray       bool
	ArrayLength   string
	IsUnionMember bool
	IsEnum        bool

	CTypeName      string // Type used in the C specification: DDS_Boolean, DDS_Long, M_MyType, etc...
	GoTypeName     string // Fully qualified go type name. Used for determining the type when storing, etc.
	GoFullTypeName string // Like GoTypeName, but possibly prepended with "[]" amd used as the type in the struct.
	CMemberName    string
	GoMemberName   string

	CFrom  string // instance.blah
	GoFrom string // data.Blah
	CTo    string // instance.blah or instance._u.blah
	GoTo   string // data.Blah

	CFromValue  string // instance.blah or *value
	GoFromValue string // Same as GoFrom or "*value" if this is an array
	CToValue    string // Same as CTo or "*value" if this is an array
	GoToValue   string // data.Blah or *value

	UnionValue string

	SeqType string // DDS_Xxxx
}

// membersOf creates an array of memberDef based on the supplied structDef. It's
// the same information, but memberDef is tailored to make it easy to use it in
// the template.
func membersOf(structDef parse.StructDef) []memberDef {
	var members []memberDef
	for _, smd := range structDef.Members {
		members = append(members, memberOf(smd))
	}
	return members
}

// memberOf creates a memberDef based on the supplied structMemberDef. It's
// the same information, but memberDef is tailored to make it easy to use it in
// the template.
func memberOf(structMemberDef parse.StructMemberDef) memberDef {
	memberDef := memberDef{}

	// Append a comment describing any array length and/or whether the member is only used for a certain union.
	if structMemberDef.SequenceLength != "" {
		memberDef.GoComment += "Max length is " + structMemberDef.SequenceLength + " "
	}
	if structMemberDef.UnionValue != "" {
		memberDef.GoComment += "For when union discriminant is " + structMemberDef.UnionValue + " "
	}
	if len(memberDef.GoComment) != 0 {
		memberDef.GoComment = "// " + memberDef.GoComment
	}

	memberDef.CTypeName = seqTypeOf(structMemberDef.TypeName)
	memberDef.GoTypeName = toGoName(structMemberDef.TypeName)
	memberDef.GoFullTypeName = memberDef.GoTypeName
	memberDef.CMemberName = structMemberDef.MemberName
	memberDef.GoMemberName = toGoName(structMemberDef.MemberName)

	memberDef.CFrom = "instance." + structMemberDef.MemberName
	memberDef.GoFrom = "data." + toGoName(structMemberDef.MemberName)
	memberDef.CTo = "instance." + structMemberDef.MemberName
	memberDef.GoTo = "data." + toGoName(structMemberDef.MemberName)

	// Default to the same as above. These will change if the type is a sequence.
	memberDef.CFromValue = memberDef.CFrom
	memberDef.GoFromValue = memberDef.GoFrom
	memberDef.CToValue = memberDef.CTo
	memberDef.GoToValue = memberDef.GoTo

	// Is this member part of a union?
	if structMemberDef.UnionValue != "" {
		memberDef.IsUnionMember = true
		memberDef.UnionValue = structMemberDef.UnionValue

		memberDef.CMemberName = "_u." + structMemberDef.MemberName

		memberDef.CFrom = "instance._u." + structMemberDef.MemberName
		memberDef.CFromValue = "instance._u." + structMemberDef.MemberName
		memberDef.CTo = "instance._u." + structMemberDef.MemberName
		memberDef.CToValue = "instance._u." + structMemberDef.MemberName
	}

	// Is this member an array (sequence)?.
	if structMemberDef.SequenceLength != "" {
		memberDef.IsArray = true
		memberDef.ArrayLength = structMemberDef.SequenceLength
		memberDef.SeqType = seqTypeOf(structMemberDef.TypeName)

		memberDef.GoFullTypeName = "[]" + memberDef.GoTypeName

		memberDef.CFromValue = "*value"
		memberDef.GoFromValue = memberDef.GoFromValue + "[index]"
		memberDef.CToValue = "*value"
		memberDef.GoToValue = memberDef.GoToValue + "[index]"
	}

	// Is this an enum?
	if structMemberDef.IsAnEnum {
		memberDef.IsEnum = true
	}

	return memberDef
}

// toGoName returns a copy of the string s with all unicode characters that
// begin words mapped to their title (upper) case. Words are separated by
// whitespace or an underscore.
//
// Except if the string is a golang builtin type, in which case the string
// is unchanged.
//
// Examples:
// com_this_that => Com_This_That
// this          => This
// thething      => Thething
// theThing      => TheThing
// string        => string
//
func toGoName(s string) string {
	switch s {
	case "int", "int16", "int32",
		"uint", "uint16", "uint32",
		"float32", "float64",
		"string",
		"bool":
		return s
	}

	prev := ' '
	return strings.Map(
		func(r rune) rune {
			if isSeparator(prev) {
				prev = r
				return unicode.ToTitle(r)
			}
			prev = r
			return r
		},
		s)
}

// isSeparator returns true if the rune is considered to separate type names
// into parts. Parts are letters and digits. Whitespace and underscore are
// separators.
func isSeparator(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	return unicode.IsSpace(r) || r == '_'
}

// seqTypeOf returns the RTI DDS C type for golang build-in types.
func seqTypeOf(t string) string {
	switch t {
	case "bool":
		return "DDS_Boolean"
	case "int16":
		return "DDS_Short"
	case "uint16":
		return "DDS_UnsignedShort"
	case "int32":
		return "DDS_Long"
	case "uint32":
		return "DDS_UnsignedLong"
	case "float32":
		return "DDS_Float"
	case "float64":
		return "DDS_Double"
	case "string":
		return "DDS_String"
	}
	return t
}

// The template for the struct/union file. If this template was a separate
// file, then it would have to travel with the parser binary.
var structFileTmpl = `// THIS IS AN AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package {{.PackageName}}

import (
  "errors"
  "fmt"
  "rtiddsgo"
  "unsafe"
)

// #cgo CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT -m64 -I{{.CRtiIncludePath}}/include -I{{.CRtiIncludePath}}/include/ndds -I/usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L{{.CRtiIncludePath}}/lib/x64Linux3gcc4.8.2 -lnddsczd -lnddscorezd -ldl -lnsl -lm -lpthread -lrt -m64 -Wl,--no-as-needed
// #include <stdlib.h>
// #include <ndds/ndds_c.h>
// #include "{{.CPath}}/{{.CFileName}}.h"
// #include "{{.CPath}}/{{.CFileName}}Support.h"
import "C"

// ==========================================================================
//  Go type definition of the IDL type
// ==========================================================================

type {{.GoStructName}} struct {
{{.DiscriminantMember.GoMemberName}} {{.DiscriminantMember.GoFullTypeName}}
{{ range $member := .Members}}
	{{$member.GoComment}}
  {{$member.GoMemberName}} {{$member.GoFullTypeName}}
{{ end }}
}

// ==========================================================================
//  store / unstore functions for publishing data
// ==========================================================================

{{define "storeBool"}}
	if {{.GoFromValue}} {
		{{.CToValue}} = 1
	} else {
		{{.CToValue}} = 0
	}
{{- end -}}

{{define "storeInt16"}}
	{{.CToValue}} = C.DDS_Short({{.GoFromValue}})
{{- end -}}

{{define "storeInt32"}}
	{{.CToValue}} = C.DDS_Long({{.GoFromValue}})
{{- end -}}

{{define "storeUint16"}}
	{{.CToValue}} = C.DDS_UnsignedShort({{.GoFromValue}})
{{- end -}}

{{define "storeUint32"}}
	{{.CToValue}} = C.DDS_UnsignedLong({{.GoFromValue}})
{{- end -}}

{{define "storeFloat32"}}
	{{.CToValue}} = C.DDS_Float({{.GoFromValue}})
{{- end -}}

{{define "storeFloat64"}}
	{{.CToValue}} = C.DDS_Double({{.GoFromValue}})
{{- end -}}

{{define "storeString"}}
	str := C.CString({{.GoFromValue}})
	C.strcpy((*C.char)({{.CToValue}}), str)
	C.free(unsafe.Pointer(str))
{{- end -}}

{{define "storeEnum"}}
	{{.CToValue}} = C.{{.CTypeName}}({{.GoFromValue}})
{{- end -}}

{{define "storeOther"}}
	{{.GoTypeName}}__store({{.GoFromValue}}, &{{.CToValue}})
{{- end -}}

{{define "store"}}
	{{- if eq .GoTypeName "bool" -}}
		{{template "storeBool" .}}
	{{- else if eq .GoTypeName "int16" -}}
		{{template "storeInt16" .}}
	{{- else if eq .GoTypeName "int32" -}}
		{{template "storeInt32" .}}
	{{- else if eq .GoTypeName "uint16" -}}
		{{template "storeUint16" .}}
	{{- else if eq .GoTypeName "uint32" -}}
		{{template "storeUint32" .}}
	{{- else if eq .GoTypeName "float32" -}}
		{{template "storeFloat32" .}}
	{{- else if eq .GoTypeName "float64" -}}
		{{template "storeFloat64" .}}
	{{- else if eq .GoTypeName "string" -}}
		{{template "storeString" . -}}
	{{- else if .IsEnum -}}
		{{template "storeEnum" . -}}
	{{- else -}}
		{{template "storeOther" . -}}
	{{- end -}}
{{- end -}}

{{define "storeSeq"}}
	C.{{.SeqType}}Seq_set_maximum(&{{.CTo}}, C.DDS_Long({{.ArrayLength}}))
	C.{{.SeqType}}Seq_set_length(&{{.CTo}}, C.DDS_Long(len({{.GoFrom}})))

	for index := range {{.GoFrom}} {
		value := C.{{.SeqType}}Seq_get_reference(&{{.CTo}}, C.DDS_Long(index))
  	{{- template "store" . }}

    _ = *value // Make sure the variable is used.
	}
{{- end -}}

func {{.GoStructName}}__store(data {{.GoStructName}}, instance *C.{{.CStructName}}) {
{{ if .IsUnion }}
  {{template "store" .DiscriminantMember }}
{{ end }}

{{ range $member := .Members -}}
	{{ if $member.IsUnionMember }}
		if data.{{$.DiscriminantMember.GoMemberName}} == {{$member.UnionValue}} {
	{{ end }}

	{{ if $member.IsArray }}
		{{template "storeSeq" $member }}
	{{ else}}
		{{template "store" $member }}
	{{ end }}

	{{ if $member.IsUnionMember }}
		}
	{{ end }}
{{ end }}
}

// ==========================================================================

{{define "unstore"}}
	{{- if eq .GoTypeName "bool"}}
	{{- else if eq .GoTypeName "int16"}}
	{{- else if eq .GoTypeName "int32"}}
	{{- else if eq .GoTypeName "uint16"}}
	{{- else if eq .GoTypeName "uint32"}}
	{{- else if eq .GoTypeName "float32"}}
	{{- else if eq .GoTypeName "float64"}}
	{{- else if eq .GoTypeName "string"}}
	{{- else if .IsEnum -}}
	{{else -}}
		{{.GoTypeName}}__unstore({{.GoFromValue}}, &instance.{{.CMemberName}})
	{{end -}}
{{end}}

{{define "unstoreSeq"}}
	for index := range {{.GoFrom}} {
		value := C.{{.SeqType}}Seq_get_reference(&{{.CTo}}, C.DDS_Long(index))
  	{{template "unstore" . }}

    _ = *value // Make sure the variable is used.
	}
{{end}}

func {{.GoStructName}}__unstore(data {{.GoStructName}}, instance *C.{{.CStructName}}) {
{{ if .IsUnion }}
  {{template "unstore" .DiscriminantMember }}
{{ end }}

{{range $member := .Members -}}
	{{ if $member.IsUnionMember }}
		if data.{{$.DiscriminantMember.GoMemberName}} == {{$member.UnionValue}} {
	{{ end }}

	{{ if $member.IsArray }}
		{{template "unstoreSeq" $member }}
	{{ else}}
		{{template "unstore" $member }}
	{{ end }}

	{{ if $member.IsUnionMember }}
		}
	{{ end }}
{{end}}
}

// ==========================================================================
//  retrieve / free functions for subscribing to data
// ==========================================================================

{{define "retrieveBool"}}
	{{.GoToValue}} =  {{.CFromValue}} == 1
{{- end -}}

{{define "retrieveInt16"}}
	{{.GoToValue}} = int16({{.CFromValue}})
{{- end -}}

{{define "retrieveInt32"}}
	{{.GoToValue}} = int32({{.CFromValue}})
{{- end -}}

{{define "retrieveUint16"}}
	{{.GoToValue}} = uint16({{.CFromValue}})
{{- end -}}

{{define "retrieveUint32"}}
	{{.GoToValue}} = uint32({{.CFromValue}})
{{- end -}}

{{define "retrieveFloat32"}}
	{{.GoToValue}} = float32({{.CFromValue}})
{{- end -}}

{{define "retrieveFloat64"}}
	{{.GoToValue}} = float64({{.CFromValue}})
{{- end -}}

{{define "retrieveEnum"}}
	{{.GoToValue}} = {{.GoTypeName}}({{.CFromValue}})
{{- end -}}

{{define "retrieveString"}}
	{{.GoToValue}} = C.GoString((*C.char)({{.CFromValue}}))
{{- end -}}

{{define "retrieveOther"}}
	{{.GoTypeName}}__retrieve(&{{.GoFromValue}}, &{{.CToValue}})
{{- end -}}

{{define "retrieve"}}
	{{- if eq .GoTypeName "bool" -}}
		{{template "retrieveBool" .}}
	{{- else if eq .GoTypeName "int16" -}}
		{{template "retrieveInt16" .}}
	{{- else if eq .GoTypeName "int32" -}}
		{{template "retrieveInt32" .}}
	{{- else if eq .GoTypeName "uint16" -}}
		{{template "retrieveUint16" .}}
	{{- else if eq .GoTypeName "uint32" -}}
		{{template "retrieveUint32" .}}
	{{- else if eq .GoTypeName "float32" -}}
		{{template "retrieveFloat32" .}}
	{{- else if eq .GoTypeName "float64" -}}
		{{template "retrieveFloat64" .}}
	{{- else if eq .GoTypeName "string" -}}
		{{template "retrieveString" . -}}
	{{- else if .IsEnum -}}
		{{template "retrieveEnum" . -}}
	{{- else -}}
		{{template "retrieveOther" . -}}
	{{- end -}}
{{- end -}}

{{define "retrieveSeq"}}
	length := C.{{.SeqType}}Seq_get_length(&{{.CFrom}})
	{{.GoTo}} = make([]{{.GoTypeName}}, length)

	for index := 0; index < int(length); index++ {
		value := C.{{.SeqType}}Seq_get_reference(&{{.CFrom}}, C.DDS_Long(index))
  	{{- template "retrieve" . }}
	}
{{- end -}}

func {{.GoStructName}}__retrieve(data *{{.GoStructName}}, instance *C.{{.CStructName}}) {
{{ if .IsUnion }}
  {{template "retrieve" .DiscriminantMember }}
{{ end }}

{{ range $member := .Members -}}
	{{ if $member.IsUnionMember }}
		if data.{{$.DiscriminantMember.GoMemberName}} == {{$member.UnionValue}} {
	{{ end }}

	{{ if $member.IsArray }}
		{{template "retrieveSeq" $member }}
	{{ else}}
		{{template "retrieve" $member }}
	{{ end }}

	{{ if $member.IsUnionMember }}
		}
	{{ end }}
{{end}}
}

func {{.GoStructName}}__free(data *{{.GoStructName}}, instance *C.{{.CStructName}}) {
{{range $member := .Members -}}
	{{ if $member.IsUnionMember }}
		if data.{{$.DiscriminantMember.GoMemberName}} == {{$member.UnionValue}} {
	{{ end }}

	{{- if eq $member.GoTypeName "bool"}}
	{{- else if eq $member.GoTypeName "int16"}}
	{{- else if eq $member.GoTypeName "int32"}}
	{{- else if eq $member.GoTypeName "uint16"}}
	{{- else if eq $member.GoTypeName "uint32"}}
	{{- else if eq $member.GoTypeName "float32"}}
	{{- else if eq $member.GoTypeName "float64"}}
	{{- else if eq $member.GoTypeName "string"}}
	{{- else if .IsEnum -}}
	{{else -}}
		{{$member.GoTypeName}}__free(&{{.GoFromValue}}, &instance.{{$member.CMemberName}})
	{{end -}}

	{{ if $member.IsUnionMember }}
		}
	{{ end }}
{{end}}
}

// ==========================================================================
//  Type functions
// ==========================================================================

func {{.GoStructName}}_GetTypeName() string {
	return C.GoString(C.{{.CStructName}}TypeSupport_get_type_name())
}

func {{.GoStructName}}_RegisterType(p rtiddsgo.Participant) error {
	rc := C.{{.CStructName}}TypeSupport_register_type(
		(*C.DDS_DomainParticipant)(p.GetUnsafe()),
		C.{{.CStructName}}TypeSupport_get_type_name())
	if rc != C.DDS_RETCODE_OK {
		return errors.New("Failed to register the type {{.GoStructName}}.")
	}
	return nil
}

// ==========================================================================
//  Data Writer
// ==========================================================================

type {{.GoStructName}}DataWriter struct {
	dw rtiddsgo.DataWriter
	cdw *C.{{.CStructName}}DataWriter
}

func New{{.GoStructName}}DataWriter(pub rtiddsgo.Publisher, t rtiddsgo.Topic, qosLibraryName, qosProfileName string) ({{.GoStructName}}DataWriter, error) {
	messageDW := {{.GoStructName}}DataWriter{}
	var err error
  messageDW.dw, err = rtiddsgo.CreateDataWriter(pub, t, qosLibraryName, qosProfileName)
	if err != nil {
		return messageDW, err
	}
	messageDW.cdw = C.{{.CStructName}}DataWriter_narrow((*C.DDS_DataWriter)(messageDW.dw.GetUnsafe()))
	return messageDW, nil
}

func (dw {{.GoStructName}}DataWriter)Free() {
	dw.dw.Free()
}

func (dw {{.GoStructName}}DataWriter)Write(m {{.GoStructName}}) error {
	instance := C.{{.CStructName}}TypeSupport_create_data()
	{{.GoStructName}}__store(m, instance)

	rc := C.{{.CStructName}}DataWriter_write(
		dw.cdw,
		instance,
		&C.DDS_HANDLE_NIL)
	defer func() {
		{{.GoStructName}}__unstore(m, instance)
	  C.{{.CStructName}}TypeSupport_delete_data(instance)
	}()
	if rc != C.DDS_RETCODE_OK {
		return fmt.Errorf("Failed to write. Return code was %s", rtiddsgo.ReturnCodeToString(int(rc)))
	}
	return nil
}

// ==========================================================================
//  Data Reader
// ==========================================================================

type {{.GoStructName}}DataReader struct {
	dr rtiddsgo.DataReader
	cdr *C.{{.CStructName}}DataReader
}

func New{{.GoStructName}}DataReader(sub rtiddsgo.Subscriber, t rtiddsgo.Topic, qosLibraryName, qosProfileName string, rxFunc func(data {{.GoStructName}})) ({{.GoStructName}}DataReader, error) {
	messageDR := {{.GoStructName}}DataReader{}

	var err error
	messageDR.dr, err = rtiddsgo.CreateDataReader(sub, t, qosLibraryName, qosProfileName, func() {
		var rc C.DDS_ReturnCode_t
		for rc = C.DDS_RETCODE_OK; rc != C.DDS_RETCODE_NO_DATA; {
			var dataSeq C.struct_{{.CStructName}}Seq
			var sampleInfoSeq C.struct_DDS_SampleInfoSeq
			rc = C.{{.CStructName}}DataReader_take(
				messageDR.cdr,
				&dataSeq,
				&sampleInfoSeq,
				C.DDS_LENGTH_UNLIMITED,
				C.DDS_ANY_SAMPLE_STATE,
				C.DDS_ANY_VIEW_STATE,
				C.DDS_ANY_INSTANCE_STATE)
			if rc != C.DDS_RETCODE_NO_DATA && rc != C.DDS_RETCODE_OK {
				return
			} else if rc == C.DDS_RETCODE_OK {
				for i := C.DDS_Long(0); i < C.{{.CStructName}}Seq_get_length(&dataSeq); i++ {
					if C.DDS_SampleInfoSeq_get_reference(&sampleInfoSeq, i).valid_data == C.DDS_BOOLEAN_TRUE {
						rxData := C.{{.CStructName}}Seq_get_reference(&dataSeq, i)

						var goData {{.GoStructName}}
						{{.GoStructName}}__retrieve(&goData, rxData)
						rxFunc(goData)
						{{.GoStructName}}__free(&goData, rxData)

						C.{{.CStructName}}DataReader_return_loan(messageDR.cdr, &dataSeq, &sampleInfoSeq)
					}
				}
			}
		}
	})
	if err != nil {
		return messageDR, err
	}

	messageDR.cdr = C.{{.CStructName}}DataReader_narrow((*C.DDS_DataReader)(messageDR.dr.GetUnsafe()))
	return messageDR, nil
}

func (dr {{.GoStructName}}DataReader)Free() {
	dr.dr.Free()
}

// ==========================================================================
//  Dummy to ensure that all imports are used
// ==========================================================================
type dummy_{{.CStructName}} struct {
	dummyPtr unsafe.Pointer
}
`
