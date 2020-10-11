package haystack

import (
	"bufio"
	"encoding/json"
	"strings"
)

// List is a list of Val items.
type List struct {
	vals []Val
}

// NewList creates a new List object.
func NewList(vals []Val) List {
	return List{vals: vals}
}

// MarshalJSON represents the object in JSON array format: "[<val1>, <val2>, ...]"
func (list List) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.vals)
}

// ToZinc representes the object as: "[<val1>, <val2>, ...]"
func (list List) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	list.WriteZincTo(out)
	out.Flush()
	return builder.String()
}

// WriteZincTo appends the Writer with the List zinc representation.
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