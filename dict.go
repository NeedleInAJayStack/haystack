package haystack

import (
	"sort"
	"strings"
)

type Dict struct {
	items map[string]Val
}

func (dict *Dict) IsEmpty() bool {
	return len(dict.items) == 0
}

func NewDict(items map[string]Val) Dict {
	return Dict{items: items}
}

// ToZinc representes the object as: "{<name1>:<val1> <name2>:<val2> ...}" with the names in alphabetical order. Markers don't require a val.
func (dict Dict) ToZinc() string {
	var buf strings.Builder
	dict.encodeTo(&buf, true)
	return buf.String()
}

// Format is {<name1>:<val1> <name2>:<val2> ...} with the names in alphabetical order. Markers don't require a val.
func (dict Dict) encodeTo(buf *strings.Builder, brackets bool) {
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
		_, isMarker := val.(Marker)
		if !isMarker {
			buf.WriteString(":" + val.ToZinc())
		}
	}

	if brackets {
		buf.WriteString("}")
	}
}
