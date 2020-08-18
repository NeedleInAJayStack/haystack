package haystack

// A ref is composed of an id and an optional dis
type Ref struct {
	val string
	dis string
}

// Convert object to zinc. Format is "@<id> [dis]"
func (ref *Ref) toZinc() string {
	result := "@" + ref.val
	if ref.dis != "" {
		dis := Str{val: ref.dis}
		result = result + " " + dis.toZinc()
	}
	return result
}
