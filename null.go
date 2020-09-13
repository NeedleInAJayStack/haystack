package haystack

type Null struct {
}

func (null Null) toZinc() string {
	return ""
}
