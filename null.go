package haystack

type Null struct {
}

// ToZinc representes the object as an empty string
func (null Null) ToZinc() string {
	return ""
}
