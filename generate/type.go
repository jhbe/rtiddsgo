package generate

func Type(goType, seqLen, arrayDims string) string {
	s := goType
	if len(seqLen) > 0 {
		s = "[]"+s
	}
	if len(arrayDims) > 0 {
		s = "[" + arrayDims + "]" + s
	}
	return s
}
