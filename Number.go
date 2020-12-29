package haystack

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Number wraps a 64-bit floating point number and unit name.
type Number struct {
	val  float64
	unit string
}

// NewNumber creates a new Number. For unitless numbers, use an empty string unit: ""
func NewNumber(val float64, unit string) Number {
	return Number{val: val, unit: unit}
}

// Float returns the numerical value
func (number Number) Float() float64 {
	return number.val
}

// Unit returns the unit symbol
func (number Number) Unit() string {
	return number.unit
}

// ToZinc representes the object as: "<val>[unit]"
func (number Number) ToZinc() string {
	return number.encode(false)
}

// MarshalJSON representes the object as: "n:<val> [unit]"
func (number Number) MarshalJSON() ([]byte, error) {
	return json.Marshal("n:" + number.encode(true))
}

// MarshalHayson representes the object as: "{"_kind":"ref","val":<id>,["dis":<dis>]}"
func (number Number) MarshalHayson() ([]byte, error) {
	buf := strings.Builder{}

	buf.WriteString("{\"_kind\":\"number\",\"val\":")
	if math.IsInf(number.val, 1) {
		buf.WriteString("\"INF\"")
	} else if math.IsInf(number.val, -1) {
		buf.WriteString("\"-INF\"")
	} else if math.IsNaN(number.val) {
		buf.WriteString("\"NaN\"")
	} else {
		buf.WriteString(fmt.Sprintf("%g", number.val))
	}
	if number.unit != "" {
		buf.WriteString(",\"unit\":\"")
		buf.WriteString(number.unit)
		buf.WriteString("\"")
	}
	buf.WriteString("}")
	return []byte(buf.String()), nil
}

func (number Number) encode(spaceBeforeUnit bool) string {
	if math.IsInf(number.val, 1) {
		return "INF"
	} else if math.IsInf(number.val, -1) {
		return "-INF"
	} else if math.IsNaN(number.val) {
		return "NaN"
	} else {
		result := fmt.Sprintf("%g", number.val)

		if number.unit != "" {
			if spaceBeforeUnit {
				result = result + " "
			}
			result = result + number.unit
		}
		return result
	}
}
