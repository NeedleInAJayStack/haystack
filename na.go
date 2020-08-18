package haystack

type NA struct {
}

// Convert object to zinc.
func (na *NA) toZinc() string {
	return "NA"
}
