package haystack

import "strings"

type List struct {
	vals []Val
}

func NewList(vals []Val) List {
	return List{vals: vals}
}

// ToZinc representes the object as: "[<val1>, <val2>, ...]"
func (list List) ToZinc() string {
	var buf strings.Builder
	list.encodeTo(&buf)
	return buf.String()
}

// Format as [<val1>, <val2>, ...]
func (list List) encodeTo(buf *strings.Builder) {
	buf.WriteString("[")
	for idx, val := range list.vals {
		if idx != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(val.ToZinc())
	}
	buf.WriteString("]")
}
