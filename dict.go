package haystack

import "strings"

type Dict struct {
	items map[string]Val
}

func (dict *Dict) isEmpty() bool {
	return len(dict.items) == 0
}

func (dict Dict) toZinc() string {
	var buf strings.Builder
	dict.encodeTo(&buf, true)
	return buf.String()
}

// Format is {name1:val1 name2:val2 ...}. Markers don't require a :val.
func (dict Dict) encodeTo(buf *strings.Builder, brackets bool) {
	if brackets {
		buf.WriteString("{")
	}
	firstVal := true
	for name, val := range dict.items {
		if firstVal {
			firstVal = false
		} else {
			buf.WriteString(" ")
		}

		buf.WriteString(name)

		_, isMarker := val.(*Marker)
		if !isMarker {
			buf.WriteString(":" + val.toZinc())
		}
	}
	if brackets {
		buf.WriteString("}")
	}
}
