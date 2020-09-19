package haystack

type Null struct {
}

func (null Null) ToZinc() string {
	return ""
}
