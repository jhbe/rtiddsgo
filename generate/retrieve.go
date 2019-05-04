package generate

import "C"

func Retrieve(goName, goType, cType, to, from, seqLen, arrayDims string, isTypeDef bool) string {
	retType := goType
	if isTypeDef {
		retType = goName
	}

	if len(arrayDims) > 0 {
		return `for index, _ := range ` + to + `{
`+Retrieve(goName, goType, cType, "("+to+")[index]", from+"[index]", seqLen, "", false)+`
}`
	}

	if len(seqLen) > 0 {
		return `	`+to+` = make([]`+goType+`, C.`+cType+`Seq_get_length(&`+from+`))
	for index, _ := range `+to+` {
		value := C.`+cType+`Seq_get_reference(&`+from+`, C.DDS_Long(index))
		`+Retrieve(goName, goType, cType, "("+to+")[index]", "*value", "", "", false)+`
	}`
	}

	switch goType {
	case "bool":
		return to + " = " + from + " == 1"
	case "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64":
		return to + " = " + retType + "(" + from + ")"
	case "string":
		return to + " = " + retType + "(C.GoString((*C.char)(" + from + ")))"
	default:
		return to + ".Retrieve(" + from + ")"
	}
}
