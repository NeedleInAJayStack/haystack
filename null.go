package haystack

// Null is the value used to indicate a Val with no type.
type Null struct {
}

// NewNull creates a new Null object.
func NewNull() Null {
	return Null{}
}

// ToZinc representes the object as an empty string
func (null Null) ToZinc() string {
	return "N"
}
