package haystack

type NA struct {
}

func (na NA) toZinc() string {
	return "NA"
}

func (na1 NA) equals(na2 NA) bool {
	return true
}
