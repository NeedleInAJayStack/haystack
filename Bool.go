package haystack

import "encoding/json"

// Bool models a boolean for true/false tag values.
type Bool struct {
	val bool
}

// NewBool creates a new Bool object.
func NewBool(val bool) Bool {
	return Bool{
		val: val,
	}
}

// ToBool returns the value of this object as a Go bool
func (b Bool) ToBool() bool {
	return b.val
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
	var jsonBool bool
	err := json.Unmarshal(buf, &jsonBool)
	if err != nil {
		return err
	}
	*b = Bool{val: jsonBool}
	return nil
}

func boolFromJSON(jsonBool bool) (*Bool, error) {
	return &Bool{val: jsonBool}, nil
}

// MarshalHayson is the same as MarshalJSON
func (b Bool) MarshalHayson() ([]byte, error) {
	return json.Marshal(b.ToBool())
}
