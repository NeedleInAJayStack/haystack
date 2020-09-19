package haystack

type NA struct {
}

// ToZinc representes the object as: "NA"
func (na NA) ToZinc() string {
	return "NA"
}
