package haystack

// Null is the value used to indicate a Val with no type.
type Null struct {
}

// NewNull creates a new Null object.
func NewNull() Null {
	return Null{}
}

// ToZinc representes the object as "N"
func (null Null) ToZinc() string {
	return "N"
}

// ToJSON representes the object as "null"
func (null Null) ToJSON() string {
	return "null"
}
