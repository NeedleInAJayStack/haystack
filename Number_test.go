package haystack

import (
	"math"
	"testing"
)

func TestNumber_ToZinc(t *testing.T) {
	valTest_ToZinc(NewNumber(100.457, ""), "100.457", t)
	valTest_ToZinc(NewNumber(100.457, "kWh"), "100.457kWh", t)
	valTest_ToZinc(NewNumber(math.Inf(1), ""), "INF", t)
	valTest_ToZinc(NewNumber(math.Inf(-1), ""), "-INF", t)
	valTest_ToZinc(NewNumber(math.NaN(), ""), "NaN", t)
}

func TestNumber_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNumber(100.457, ""), "\"n:100.457\"", t)
	valTest_MarshalJSON(NewNumber(100.457, "kWh"), "\"n:100.457 kWh\"", t)
	valTest_MarshalJSON(NewNumber(math.Inf(1), ""), "\"n:INF\"", t)
	valTest_MarshalJSON(NewNumber(math.Inf(-1), ""), "\"n:-INF\"", t)
	valTest_MarshalJSON(NewNumber(math.NaN(), ""), "\"n:NaN\"", t)
}

func TestNumber_UnmarshalJSON(t *testing.T) {
	var number Number
	valTest_UnmarshalJSON("\"n:100.457\"", number, "100.457", t)
	var numberUnit Number
	valTest_UnmarshalJSON("\"n:100.457 kWh\"", numberUnit, "100.457kWh", t)
	var inf Number
	valTest_UnmarshalJSON("\"n:INF\"", inf, "INF", t)
	var negInf Number
	valTest_UnmarshalJSON("\"n:-INF\"", negInf, "-INF", t)
	var nan Number
	valTest_UnmarshalJSON("\"n:NaN\"", nan, "NaN", t)
}

func TestNumber_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNumber(100.457, ""), "{\"_kind\":\"number\",\"val\":100.457}", t)
	valTest_MarshalHayson(NewNumber(100.457, "kWh"), "{\"_kind\":\"number\",\"val\":100.457,\"unit\":\"kWh\"}", t)
	valTest_MarshalHayson(NewNumber(math.Inf(1), ""), "{\"_kind\":\"number\",\"val\":\"INF\"}", t)
	valTest_MarshalHayson(NewNumber(math.Inf(-1), ""), "{\"_kind\":\"number\",\"val\":\"-INF\"}", t)
	valTest_MarshalHayson(NewNumber(math.NaN(), ""), "{\"_kind\":\"number\",\"val\":\"NaN\"}", t)
}
