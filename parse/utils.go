package parse

import (
	"strings"
)

// Returns the path concatenated with the name and an underscore in between or an empty string if the name is an empty
// string. All instances of "::" are replaced with an underscore.
// If name is a golang reserved ("type", "range") and there is no path, then a prefix "_" is added.
func cNameOf(path, name string) string {
	if len(name) == 0 {
		return ""
	}

	// If the name is a reserved word, then an underscore is prepended when accessing the variable from Go. This only
	// applies to names without a path.
	if len(path) == 0  && (name == "type" || name == "range"){
		return "_" + name
	}

	path = strings.Replace(strings.Trim(path, "_"), "::", "_", -1)
	name = strings.Replace(strings.Trim(name, "_"), "::", "_", -1)
	return strings.Trim(path + "_" + name, "_")
}

// Returns the path concatenated with the name and an underscore in between or an empty string if name is empty.
// Path and Name are both set in Title. Words after a double colon are also Title'ed.
// All instances of "::" are replaced with an underscore.
// Returns name if the path is empty or contains only an underscore. Leading underscores are removed from the path.
func goNameOf(path, name string) string {
	if len(name) == 0 {
		return ""
	}
	path = strings.Replace(strings.Trim(toTitle(path), "_"), "::", "_", -1)
	name = strings.Replace(strings.Trim(toTitle(name), "_"), "::", "_", -1)
	return strings.Trim(path + "_" + name, "_")
}

// Returns the string with each word set in the Title case. Words are separated by a double-colon or underscore.
// The string is returned unchanged if it contains spaces.
func toTitle(s string) string {
	if strings.Contains(s, " ") {
		return s
	}

	// Convert all double colons to spaces, Title the whole string then convert the spaces back to double-colons.
	s = strings.Replace(s, "::", " ", -1)
	s = strings.Title(s)
	s = strings.Replace(s, " ", "::", -1)

	// Convert all underscores to spaces, Title the whole string then convert the spaces back to underscores.
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	s = strings.Replace(s, " ", "_", -1)

	return s
}

// Returns true if the string looks like a qualified statement, i.e. contains at least one
// double-colon.
func isAQualifiedValue(s string) bool {
	// Shortest possible qualified string is "a::b"
	if len(s) < 4 {
		return false
	}
	return strings.Contains(s, "::")
}


// Returns the golang type that corresponds to the given XML type. If "xmlType" is "nonBasic", then the nonBasic type
// is used.
func goTypeOf(xmlType, nonBasic string) string {
	return goPathTypeOf("", xmlType, nonBasic)
}

// Returns the golang type that corresponds to the given path and XML type. If "xmlType" is "nonBasic", then the
// nonBasic type is used.
func goPathTypeOf(goPath, xmlType, nonBasic string) string {
	if xmlType == "nonBasic" {
		xmlType = nonBasic
	}
	switch (xmlType) {
	case "boolean":
		return "bool"
	case "char8":
		return "int8"
	case "byte":
		return "uint8"
	case "octet":
		return "uint8"
	case "int16":
		return "int16"
	case "uint16":
		return "uint16"
	case "int32":
		return "int32"
	case "uint32":
		return "uint32"
	case "int64":
		return "int64"
	case "uint64":
		return "uint64"
	case "float32":
		return "float32"
	case "float64":
		return "float64"
	case "string":
		return "string"
	}
	return goNameOf(goPath, xmlType)
}

// Returns the DDS type that corresponds to the given XML type. If "t" is "nonBasic", then the nonBasic type is used.
// The Dds type of a nonBasic type is on the form "a_b_c" with casing preserved.
func ddsTypeOf(xmlType, nonBasic string) string {
	if xmlType == "nonBasic" {
		xmlType = nonBasic
	}
	switch xmlType {
	case "boolean":
		return "DDS_Boolean"
	case "byte":
		return "DDS_Octet"
	case "char8":
		return "DDS_Char"
	case "octet":
		return "DDS_Octet"
	case "int16":
		return "DDS_Short"
	case "uint16":
		return "DDS_UnsignedShort"
	case "int32":
		return "DDS_Long"
	case "uint32":
		return "DDS_UnsignedLong"
	case "int64":
		return "DDS_LongLong"
	case "uint64":
		return "DDS_UnsignedLongLong"
	case "float32":
		return "DDS_Float"
	case "float64":
		return "DDS_Double"
	case "string":
		return "DDS_String"
	}
	return cNameOf("", xmlType)
}

// Returns xmlType unless it is "nonbasic", in which case the nonBasic argument is returned
func xmlTypeOf(xmlType, nonBasic string) string {
	if xmlType == "nonBasic" {
		xmlType = nonBasic
	}
	return xmlType
}

// Returns the sequence length provided with any qualified names converted to cNamestyle ("::" -> "_"). Allow
// for simple arithmetic.
func cSeqLenOf(seqLen string) string {
	if strings.Contains(seqLen, " ") {
		panic ("Must not have spaces in the sequence length")
	}

	// A single character, including the "+-*/" is always left unchanged.
	if len(seqLen) <= 1 {
		return seqLen
	}

	// If the string contains arithmetic, then break up the string in substrings and
	// process each recursively.
	if strings.ContainsAny(seqLen, "+-*/") {
		s := strings.ReplaceAll(seqLen, "+", " + ")
		s = strings.ReplaceAll(s, "-", " - ")
		s = strings.ReplaceAll(s, "*", " * ")
		s = strings.ReplaceAll(s, "/", " / ")

		f := strings.Fields(s)
		result := ""
		for _, field := range f {
			result += cSeqLenOf (field)
		}
		return result
	} else if strings.Contains(seqLen, "::") {
		return "C." + cNameOf("", seqLen)
	}
	return seqLen
}

func goArrayDimsOf(arrayDims string) string {
	if strings.Contains(arrayDims, ",") {
		panic("Multi dimensional arrays are not yet supported")
	}
	arrayDims = strings.TrimLeft(arrayDims, "(")
	arrayDims = strings.TrimRight(arrayDims, ")")
	if strings.Contains(arrayDims, "::") {
		return goNameOf("", arrayDims)
	}
	return arrayDims
}