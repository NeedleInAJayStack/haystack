package haystack

import "encoding/json"

// Bool models a boolean for true/false tag values.
type Bool bool

const (
	TRUE  Bool = true
	FALSE Bool = false
)

// ToBool returns the value of this object as a Go bool
func (b Bool) ToBool() bool {
	if b == TRUE {
		return true
	}
	return false
}

// ToZinc representes the object as: "T" or "F"
func (b Bool) ToZinc() string {
	if b.ToBool() {
		return "T"
	}
	return "F"
}

// MarshalJSON representes the object as: "true" or "false"
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.ToBool())
}

// UnmarshalJSON interprets the json value: "true" or "false"
func (b *Bool) UnmarshalJSON(buf []byte) error {
	var boolVal bool
	err := json.Unmarshal(buf, &boolVal)
	if err != nil {
		return err
	}

	if boolVal {
		*b = TRUE
	} else {
		*b = FALSE
	}
	return nil
}

// MarshalHayson is the same as MarshalJSON
func (b Bool) MarshalHayson() ([]byte, error) {
	return json.Marshal(b.ToBool())
}
