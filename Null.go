package haystack

import "encoding/json"

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

// MarshalJSON representes the object as "null"
func (null Null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

// MarshalHAYSON representes the object as "null"
func (null Null) MarshalHAYSON() ([]byte, error) {
	return json.Marshal(nil)
}
