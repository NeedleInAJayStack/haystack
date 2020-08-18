package haystack

// A number is composed of a val and unit
type Bool struct {
	val bool
}

// Convert object to zinc
func (b *Bool) toZinc() string {
	if b.val == true {
		return "T"
	} else {
		return "F"
	}
}
