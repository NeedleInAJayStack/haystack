package haystack

type Uri struct {
	val string
}

// Format is `<val>`
func (uri Uri) ToZinc() string {
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
