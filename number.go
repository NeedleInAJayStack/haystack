package haystack

import (
	"fmt"
	"math"
)

type Number struct {
	val  float64
	unit string // Optional
}

func NewNumber(val float64, unit string) Number {
	return Number{val: val, unit: unit}
}

func (number Number) Float() float64 {
	return number.val
}

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
