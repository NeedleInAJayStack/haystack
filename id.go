package haystack

type Id struct {
	val string
}

// Encode using double quotes and back slash escapes
func (id Id) toZinc() string {
	return id.val
}
