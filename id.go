package haystack

type Id struct {
	val string
}

// ToZinc representes the object as: "<val>"
func (id Id) ToZinc() string {
	return id.val
}
