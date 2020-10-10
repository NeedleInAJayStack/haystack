package haystack

// Ref wraps a string reference identifier and display name.
type Ref struct {
	id  string
	dis string // Optional
}

// NewRef creates a new Ref. For display-less refs, use an empty string dis: ""
func NewRef(id string, dis string) Ref {
	return Ref{id: id, dis: dis}
}

// Id returns the ref identifier
func (ref Ref) Id() string {
	return ref.id
}

// Dis returns the ref display string
func (ref Ref) Dis() string {
	return ref.dis
}

// ToZinc representes the object as: "@<id> \"[dis]\""
func (ref Ref) ToZinc() string {
	result := "@" + ref.id
	if ref.dis != "" {
		dis := Str{val: ref.dis}
		result = result + " " + dis.ToZinc()
	}
	return result
}

// ToJSON representes the object as: "r:<id> [dis]"
func (ref Ref) ToJSON() string {
	result := "r:" + ref.id
	if ref.dis != "" {
		result = result + " " + ref.dis
	}
	return result
}
