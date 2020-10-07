package haystack

import (
	"fmt"
	"math"
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
