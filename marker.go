package haystack

// Marker is the value for a marker tag.
type Marker struct {
}

// NewMarker creates a new Marker object.
func NewMarker() Marker {
	return Marker{}
}

// ToZinc representes the object as: "M"
func (marker Marker) ToZinc() string {
	return "M"
}

// Remove is the value used to indicate a tag remove.
type Remove struct {
}

// NewRemove creates a new Remove object.
func NewRemove() Remove {
	return Remove{}
}

// ToZinc representes the object as: "R"
func (remove Remove) ToZinc() string {
	return "R"
}
