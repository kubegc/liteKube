package help

import (
	"fmt"
	"io"
	"strings"
)

type Tip struct {
	Name         string
	ValueType    string
	TipStr       string
	DefaultValue string
}

func (tip Tip) Fprint(w io.Writer, offset int, keyLength int, typeLength int) {
	fmt.Fprintf(w, "%s%s %s,  %s", strings.Repeat(" ", offset), formatString(tip.Name, keyLength, ":"), formatString(tip.ValueType, typeLength, ""), tip.TipStr)

	if len(tip.DefaultValue) > 0 {
		fmt.Fprintf(w, "(Default: %s)\n", tip.DefaultValue)
	} else {
		fmt.Fprintf(w, "\n")
	}

}

func formatString(str string, length int, insert string) string {
	s := str
	s += insert
	if length > len(str) {
		s += strings.Repeat(" ", length-len(str))
	}

	return s
}
