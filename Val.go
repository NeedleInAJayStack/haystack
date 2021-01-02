package haystack

import "encoding/json"

// Val represents a haystack tag value.
type Val interface {
	ToZinc() string
	json.Marshaler
	// json.Unmarshaler
	MarshalHayson() ([]byte, error)
}
