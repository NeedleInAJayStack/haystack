package haystack

type XStr struct {
	valType string
	val     string
}

func NewXStr(valType string, val string) XStr {
	return XStr{valType: valType, val: val}
}

func (xStr XStr) Type() string {
	return xStr.valType
}

func (xStr XStr) String() string {
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
