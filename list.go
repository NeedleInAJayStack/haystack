package haystack

import (
	"bufio"
	"strings"
)

type List struct {
	vals []Val
}

func NewList(vals []Val) List {
	return List{vals: vals}
}

// ToZinc representes the object as: "[<val1>, <val2>, ...]"
func (list List) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	list.WriteZincTo(out)
	out.Flush()
	return builder.String()
}

// Format as [<val1>, <val2>, ...]
func (list List) WriteZincTo(buf *bufio.Writer) {
	buf.WriteString("[")
	for idx, val := range list.vals {
		if idx != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(val.ToZinc())
	}
	buf.WriteString("]")
}
