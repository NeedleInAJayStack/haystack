package haystack

type Marker struct {
}

// ToZinc representes the object as: "M"
func (marker Marker) ToZinc() string {
	return "M"
}

type Remove struct {
}

// ToZinc representes the object as: "R"
func (remove Remove) ToZinc() string {
	return "R"
}
