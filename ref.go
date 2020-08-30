package haystack

type Ref struct {
	val string
	dis string // Optional
}

// Format is "@<id> [dis]"
func (ref *Ref) toZinc() string {
	result := "@" + ref.val
	if ref.dis != "" {
		dis := Str{val: ref.dis}
		result = result + " " + dis.toZinc()
	}
	return result
}
