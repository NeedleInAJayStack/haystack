package haystack

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

// ToZinc representes the object as: "<valType>(<val>)"
func (xStr XStr) ToZinc() string {
	result := xStr.valType
	result = result + "(\""
	result = result + xStr.val
	result = result + "\")"

	return result
}
