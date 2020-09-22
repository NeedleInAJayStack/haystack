package haystack

import (
	"bufio"
	"sort"
	"strings"
)

type Dict struct {
	items map[string]Val
}

func NewDict(items map[string]Val) Dict {
	return Dict{items: items}
}

func (dict *Dict) Get(key string) Val {
	val := dict.items[key]
	if val == nil {
		val = NewNull()
	}
	return val
}

func (dict *Dict) Size() int {
	return len(dict.items)
}

func (dict *Dict) IsEmpty() bool {
	return len(dict.items) == 0
}

// ToZinc representes the object as: "{<name1>:<val1> <name2>:<val2> ...}" with the names in alphabetical order. Markers don't require a val.
func (dict Dict) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	dict.WriteZincTo(out, true)
	out.Flush()
	return builder.String()
}

// Format is {<name1>:<val1> <name2>:<val2> ...} with the names in alphabetical order. Markers don't require a val.
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
