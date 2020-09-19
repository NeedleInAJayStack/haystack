package haystack

type Marker struct {
}

func (marker Marker) ToZinc() string {
	return "M"
}

type Remove struct {
}

func (remove Remove) ToZinc() string {
	return "R"
}
