package generate

import "C"

func Retrieve(goName, goType, to, from string, isTypeDef bool) string {
	retType := goType
	if isTypeDef {
		retType = goName
	}

	switch goType {
	case "bool":
		return to + " = " + from + " == 1"
	case "int16", "uint16", "int32", "uint32", "float32", "float64":
		return to + " = " + retType + "(" + from + ")"
	case "string":
		return to + " = " + retType + "(C.GoString((*C.char)(" + from + ")))"
	default:
		return to + ".Retrieve(" + from + ")"
	}
}
