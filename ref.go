package haystack

type Ref struct {
	id  string
	dis string // Optional
}

func NewRef(id string, dis string) Ref {
	return Ref{id: id, dis: dis}
}

func (ref Ref) Id() string {
	return ref.id
}

func (ref Ref) Dis() string {
	return ref.dis
}

// ToZinc representes the object as: "@<id> [dis]"
func (ref Ref) ToZinc() string {
	result := "@" + ref.id
	if ref.dis != "" {
		dis := Str{val: ref.dis}
		result = result + " " + dis.ToZinc()
	}
	return result
}
