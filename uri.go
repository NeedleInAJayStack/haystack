package haystack

type Uri struct {
	val string
}

func NewUri(val string) Uri {
	return Uri{val: val}
}

func (uri Uri) String() string {
	return uri.val
}

// ToZinc representes the object as: "`<val>`"
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
