package haystack

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
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

// Inf returns a positive infinity number
func Inf() Number {
	return NewNumber(math.Inf(1), "")
}

// NegInf returns a negative infinity number
func NegInf() Number {
	return NewNumber(math.Inf(-1), "")
}

// NaN returns a not-a-number number
func NaN() Number {
	return NewNumber(math.NaN(), "")
}

// newNumberFromStr creates a new number from a string. The string representation must have a space between the number and unit
func newNumberFromStr(str string) (Number, error) {
	if str == "INF" {
		return Inf(), nil
	} else if str == "-INF" {
		return NegInf(), nil
	} else if str == "NaN" {
		return NaN(), nil
	} else {
		numberSplit := strings.Split(str, " ")
		val, valErr := strconv.ParseFloat(numberSplit[0], 64)
		if valErr != nil {
			return NewNumber(0.0, ""), valErr
		}
		unit := ""
		if len(numberSplit) > 1 {
			unit = numberSplit[1]
		}
		return NewNumber(val, unit), nil
	}
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
	return number.toStr(false)
}

// MarshalJSON representes the object as: "n:<val> [unit]"
func (number Number) MarshalJSON() ([]byte, error) {
	return json.Marshal("n:" + number.toStr(true))
}

// UnmarshalJSON interprets the json value: "n:<val> [unit]"
func (number *Number) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newNumber, newErr := numberFromJSON(jsonStr)
	*number = newNumber
	return newErr
}

func numberFromJSON(jsonStr string) (Number, error) {
	if !strings.HasPrefix(jsonStr, "n:") {
		return Number{}, errors.New("value does not begin with 'n:'")
	}
	numberStr := jsonStr[2:]

	return newNumberFromStr(numberStr)
}

// MarshalHayson representes the object as: "{"_kind":"number","val":<val>,["unit":<unit>]}"
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

func (number Number) toStr(spaceBeforeUnit bool) string {
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
