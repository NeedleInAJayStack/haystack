package io

import (
	"testing"
	"time"

	"github.com/NeedleInAJayStack/haystack"
	"github.com/stretchr/testify/assert"
)

func TestTokenizer_empty(t *testing.T) {
	testTokenizerMulti(t, "", []Token{}, []haystack.Val{}, false)
}
func TestTokenizer_testId(t *testing.T) {
	testTokenizerSingle(t, "x", ID, haystack.NewId("x"), false)
	testTokenizerSingle(t, "fooBar", ID, haystack.NewId("fooBar"), false)
	testTokenizerSingle(t, "fooBar1999x", ID, haystack.NewId("fooBar1999x"), false)
	testTokenizerSingle(t, "foo_23", ID, haystack.NewId("foo_23"), false)
	testTokenizerSingle(t, "Foo", ID, haystack.NewId("Foo"), false)
}
func TestTokenizer_testInts(t *testing.T) {
	testTokenizerSingle(t, "5", NUMBER, haystack.NewNumber(5, ""), false)
	testTokenizerSingle(t, "0x1234_abcd", NUMBER, haystack.NewNumber(0x1234_abcd, ""), false)
}
func TestTokenizer_testFloats(t *testing.T) {
	testTokenizerSingle(t, "5.0", NUMBER, haystack.NewNumber(5.0, ""), false)
	testTokenizerSingle(t, "5.42", NUMBER, haystack.NewNumber(5.42, ""), false)
	testTokenizerSingle(t, "123.2e32", NUMBER, haystack.NewNumber(123.2e32, ""), false)
	testTokenizerSingle(t, "123.2e+32", NUMBER, haystack.NewNumber(123.2e+32, ""), false)
	testTokenizerSingle(t, "2_123.2e+32", NUMBER, haystack.NewNumber(2_123.2e+32, ""), false)
	testTokenizerSingle(t, "4.2e-7", NUMBER, haystack.NewNumber(4.2e-7, ""), false)
}
func TestTokenizer_testNumberWithUnits(t *testing.T) {
	testTokenizerSingle(t, "-40ms", NUMBER, haystack.NewNumber(-40, "ms"), false)
	testTokenizerSingle(t, "1sec", NUMBER, haystack.NewNumber(1, "sec"), false)
	testTokenizerSingle(t, "2.5day", NUMBER, haystack.NewNumber(2.5, "day"), false)
	testTokenizerSingle(t, "12%", NUMBER, haystack.NewNumber(12, "%"), false)
	testTokenizerSingle(t, "987_foo", NUMBER, haystack.NewNumber(987, "_foo"), false)
	testTokenizerSingle(t, "-1.2m/s", NUMBER, haystack.NewNumber(-1.2, "m/s"), false)
	testTokenizerSingle(t, "12kWh/ft\u00B2", NUMBER, haystack.NewNumber(12, "kWh/ft\u00B2"), false)
	testTokenizerSingle(t, "3_000.5J/kg_dry", NUMBER, haystack.NewNumber(3000.5, "J/kg_dry"), false)
}
func TestTokenizer_testStr(t *testing.T) {
	testTokenizerSingle(t, "\"\"", STR, haystack.NewStr(""), false)
	testTokenizerSingle(t, "\"x y\"", STR, haystack.NewStr("x y"), false)
	testTokenizerSingle(t, "\"x\\\"y\"", STR, haystack.NewStr("x\"y"), false)
	testTokenizerSingle(t, "\"_\\u012f \\n \\t\\\" \\\\_\"", STR, haystack.NewStr("_\u012f \n \t\" \\_"), false)
}
func TestTokenizer_testDate(t *testing.T) {
	testTokenizerSingle(t, "2016-06-06", DATE, haystack.NewDate(2016, 6, 6), false)
}
func TestTokenizer_testTime(t *testing.T) {
	testTokenizerSingle(t, "8:30", TIME, haystack.NewTime(8, 30, 0, 0), false)
	testTokenizerSingle(t, "20:15", TIME, haystack.NewTime(20, 15, 0, 0), false)
	testTokenizerSingle(t, "00:00", TIME, haystack.NewTime(0, 0, 0, 0), false)
	testTokenizerSingle(t, "00:00:00", TIME, haystack.NewTime(0, 0, 0, 0), false)
	testTokenizerSingle(t, "01:02:03", TIME, haystack.NewTime(1, 2, 3, 0), false)
	testTokenizerSingle(t, "01:02:03", TIME, haystack.NewTime(1, 2, 3, 0), false)
	testTokenizerSingle(t, "23:59:59", TIME, haystack.NewTime(23, 59, 59, 0), false)
	testTokenizerSingle(t, "12:00:12.9", TIME, haystack.NewTime(12, 00, 12, 900), false)
	testTokenizerSingle(t, "12:00:12.9", TIME, haystack.NewTime(12, 00, 12, 900), false)
	testTokenizerSingle(t, "12:00:12.9", TIME, haystack.NewTime(12, 00, 12, 900), false)
	testTokenizerSingle(t, "12:00:12.99", TIME, haystack.NewTime(12, 00, 12, 990), false)
	testTokenizerSingle(t, "12:00:12.999", TIME, haystack.NewTime(12, 00, 12, 999), false)
	testTokenizerSingle(t, "12:00:12.000", TIME, haystack.NewTime(12, 00, 12, 0), false)
	testTokenizerSingle(t, "12:00:12.001", TIME, haystack.NewTime(12, 00, 12, 1), false)
}
func TestTokenizer_testDateTime(t *testing.T) {
	newYorkLocation, _ := time.LoadLocation("America/New_York")
	testTokenizerSingle(t, "2016-01-13T09:51:33-05:00 New_York", DATETIME,
		haystack.NewDateTimeFromGo(time.Date(2016, 1, 13, 9, 51, 33, 0, newYorkLocation)), true,
	)
	testTokenizerSingle(t, "2016-01-13T09:51:33.353-05:00 New_York", DATETIME,
		haystack.NewDateTimeFromGo(time.Date(2016, 1, 13, 9, 51, 33, 353e6, newYorkLocation)), true,
	)
	testTokenizerSingle(t, "2010-12-18T14:11:30.924Z", DATETIME,
		haystack.NewDateTimeFromGo(time.Date(2010, 12, 18, 14, 11, 30, 924e6, time.UTC)), true,
	)
	testTokenizerSingle(t, "2010-12-18T14:11:30.924Z UTC", DATETIME,
		haystack.NewDateTimeFromGo(time.Date(2010, 12, 18, 14, 11, 30, 924e6, time.UTC)), true,
	)
	// TODO: extract tzOffset from timezone name (go has no tz lookup)
	// testTokenizerSingle(t, "2010-12-18T14:11:30.924Z London", DATETIME,
	// 	haystack.NewDateTimeRaw(2010, 12, 18, 14, 11, 30, 924, 0, "London"),
	// )
	gmtPlus5Location, _ := time.LoadLocation("Etc/GMT+5")
	testTokenizerSingle(t, "2010-03-01T23:55:00.013-05:00 GMT+5", DATETIME,
		haystack.NewDateTimeFromGo(time.Date(2010, 3, 1, 23, 55, 00, 13e6, gmtPlus5Location)), true,
	)
	gmtMinus10Location, _ := time.LoadLocation("Etc/GMT-10")
	testTokenizerSingle(t, "2010-03-01T23:55:00.013+10:00 GMT-10", DATETIME,
		haystack.NewDateTimeFromGo(time.Date(2010, 3, 1, 23, 55, 00, 13e6, gmtMinus10Location)), true,
	)
}
func TestTokenizer_testRef(t *testing.T) {
	testTokenizerSingle(t, "@125b780e-0684e169", REF, haystack.NewRef("125b780e-0684e169", ""), false)
	testTokenizerSingle(t, "@demo:_:-.~", REF, haystack.NewRef("demo:_:-.~", ""), false)
}
func TestTokenizer_testUri(t *testing.T) {
	testTokenizerSingle(t, "`http://foo/`", URI, haystack.NewUri("http://foo/"), false)
	testTokenizerSingle(t, "`_ \\n \\\\ \\`_`", URI, haystack.NewUri("_ \n \\\\ `_"), false)
}
func TestTokenizer_testWhitespace(t *testing.T) {
	testTokenizerMulti(t, "a\n  b   \rc \r\nd\n\ne",
		[]Token{
			ID,
			NL,
			ID,
			NL,
			ID,
			NL,
			ID,
			NL,
			NL,
			ID,
		},
		[]haystack.Val{
			haystack.NewId("a"),
			haystack.NewNull(),
			haystack.NewId("b"),
			haystack.NewNull(),
			haystack.NewId("c"),
			haystack.NewNull(),
			haystack.NewId("d"),
			haystack.NewNull(),
			haystack.NewNull(),
			haystack.NewId("e"),
		},
		false,
	)
}

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'ToZinc' method
func testTokenizerSingle(t *testing.T, str string, expectedToken Token, expectedVal haystack.Val, compareByZinc bool) {
	testTokenizerMulti(t, str, []Token{expectedToken}, []haystack.Val{expectedVal}, compareByZinc)
}

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'ToZinc' method
func testTokenizerMulti(t *testing.T, str string, expectedTokens []Token, expectedVals []haystack.Val, compareByZinc bool) {
	tokens, vals := testTokenizerRead(t, str)

	assert.Equal(t, len(tokens), len(expectedTokens), "Actual and expected token list lengths don't match")
	for index, token := range tokens {
		assert.Equal(t, token, expectedTokens[index], "Tokens don't match\n"+
			"\tACTUAL:   "+token.String()+"\n"+
			"\tEXPECTED: "+expectedTokens[index].String())
	}

	assert.Equal(t, len(vals), len(expectedVals), "Actual and expected value list lengths don't match")
	for index, val := range vals {
		if compareByZinc {
			assert.Equal(t, val.ToZinc(), expectedVals[index].ToZinc(), "Val doesn't match expected\n"+
				"\tACTUAL:   "+val.ToZinc()+"\n"+
				"\tEXPECTED: "+expectedVals[index].ToZinc())
		} else {
			assert.Equal(t, val, expectedVals[index], "Val doesn't match expected\n"+
				"\tACTUAL:   "+val.ToZinc()+"\n"+
				"\tEXPECTED: "+expectedVals[index].ToZinc())
		}
	}
}

func testTokenizerRead(t *testing.T, str string) ([]Token, []haystack.Val) {
	var tokenizer Tokenizer
	tokenizer.InitString(str)

	var tokens []Token
	var vals []haystack.Val

	for {
		nextToken, err := tokenizer.Next()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, nextToken, tokenizer.token, "The same object doesn't equal itself")
		if nextToken == EOF {
			break
		} else {
			tokens = append(tokens, tokenizer.token)
			vals = append(vals, tokenizer.val)
		}
	}

	return tokens, vals
}
