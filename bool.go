package haystack

type Bool struct {
	val bool
}

func NewBool(val bool) Bool {
	return Bool{val: val}
}

// ToZinc representes the object as: "T" or "F"
func (b Bool) ToZinc() string {
	if b.val {
		return "T"
	}
	return "F"
}
