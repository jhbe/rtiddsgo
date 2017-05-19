package rtiddsgo

var returnCodeStrings = []string{
	"OK",
	"ERROR",
	"UNSUPPORTED",
	"BAD_PARAMETER",
	"PRECONDITION_NOT_MET",
	"OUT_OF_RESOURCES",
	"NOT_ENABLED",
	"IMMUTABLE_POLICY",
	"INCONSISTENT_POLICY",
	"ALREADY_DELETED",
	"TIMEOUT",
	"NO_DATA",
	"ILLEGAL_OPERATION",
}

func ReturnCodeToString(rc int) string {
	if 0 <= rc && rc < 13 {
		return returnCodeStrings[rc]
	}
	return "<unknown>"
}
