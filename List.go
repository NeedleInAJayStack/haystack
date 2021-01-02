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
func NewList(vals []Val) *List {
	return &List{vals: vals}
}

// Get returns the val at the given index
func (list *List) Get(index int) Val {
	return list.vals[index]
}

// Size returns the number of vals in the list
func (list *List) Size() int {
	return len(list.vals)
}

// MarshalJSON represents the object in JSON array format: "[<val1>, <val2>, ...]"
func (list *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.vals)
}

// UnmarshalJSON interprets the JSON array format: "[<val1>, <val2>, ...]"
func (list *List) UnmarshalJSON(buf []byte) error {
	var jsonList []interface{}
	err := json.Unmarshal(buf, &jsonList)
	if err != nil {
		return err
	}

	newList, newErr := listFromJSON(jsonList)
	*list = *newList
	return newErr
}

func listFromJSON(jsonList []interface{}) (*List, error) {
	items := []Val{}
	for _, jsonVal := range jsonList {
		val, err := ValFromJSON(jsonVal)
		if err != nil {
			return nil, err
		}
		items = append(items, val)
	}

	return NewList(items), nil
}

// MarshalHayson represents the object in JSON array format: "[<val1>, <val2>, ...]"
func (list *List) MarshalHayson() ([]byte, error) {
	builder := new(strings.Builder)
	builder.WriteString("[")
	for idx, val := range list.vals {
		if idx != 0 {
			builder.WriteString(",")
		}
		valHayson, valErr := val.MarshalHayson()
		if valErr != nil {
			return []byte{}, valErr
		}
		builder.Write(valHayson)
	}
	builder.WriteString("]")
	return []byte(builder.String()), nil
}

// ToZinc representes the object as: "[<val1>, <val2>, ...]"
func (list *List) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	list.WriteZincTo(out)
	out.Flush()
	return builder.String()
}

// WriteZincTo appends the Writer with the List zinc representation.
func (list *List) WriteZincTo(buf *bufio.Writer) {
	buf.WriteString("[")
	for idx, val := range list.vals {
		if idx != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(val.ToZinc())
	}
	buf.WriteString("]")
}
