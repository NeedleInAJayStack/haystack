package haystack

import (
	"bufio"
	"encoding/json"
	"errors"
	"strings"
)

// Symbol is a name to a def in the meta-model namespace.
type Symbol struct {
	val string
}

// NewSymbol creates a new Symbol object
func NewSymbol(val string) Symbol {
	return Symbol{val: val}
}

// String returns the object's value directly as a Go string
func (symbol Symbol) String() string {
	return symbol.val
}

// MarshalJSON representes the object as "y:<val>"
func (symbol Symbol) MarshalJSON() ([]byte, error) {
	return json.Marshal("y:" + symbol.val)
}

// UnmarshalJSON interprets the json value: "y:<val>"
func (symbol *Symbol) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newSymbol, newErr := symbolFromJSON(jsonStr)
	*symbol = newSymbol
	return newErr
}

func symbolFromJSON(jsonStr string) (Symbol, error) {
	if !strings.HasPrefix(jsonStr, "y:") {
		return Symbol{}, errors.New("value does not begin with 'y:'")
	}
	val := jsonStr[2:]

	return NewSymbol(val), nil
}

// MarshalHayson representes the object as "{\"_kind\":\"symbol\",\"val\":\"<val>\"}"
func (symbol Symbol) MarshalHayson() ([]byte, error) {
	builder := new(strings.Builder)
	builder.WriteString("{\"_kind\":\"symbol\",\"val\":\"")
	builder.WriteString(symbol.val)
	builder.WriteString("\"}")
	return []byte(builder.String()), nil
}

// ToZinc represents the symbol with a prefix `^`
func (symbol Symbol) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	symbol.WriteZincTo(out)
	out.Flush()
	return builder.String()
}

// WriteZincTo writes the symbol with a prefix `^`
func (symbol Symbol) WriteZincTo(buf *bufio.Writer) {
	buf.WriteRune('^')
	buf.WriteString(symbol.val)
}
