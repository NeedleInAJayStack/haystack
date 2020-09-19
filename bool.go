package haystack

type Bool struct {
	val bool
}

// ToZinc representes the object as: "T" or "F"
func (b Bool) ToZinc() string {
	if b.val {
		return "T"
	}
	return "F"
}
