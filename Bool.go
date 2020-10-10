package haystack

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

// ToJSON representes the object as: "true" or "false"
func (b Bool) ToJSON() string {
	if b.ToBool() {
		return "true"
	}
	return "false"
}