package haystack

type Marker struct {
}

func NewMarker() Marker {
	return Marker{}
}

// ToZinc representes the object as: "M"
func (marker Marker) ToZinc() string {
	return "M"
}

type Remove struct {
}

func NewRemove() Remove {
	return Remove{}
}

// ToZinc representes the object as: "R"
func (remove Remove) ToZinc() string {
	return "R"
}
