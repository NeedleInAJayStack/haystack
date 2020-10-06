package haystack

type Bin struct {
	mime string
}

func NewBin(mime string) Bin {
	return Bin{mime: mime}
}

// ToZinc representes the object as: "{@code Bin("<mime>")}"
func (bin Bin) ToZinc() string {
	return "Bin(\"" + bin.mime + "\")"
}
