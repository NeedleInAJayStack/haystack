package haystack

import "encoding/json"

// Val represents a haystack tag value.
type Val interface {
	ToZinc() string
	json.Marshaler
	json.Unmarshaler
	HaysonMarshaller
}

type HaysonMarshaller interface {
	MarshalHayson() ([]byte, error)
}
