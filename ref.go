package haystack

type Ref struct {
	val string
	dis string // Optional
}

// Format is "@<id> [dis]"
func (ref Ref) ToZinc() string {
	result := "@" + ref.val
	if ref.dis != "" {
		dis := Str{val: ref.dis}
		result = result + " " + dis.ToZinc()
	}
	return result
}
