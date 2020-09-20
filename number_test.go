package haystack

import (
	"math"
	"testing"
)

func TestNumber_ToZinc(t *testing.T) {
	number := NewNumber(100.457, "kWh")
	numberStr := number.ToZinc()
	if numberStr != "100.457kWh" {
		t.Error(numberStr)
	}

	inf := NewNumber(math.Inf(1), "")
	infStr := inf.ToZinc()
	if infStr != "INF" {
		t.Error(infStr)
	}

	negInf := NewNumber(math.Inf(-1), "")
	negInfStr := negInf.ToZinc()
	if negInfStr != "-INF" {
		t.Error(negInfStr)
	}

	nan := NewNumber(math.NaN(), "")
	nanStr := nan.ToZinc()
	if nanStr != "NaN" {
		t.Error(nanStr)
	}
}
