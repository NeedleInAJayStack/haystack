package haystack

// Bool models a boolean for true/false tag values.
type Bool struct {
	val bool
}

// NewBool creates a new Bool object.
func NewBool(val bool) Bool {
	return Bool{val: val}
}

// ToBool returns the value of this object as a Go bool
func (b Bool) ToBool() bool {
	return b.val
}

// ToZinc representes the object as: "T" or "F"
func (b Bool) ToZinc() string {
	if b.val {
		return "T"
	}
	return "F"
}
