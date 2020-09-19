package haystack

type XStr struct {
	valType string
	val     string
}

// ToZinc representes the object as: "<valType>(<val>)"
func (xStr XStr) ToZinc() string {
	result := xStr.valType
	result = result + "(\""
	result = result + xStr.val
	result = result + "\")"

	return result
}
