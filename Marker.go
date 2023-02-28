package haystack

import (
	"encoding/json"
	"errors"
	"strings"
)

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

// UnmarshalJSON interprets the json value: "m:"
func (marker *Marker) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newMarker, newErr := markerFromJSON(jsonStr)
	*marker = newMarker
	return newErr
}

func markerFromJSON(jsonStr string) (Marker, error) {
	if !strings.HasPrefix(jsonStr, "m:") {
		return NewMarker(), errors.New("Input value does not begin with 'm:'")
	}
	return NewMarker(), nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"marker\"}"
func (marker Marker) MarshalHayson() ([]byte, error) {
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

// UnmarshalJSON interprets the json value: "-:"
func (remove Remove) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newRemove, newErr := removeFromJSON(jsonStr)
	remove = newRemove
	return newErr
}

func removeFromJSON(jsonStr string) (Remove, error) {
	if !strings.HasPrefix(jsonStr, "-:") {
		return NewRemove(), errors.New("Input value does not begin with '-:'")
	}
	return NewRemove(), nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"remove\"}"
func (remove Remove) MarshalHayson() ([]byte, error) {
	return []byte("{\"_kind\":\"remove\"}"), nil
}
