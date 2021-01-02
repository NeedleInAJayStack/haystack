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

// UnmarshalJSON interprets the json value: "z:"
func (na *NA) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(jsonStr, "z:") {
		return errors.New("Input value does not begin with z:")
	}

	*na = NA{}

	return nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"na\"}"
func (na NA) MarshalHayson() ([]byte, error) {
	return []byte("{\"_kind\":\"na\"}"), nil
}
