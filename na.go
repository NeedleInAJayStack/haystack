package haystack

type NA struct {
}

func (na NA) toZinc() string {
	return "NA"
}
