package haystack

type Marker struct {
}

func (marker *Marker) toZinc() string {
	return "M"
}
