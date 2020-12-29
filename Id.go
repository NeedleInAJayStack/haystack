package haystack

import "encoding/json"

// Id models a simple text object, typically a tag name or keyword.
// It's not used in the tagging model, but is used by the parser and tokenizer.
type Id struct {
	val string
}

// NewId creates a new Id object.
func NewId(val string) Id {
	return Id{val: val}
}

// String returns the value of the string.
func (id Id) String() string {
	return id.val
}

// ToZinc representes the object as: "<val>"
func (id Id) ToZinc() string {
	return id.val
}

// MarshalJSON representes the object as: "<val>"
func (id Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.val)
}

// MarshalHayson representes the object as: "<val>"
func (id Id) MarshalHayson() ([]byte, error) {
	return json.Marshal(id.val)
}
