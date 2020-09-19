package haystack

type XStr struct {
	valType string
	val     string
}

// Format is <valType>(<val>)
func (xStr XStr) ToZinc() string {
	result := xStr.valType
	result = result + "(\""
	result = result + xStr.val
	result = result + "\")"

	return result
}
