package haystack

type NA struct {
}

func NewNA() NA {
	return NA{}
}

// ToZinc representes the object as: "NA"
func (na NA) ToZinc() string {
	return "NA"
}
