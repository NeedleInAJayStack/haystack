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

func TestNumber_ToJSON(t *testing.T) {
	number := NewNumber(100.457, "kWh")
	numberStr := number.ToJSON()
	if numberStr != "n:100.457 kWh" {
		t.Error(numberStr)
	}

	inf := NewNumber(math.Inf(1), "")
	infStr := inf.ToJSON()
	if infStr != "n:INF" {
		t.Error(infStr)
	}

	negInf := NewNumber(math.Inf(-1), "")
	negInfStr := negInf.ToJSON()
	if negInfStr != "n:-INF" {
		t.Error(negInfStr)
	}

	nan := NewNumber(math.NaN(), "")
	nanStr := nan.ToJSON()
	if nanStr != "n:NaN" {
		t.Error(nanStr)
	}
}
