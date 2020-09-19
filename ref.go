package haystack

type Ref struct {
	val string
	dis string // Optional
}

func NewRef(val string, dis string) Ref {
	return Ref{val: val, dis: dis}
}

// ToZinc representes the object as: "@<id> [dis]"
func (ref Ref) ToZinc() string {
	result := "@" + ref.val
	if ref.dis != "" {
		dis := Str{val: ref.dis}
		result = result + " " + dis.ToZinc()
	}
	return result
}
