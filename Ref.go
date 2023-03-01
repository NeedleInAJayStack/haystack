package haystack

import (
	"encoding/json"
	"errors"
	"strings"
)

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

// MarshalJSON representes the object as: "r:<id> [dis]"
func (ref Ref) MarshalJSON() ([]byte, error) {
	result := "r:" + ref.id
	if ref.dis != "" {
		result = result + " " + ref.dis
	}
	return json.Marshal(result)
}

// UnmarshalJSON interprets the json value: "r:<id> [dis]"
func (ref *Ref) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newRef, newErr := refFromJSON(jsonStr)
	*ref = newRef
	return newErr
}

func refFromJSON(jsonStr string) (Ref, error) {
	if !strings.HasPrefix(jsonStr, "r:") {
		return Ref{}, errors.New("value does not begin with 'r:'")
	}
	refStr := jsonStr[2:]
	firstSpaceIndex := strings.Index(refStr, " ")

	if firstSpaceIndex == -1 {
		return NewRef(refStr, ""), nil
	} else {
		id := refStr[:firstSpaceIndex]
		dis := refStr[firstSpaceIndex+1:]
		return NewRef(id, dis), nil
	}
}

// MarshalHayson representes the object as: "{"_kind":"ref","val":<id>,["dis":<dis>]}"
func (ref Ref) MarshalHayson() ([]byte, error) {
	buf := strings.Builder{}

	buf.WriteString("{\"_kind\":\"ref\",\"val\":\"")
	buf.WriteString(ref.id)
	buf.WriteString("\"")
	if ref.dis != "" {
		buf.WriteString(",\"dis\":\"")
		buf.WriteString(ref.dis)
		buf.WriteString("\"")
	}
	buf.WriteString("}")
	return []byte(buf.String()), nil
}
