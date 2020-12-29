package haystack

import (
	"encoding/json"
	"strings"
)

// XStr is an extended string which is a type name and value encoded as a string.
// It is used as a generic value when an XStr is decoded without any predefined type.
type XStr struct {
	valType string
	val     string
}

// NewXStr creates a new XStr object.
func NewXStr(valType string, val string) XStr {
	return XStr{valType: valType, val: val}
}

// Type returns the XStr object type
func (xStr XStr) Type() string {
	return xStr.valType
}

// Val returns the XStr object value
func (xStr XStr) Val() string {
	return xStr.val
}

// ToZinc representes the object as: <valType>("<val>")
func (xStr XStr) ToZinc() string {
	result := xStr.valType
	result = result + "(\""
	result = result + xStr.val
	result = result + "\")"

	return result
}

// MarshalJSON representes the object as: "x:<valType>:<val>"
func (xStr XStr) MarshalJSON() ([]byte, error) {
	result := "x:"
	result = result + xStr.valType
	result = result + ":"
	result = result + xStr.val

	return json.Marshal(result)
}

// MarshalHayson representes the object as: "{\"_kind\":\"xstr\",\"type\":\"<valType>\",\"val\":\"<val>\"}"
func (xStr XStr) MarshalHayson() ([]byte, error) {
	builder := new(strings.Builder)
	builder.WriteString("{\"_kind\":\"xstr\",\"type\":\"")
	builder.WriteString(xStr.valType)
	builder.WriteString("\",\"val\":\"")
	builder.WriteString(xStr.val)
	builder.WriteString("\"}")
	return []byte(builder.String()), nil
}
