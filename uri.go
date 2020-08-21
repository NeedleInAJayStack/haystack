package haystack

// A number is composed of a val and unit
type Uri struct {
	val string
}

/** Encode using double quotes and back slash escapes */
func (uri *Uri) toZinc() string {
	result := "`"

	for i := 0; i < len(uri.val); i++ {
		char := uri.val[i]
		// URIs cannot contain characters < ' ', so just ignore them.
		if char > ' ' {
			result = result + string(char)
		}
	}
	result = result + "`"

	return result
}
