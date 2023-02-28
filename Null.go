package haystack

import (
	"encoding/json"
	"errors"
)

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

// UnmarshalJSON interprets the json value: "null"
// NOTE THAT THIS DOESN'T WORK with json.Unmarshal because Unmarshallers implement UnmarshalJSON([]byte("null")) as a no-op,
// meaning that we always short-circuit and get nil instead of NewNull(). See https://golang.org/pkg/encoding/json/#Unmarshal
func (null *Null) UnmarshalJSON(buf []byte) error {
	var jsonNull interface{}
	err := json.Unmarshal(buf, &jsonNull)
	if err != nil {
		return err
	}

	newNull, newErr := nullFromJSON(jsonNull)
	*null = newNull
	return newErr
}

func nullFromJSON(jsonNull interface{}) (Null, error) {
	if jsonNull != nil {
		return NewNull(), errors.New("json value was not unmarshalled as nil")
	}
	return NewNull(), nil
}

// MarshalHayson representes the object as "null"
func (null Null) MarshalHayson() ([]byte, error) {
	return json.Marshal(nil)
}
