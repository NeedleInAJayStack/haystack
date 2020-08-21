package haystack

// A number is composed of a val and unit
type XStr struct {
	valType string
	val     string
}

/** Encode using double quotes and back slash escapes */
func (xStr *XStr) toZinc() string {
	result := xStr.valType
	result = result + "(\""
	result = result + xStr.val
	result = result + "\")"

	return result
}
