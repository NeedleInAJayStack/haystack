package haystack

import (
	"math"
	"testing"
)

func TestNumber_ToZinc(t *testing.T) {
	number := NewNumber(100.457, "kWh")
	valTest_ToZinc(number, "100.457kWh", t)

	inf := NewNumber(math.Inf(1), "")
	valTest_ToZinc(inf, "INF", t)

	negInf := NewNumber(math.Inf(-1), "")
	valTest_ToZinc(negInf, "-INF", t)

	nan := NewNumber(math.NaN(), "")
	valTest_ToZinc(nan, "NaN", t)
}

func TestNumber_MarshalJSON(t *testing.T) {
	number := NewNumber(100.457, "kWh")
	valTest_MarshalJSON(number, "\"n:100.457 kWh\"", t)

	inf := NewNumber(math.Inf(1), "")
	valTest_MarshalJSON(inf, "\"n:INF\"", t)

	negInf := NewNumber(math.Inf(-1), "")
	valTest_MarshalJSON(negInf, "\"n:-INF\"", t)

	nan := NewNumber(math.NaN(), "")
	valTest_MarshalJSON(nan, "\"n:NaN\"", t)
}

func TestNumber_MarshalHAYSON(t *testing.T) {
	number := NewNumber(100.457, "kWh")
	valTest_MarshalHAYSON(number, "{\"_kind\":\"number\",\"val\":100.457,\"unit\":\"kWh\"}", t)

	inf := NewNumber(math.Inf(1), "")
	valTest_MarshalHAYSON(inf, "{\"_kind\":\"number\",\"val\":\"INF\"}", t)

	negInf := NewNumber(math.Inf(-1), "")
	valTest_MarshalHAYSON(negInf, "{\"_kind\":\"number\",\"val\":\"-INF\"}", t)

	nan := NewNumber(math.NaN(), "")
	valTest_MarshalHAYSON(nan, "{\"_kind\":\"number\",\"val\":\"NaN\"}", t)
}
