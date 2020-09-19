package haystack

type Bool struct {
	val bool
}

func (b Bool) ToZinc() string {
	if b.val == true {
		return "T"
	} else {
		return "F"
	}
}
