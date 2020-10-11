package haystack

import "encoding/json"

// Bin models a binary file with a MIME type.
type Bin struct {
	mime string
}

// NewBin creates a new Bin object.
func NewBin(mime string) Bin {
	return Bin{mime: mime}
}

// ToZinc representes the object as: "{@code Bin("<mime>")}"
func (bin Bin) ToZinc() string {
	return "Bin(\"" + bin.mime + "\")"
}

// MarshalJSON representes the object as: "b:<mime>"
func (bin Bin) MarshalJSON() ([]byte, error) {
	return json.Marshal("b:" + bin.mime)
}