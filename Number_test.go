package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber_ToZinc(t *testing.T) {
	assert.Equal(t, NewNumber(100.457, "").ToZinc(), "100.457")
	assert.Equal(t, NewNumber(100.457, "kWh").ToZinc(), "100.457kWh")
	assert.Equal(t, Inf().ToZinc(), "INF")
	assert.Equal(t, NegInf().ToZinc(), "-INF")
	assert.Equal(t, NaN().ToZinc(), "NaN")
}

func TestNumber_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNumber(100.457, ""), "\"n:100.457\"", t)
	valTest_MarshalJSON(NewNumber(100.457, "kWh"), "\"n:100.457 kWh\"", t)
	valTest_MarshalJSON(Inf(), "\"n:INF\"", t)
	valTest_MarshalJSON(NegInf(), "\"n:-INF\"", t)
	valTest_MarshalJSON(NaN(), "\"n:NaN\"", t)
}

func TestNumber_UnmarshalJSON(t *testing.T) {
	var number Number
	json.Unmarshal([]byte("\"n:100.457\""), &number)
	assert.Equal(t, number, NewNumber(100.457, ""))

	var numberUnit Number
	json.Unmarshal([]byte("\"n:100.457 kWh\""), &numberUnit)
	assert.Equal(t, numberUnit, NewNumber(100.457, "kWh"))

	var inf Number
	json.Unmarshal([]byte("\"n:INF\""), &inf)
	assert.Equal(t, inf, Inf())

	var negInf Number
	json.Unmarshal([]byte("\"n:-INF\""), &negInf)
	assert.Equal(t, negInf, NegInf())

	var nan Number
	json.Unmarshal([]byte("\"n:NaN\""), &nan)
	assert.Equal(t, nan.ToZinc(), "NaN")
}

func TestNumber_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNumber(100.457, ""), "{\"_kind\":\"number\",\"val\":100.457}", t)
	valTest_MarshalHayson(NewNumber(100.457, "kWh"), "{\"_kind\":\"number\",\"val\":100.457,\"unit\":\"kWh\"}", t)
	valTest_MarshalHayson(Inf(), "{\"_kind\":\"number\",\"val\":\"INF\"}", t)
	valTest_MarshalHayson(NegInf(), "{\"_kind\":\"number\",\"val\":\"-INF\"}", t)
	valTest_MarshalHayson(NaN(), "{\"_kind\":\"number\",\"val\":\"NaN\"}", t)
}
