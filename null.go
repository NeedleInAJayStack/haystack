package haystack

type Null struct {
}

func NewNull() Null {
	return Null{}
}

// ToZinc representes the object as an empty string
func (null Null) ToZinc() string {
	return "N"
}
