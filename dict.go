package haystack

import (
	"bufio"
	"sort"
	"strings"
)

// Dict is a map of name/Val pairs.
type Dict struct {
	items map[string]Val
}

// NewDict creates a new Dict object.
func NewDict(items map[string]Val) Dict {
	return Dict{items: items}
}

// Get returns the Val of the given name. If the name is not found, Null is returned.
func (dict *Dict) Get(name string) Val {
	val := dict.items[name]
	if val == nil {
		val = NewNull()
	}
	return val
}

// Size returns the number of name/Val pairs.
func (dict *Dict) Size() int {
	return len(dict.items)
}

// IsEmpty returns true if there is nothing in the Dict.
func (dict *Dict) IsEmpty() bool {
	return len(dict.items) == 0
}

// ToZinc representes the object as: "{<name1>:<val1> <name2>:<val2> ...}" with the names in alphabetical order.
// Markers don't require a val.
func (dict Dict) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	dict.WriteZincTo(out, true)
	out.Flush()
	return builder.String()
}

// WriteZincTo appends the Writer with the Dict zinc representation.
func (dict Dict) WriteZincTo(buf *bufio.Writer, brackets bool) {
	if brackets {
		buf.WriteString("{")
	}

	var names []string
	for name := range dict.items {
		names = append(names, name)
	}
	sort.Strings(names)
	for idx, name := range names {
		if idx != 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(name)

		val := dict.items[name]
		switch val := val.(type) {
		case Grid:
			val.WriteZincTo(buf, 1)
		case Marker:
			break
		default:
			buf.WriteString(":" + val.ToZinc())
		}
	}

	if brackets {
		buf.WriteString("}")
	}
}
