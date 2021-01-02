package haystack

import (
	"encoding/json"
	"errors"
	"strings"
)

// Bin models a binary file with a MIME type.
type Bin struct {
	mime string
}

// NewBin creates a new Bin object.
func NewBin(mime string) *Bin {
	return &Bin{mime: mime}
}

// ToZinc representes the object as: "{@code Bin("<mime>")}"
func (bin *Bin) ToZinc() string {
	return "Bin(\"" + bin.mime + "\")"
}

// MarshalJSON representes the object as: "b:<mime>"
func (bin *Bin) MarshalJSON() ([]byte, error) {
	return json.Marshal("b:" + bin.mime)
}

// UnmarshalJSON interprets the json value: "b:<val>"
func (bin *Bin) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newBin, newErr := binFromJSON(jsonStr)
	*bin = *newBin
	return newErr
}

func binFromJSON(jsonStr string) (*Bin, error) {
	if !strings.HasPrefix(jsonStr, "b:") {
		return nil, errors.New("Input value does not begin with b:")
	}
	mime := jsonStr[2:len(jsonStr)]
	return NewBin(mime), nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"bin\",\"mime\":\"<mime>\"}"
// This representation is unofficial, but it fills out the Val interface
func (bin *Bin) MarshalHayson() ([]byte, error) {
	builder := new(strings.Builder)
	builder.WriteString("{\"_kind\":\"bin\",\"mime\":\"")
	builder.WriteString(bin.mime)
	builder.WriteString("\"}")
	return []byte(builder.String()), nil
}
