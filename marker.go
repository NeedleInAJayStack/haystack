package haystack

type Marker struct {
}

// Convert object to zinc.
func (marker *Marker) toZinc() string {
	return "M"
}
