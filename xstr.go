package haystack

type XStr struct {
	valType string
	val     string
}

func NewXStr(valType string, val string) XStr {
	return XStr{valType: valType, val: val}
}

// ToZinc representes the object as: "<valType>(<val>)"
func (xStr XStr) ToZinc() string {
	result := xStr.valType
	result = result + "(\""
	result = result + xStr.val
	result = result + "\")"

	return result
}
