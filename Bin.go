package haystack

import (
	"encoding/json"
	"strings"
)

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

// MarshalHAYSON representes the object as: "{\"_kind\":\"bin\",\"mime\":\"<mime>\"}"
// This representation is unofficial, but it fills out the Val interface
func (bin Bin) MarshalHAYSON() ([]byte, error) {
	builder := new(strings.Builder)
	builder.WriteString("{\"_kind\":\"bin\",\"mime\":\"")
	builder.WriteString(bin.mime)
	builder.WriteString("\"}")
	return []byte(builder.String()), nil
}
