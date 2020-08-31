package haystack

import "strings"

type List struct {
	vals []Val
}

// Format as [<val1>, <val2>, ...]
func (list *List) toZinc() string {
	var buf strings.Builder
	list.encodeTo(&buf)
	return buf.String()
}

func (list *List) encodeTo(buf *strings.Builder) {
	buf.WriteString("[")
	for idx, val := range list.vals {
		if idx != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(val.toZinc())
	}
	buf.WriteString("]")
}
