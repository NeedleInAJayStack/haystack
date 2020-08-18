package haystack

import (
	"math"
	"testing"
)

func TestNumber_toZinc(t *testing.T) {
	number := Number{val: 100.457, unit: "kWh"}
	numberStr := number.toZinc()
	if numberStr != "n:100.457kWh" {
		t.Error(numberStr)
	}

	inf := Number{val: math.Inf(1)}
	infStr := inf.toZinc()
	if infStr != "n:INF" {
		t.Error(infStr)
	}

	negInf := Number{val: math.Inf(-1)}
	negInfStr := negInf.toZinc()
	if negInfStr != "n:-INF" {
		t.Error(negInfStr)
	}

	nan := Number{val: math.NaN()}
	nanStr := nan.toZinc()
	if nanStr != "n:NaN" {
		t.Error(nanStr)
	}
}
