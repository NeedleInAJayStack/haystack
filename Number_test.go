package haystack

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber_ToZinc(t *testing.T) {
	assert.Equal(t, NewNumber(100.457, "").ToZinc(), "100.457")
	assert.Equal(t, NewNumber(100.457, "kWh").ToZinc(), "100.457kWh")
	assert.Equal(t, NewNumber(math.Inf(1), "").ToZinc(), "INF")
	assert.Equal(t, NewNumber(math.Inf(-1), "").ToZinc(), "-INF")
	assert.Equal(t, NewNumber(math.NaN(), "").ToZinc(), "NaN")
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
	json.Unmarshal([]byte("\"n:100.457\""), &number)
	assert.Equal(t, number, NewNumber(100.457, ""))

	var numberUnit Number
	json.Unmarshal([]byte("\"n:100.457 kWh\""), &numberUnit)
	assert.Equal(t, numberUnit, NewNumber(100.457, "kWh"))

	var inf Number
	json.Unmarshal([]byte("\"n:INF\""), &inf)
	assert.Equal(t, inf, NewNumber(math.Inf(1), ""))

	var negInf Number
	json.Unmarshal([]byte("\"n:-INF\""), &negInf)
	assert.Equal(t, negInf, NewNumber(math.Inf(-1), ""))

	var nan Number
	json.Unmarshal([]byte("\"n:NaN\""), &nan)
	assert.Equal(t, nan.ToZinc(), "NaN")
}

func TestNumber_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNumber(100.457, ""), "{\"_kind\":\"number\",\"val\":100.457}", t)
	valTest_MarshalHayson(NewNumber(100.457, "kWh"), "{\"_kind\":\"number\",\"val\":100.457,\"unit\":\"kWh\"}", t)
	valTest_MarshalHayson(NewNumber(math.Inf(1), ""), "{\"_kind\":\"number\",\"val\":\"INF\"}", t)
	valTest_MarshalHayson(NewNumber(math.Inf(-1), ""), "{\"_kind\":\"number\",\"val\":\"-INF\"}", t)
	valTest_MarshalHayson(NewNumber(math.NaN(), ""), "{\"_kind\":\"number\",\"val\":\"NaN\"}", t)
}
