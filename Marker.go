package haystack

import "encoding/json"

// Marker is the value for a marker tag.
type Marker struct {
}

// NewMarker creates a new Marker object.
func NewMarker() Marker {
	return Marker{}
}

// ToZinc representes the object as: "M"
func (marker Marker) ToZinc() string {
	return "M"
}

// MarshalJSON representes the object as: "m:"
func (marker Marker) MarshalJSON() ([]byte, error) {
	return json.Marshal("m:")
}

// MarshalHAYSON representes the object as: "{\"_kind\":\"marker\"}"
func (marker Marker) MarshalHAYSON() ([]byte, error) {
	return []byte("{\"_kind\":\"marker\"}"), nil
}

// Remove is the value used to indicate a tag remove.
type Remove struct {
}

// NewRemove creates a new Remove object.
func NewRemove() Remove {
	return Remove{}
}

// ToZinc representes the object as: "R"
func (remove Remove) ToZinc() string {
	return "R"
}

// MarshalJSON representes the object as: "-:"
func (remove Remove) MarshalJSON() ([]byte, error) {
	return json.Marshal("-:")
}

// MarshalHAYSON representes the object as: "{\"_kind\":\"remove\"}"
func (remove Remove) MarshalHAYSON() ([]byte, error) {
	return []byte("{\"_kind\":\"remove\"}"), nil
}
