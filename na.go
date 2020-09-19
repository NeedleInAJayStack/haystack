package haystack

type NA struct {
}

func (na NA) ToZinc() string {
	return "NA"
}
