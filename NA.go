package haystack

import "encoding/json"

// NA is the value used to indicate not available.
type NA struct {
}

// NewNA creates a new NA object.
func NewNA() NA {
	return NA{}
}

// ToZinc representes the object as: "NA"
func (na NA) ToZinc() string {
	return "NA"
}

// MarshalJSON representes the object as: "z:"
func (na NA) MarshalJSON() ([]byte, error) {
	return json.Marshal("z:")
}

// MarshalHAYSON representes the object as: "{\"_kind\":\"na\"}"
func (na NA) MarshalHAYSON() ([]byte, error) {
	return []byte("{\"_kind\":\"na\"}"), nil
}
