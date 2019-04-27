package generate
import "C"

func Store(goType, cType, from, to, toRef string) string {
	switch goType {
	case "bool":
		return "if " + from + " { " + to + " = 1 } else { " + to + " = 0 }"
	case "int16", "uint16", "int32", "uint32", "float32", "float64":
		return to + " = C." + cType + "(" + from + ")"
	case "string":
		return `
{
    str := C.CString(string(` + from + `))
    C.strcpy((*C.char)(` + to + `), str)
    C.free(unsafe.Pointer(str))
}
`
	default:
		return from + ".Store(" + toRef + ")"
	}
}
