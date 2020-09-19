package haystack

import (
	"math"
	"testing"
)

func TestNumber_ToZinc(t *testing.T) {
	number := Number{val: 100.457, unit: "kWh"}
	numberStr := number.ToZinc()
	if numberStr != "100.457kWh" {
		t.Error(numberStr)
	}

	inf := Number{val: math.Inf(1)}
	infStr := inf.ToZinc()
	if infStr != "INF" {
		t.Error(infStr)
	}

	negInf := Number{val: math.Inf(-1)}
	negInfStr := negInf.ToZinc()
	if negInfStr != "-INF" {
		t.Error(negInfStr)
	}

	nan := Number{val: math.NaN()}
	nanStr := nan.ToZinc()
	if nanStr != "NaN" {
		t.Error(nanStr)
	}
}
