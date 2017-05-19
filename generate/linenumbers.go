package generate

import (
	"bufio"
	"bytes"
	"fmt"
)

func withLineNumbers (in string) string {
	lineNumber := 1
	s := ""
	scanner := bufio.NewScanner(bytes.NewBufferString(in))
	for scanner.Scan() {
		s += fmt.Sprintf("%d: %s\n", lineNumber, scanner.Text())
		lineNumber++
	}
	return s
}