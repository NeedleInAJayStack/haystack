package haystack

type Marker struct {
}

func (marker Marker) toZinc() string {
	return "M"
}

type Remove struct {
}

func (remove Remove) toZinc() string {
	return "R"
}
