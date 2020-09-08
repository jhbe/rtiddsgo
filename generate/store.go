package generate
import "C"

func Store(goType, cType, from, to, toRef, seqLen, arrayDims string) string {
	if len(arrayDims) > 0 {
		return `for index, _ := range `+from+` {
        `+Store(goType, cType, from+"[index]", to+"[index]", "&("+to+"[index])", seqLen, "")+`
    }`
	}

	if len(seqLen) > 0 {
	return `C.`+cType+`Seq_set_maximum(`+toRef+`, C.DDS_Long(`+seqLen+`))
    C.`+cType+`Seq_set_length(`+toRef+`, C.DDS_Long(len(`+from+`)))
    for index, _ := range `+from+` {
	    value := C.`+cType+`Seq_get_reference(`+toRef+`, C.DDS_Long(index))
        `+Store(goType, cType, from+"[index]", "*value", "value", "", "")+`
    }`
	}

	switch goType {
	case "bool":
		return "if " + from + " { " + to + " = 1 } else { " + to + " = 0 }"
	case "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64":
		return to + " = C." + cType + "(" + from + ")"
	case "string":
		return `{
    str := C.CString(string(` + from + `))
    C.strcpy((*C.char)(` + to + `), str)
    C.free(unsafe.Pointer(str))
}`
	default:
		return from + ".Store((*C."+cType+")(" + toRef + "))"
	}
}
