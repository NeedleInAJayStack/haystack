package haystack

type Id struct {
	val string
}

func NewId(val string) Id {
	return Id{val: val}
}

func (id Id) String() string {
	return id.val
}

// ToZinc representes the object as: "<val>"
func (id Id) ToZinc() string {
	return id.val
}
