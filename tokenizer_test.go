package haystack

import (
	"strings"
	"testing"
)

func TestTokenizer_empty(t *testing.T) {
	testTokenizerMulti(t, "", []Token{}, []Val{})
}
func TestTokenizer_testId(t *testing.T) {
	testTokenizerSingle(t, "x", ID, Id{val: "x"})
	testTokenizerSingle(t, "fooBar", ID, Id{val: "fooBar"})
	testTokenizerSingle(t, "fooBar1999x", ID, Id{val: "fooBar1999x"})
	testTokenizerSingle(t, "foo_23", ID, Id{val: "foo_23"})
	testTokenizerSingle(t, "Foo", ID, Id{val: "Foo"})
}
func TestTokenizer_testInts(t *testing.T) {
	testTokenizerSingle(t, "5", NUMBER, Number{val: 5})
	testTokenizerSingle(t, "0x1234_abcd", NUMBER, Number{val: 0x1234_abcd})
}
func TestTokenizer_testFloats(t *testing.T) {
	testTokenizerSingle(t, "5.0", NUMBER, Number{val: 5.0})
	testTokenizerSingle(t, "5.42", NUMBER, Number{val: 5.42})
	testTokenizerSingle(t, "123.2e32", NUMBER, Number{val: 123.2e32})
	testTokenizerSingle(t, "123.2e+32", NUMBER, Number{val: 123.2e+32})
	testTokenizerSingle(t, "2_123.2e+32", NUMBER, Number{val: 2_123.2e+32})
	testTokenizerSingle(t, "4.2e-7", NUMBER, Number{val: 4.2e-7})
}
func TestTokenizer_testNumberWithUnits(t *testing.T) {
	testTokenizerSingle(t, "-40ms", NUMBER, Number{val: -40, unit: "ms"})
	testTokenizerSingle(t, "1sec", NUMBER, Number{val: 1, unit: "sec"})
	testTokenizerSingle(t, "2.5day", NUMBER, Number{val: 2.5, unit: "day"})
	testTokenizerSingle(t, "12%", NUMBER, Number{val: 12, unit: "%"})
	testTokenizerSingle(t, "987_foo", NUMBER, Number{val: 987, unit: "_foo"})
	testTokenizerSingle(t, "-1.2m/s", NUMBER, Number{val: -1.2, unit: "m/s"})
	testTokenizerSingle(t, "12kWh/ft\u00B2", NUMBER, Number{val: 12, unit: "kWh/ft\u00B2"})
	testTokenizerSingle(t, "3_000.5J/kg_dry", NUMBER, Number{val: 3000.5, unit: "J/kg_dry"})
}
func TestTokenizer_testStr(t *testing.T) {
	testTokenizerSingle(t, "\"\"", STR, Str{val: ""})
	testTokenizerSingle(t, "\"x y\"", STR, Str{val: "x y"})
	testTokenizerSingle(t, "\"x\\\"y\"", STR, Str{val: "x\"y"})
	testTokenizerSingle(t, "\"_\\u012f \\n \\t \\\\_\"", STR, Str{val: "_\u012f \n \t \\_"})
}
func TestTokenizer_testDate(t *testing.T) {
	testTokenizerSingle(t, "2016-06-06", DATE, Date{year: 2016, month: 6, day: 6})
}
func TestTokenizer_testTime(t *testing.T) {
	testTokenizerSingle(t, "8:30", TIME, Time{hour: 8, min: 30})
	testTokenizerSingle(t, "20:15", TIME, Time{hour: 20, min: 15})
	testTokenizerSingle(t, "00:00", TIME, Time{hour: 0, min: 0})
	testTokenizerSingle(t, "00:00:00", TIME, Time{hour: 0, min: 0, sec: 0})
	testTokenizerSingle(t, "01:02:03", TIME, Time{hour: 1, min: 2, sec: 3})
	testTokenizerSingle(t, "01:02:03", TIME, Time{hour: 1, min: 2, sec: 3})
	testTokenizerSingle(t, "23:59:59", TIME, Time{hour: 23, min: 59, sec: 59})
	testTokenizerSingle(t, "12:00:12.9", TIME, Time{hour: 12, min: 00, sec: 12, ms: 900})
	testTokenizerSingle(t, "12:00:12.9", TIME, Time{hour: 12, min: 00, sec: 12, ms: 900})
	testTokenizerSingle(t, "12:00:12.9", TIME, Time{hour: 12, min: 00, sec: 12, ms: 900})
	testTokenizerSingle(t, "12:00:12.99", TIME, Time{hour: 12, min: 00, sec: 12, ms: 990})
	testTokenizerSingle(t, "12:00:12.999", TIME, Time{hour: 12, min: 00, sec: 12, ms: 999})
	testTokenizerSingle(t, "12:00:12.000", TIME, Time{hour: 12, min: 00, sec: 12, ms: 0})
	testTokenizerSingle(t, "12:00:12.001", TIME, Time{hour: 12, min: 00, sec: 12, ms: 1})
}
func TestTokenizer_testDateTime(t *testing.T) {
	testTokenizerSingle(t, "2016-01-13T09:51:33-05:00 New_York", DATETIME,
		DateTime{
			date:     Date{year: 2016, month: 1, day: 13},
			time:     Time{hour: 9, min: 51, sec: 33},
			tz:       "New_York",
			tzOffset: -18000,
		},
	)
	testTokenizerSingle(t, "2016-01-13T09:51:33.353-05:00 New_York", DATETIME,
		DateTime{
			date:     Date{year: 2016, month: 1, day: 13},
			time:     Time{hour: 9, min: 51, sec: 33, ms: 353},
			tz:       "New_York",
			tzOffset: -18000,
		},
	)
	testTokenizerSingle(t, "2010-12-18T14:11:30.924Z", DATETIME,
		DateTime{
			date:     Date{year: 2010, month: 12, day: 18},
			time:     Time{hour: 14, min: 11, sec: 30, ms: 924},
			tz:       "UTC",
			tzOffset: 0,
		},
	)
	testTokenizerSingle(t, "2010-12-18T14:11:30.924Z UTC", DATETIME,
		DateTime{
			date:     Date{year: 2010, month: 12, day: 18},
			time:     Time{hour: 14, min: 11, sec: 30, ms: 924},
			tz:       "UTC",
			tzOffset: 0,
		},
	)
	// TODO: extract tzOffset from timezone name (go has no tz lookup)
	// testTokenizerSingle(t, "2010-12-18T14:11:30.924Z London", DATETIME,
	//	 DateTime{
	// 		date: Date{year: 2010, month: 12, day: 18},
	// 		time: Time{hour: 14, min: 11, sec: 30, ms: 924},
	// 		tz: "London",
	// 		tzOffset: 0,
	//},
	// )
	testTokenizerSingle(t, "2010-03-01T23:55:00.013-05:00 GMT+5", DATETIME,
		DateTime{
			date:     Date{year: 2010, month: 3, day: 1},
			time:     Time{hour: 23, min: 55, sec: 00, ms: 13},
			tz:       "GMT+5",
			tzOffset: -18000,
		},
	)
	testTokenizerSingle(t, "2010-03-01T23:55:00.013+10:00 GMT-10 ", DATETIME,
		DateTime{
			date:     Date{year: 2010, month: 3, day: 1},
			time:     Time{hour: 23, min: 55, sec: 00, ms: 13},
			tz:       "GMT-10",
			tzOffset: 36000,
		},
	)
}
func TestTokenizer_testRef(t *testing.T) {
	testTokenizerSingle(t, "@125b780e-0684e169", REF, Ref{val: "125b780e-0684e169"})
	testTokenizerSingle(t, "@demo:_:-.~", REF, Ref{val: "demo:_:-.~"})
}
func TestTokenizer_testUri(t *testing.T) {
	testTokenizerSingle(t, "`http://foo/`", URI, Uri{val: "http://foo/"})
	testTokenizerSingle(t, "`_ \\n \\\\ \\`_`", URI, Uri{val: "_ \n \\\\ `_"})
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
		[]Val{
			Id{val: "a"},
			Null{},
			Id{val: "b"},
			Null{},
			Id{val: "c"},
			Null{},
			Id{val: "d"},
			Null{},
			Null{},
			Id{val: "e"},
		},
	)
}

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'ToZinc' method
func testTokenizerSingle(t *testing.T, str string, expectedToken Token, expectedVal Val) {
	testTokenizerMulti(t, str, []Token{expectedToken}, []Val{expectedVal})
}

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'ToZinc' method
func testTokenizerMulti(t *testing.T, str string, expectedTokens []Token, expectedVals []Val) {
	tokens, vals := testTokenizerRead(t, str)

	if len(tokens) != len(expectedTokens) {
		t.Error(str + " - Actual and expected token list lengths don't match")
	}
	for index, token := range tokens {
		if token != expectedTokens[index] {
			t.Error(str + " - Tokens don't match:\n" +
				"\tactual:   " + token.String() + "\n" +
				"\texpected: " + expectedTokens[index].String())
		}
	}

	if len(vals) != len(expectedVals) {
		t.Error(str + " - Actual and expected value list lengths don't match")
	}
	for index, val := range vals {
		if val.ToZinc() != expectedVals[index].ToZinc() {
			t.Error(str + " - Val doesn't match expected\n" +
				"\tactual:   " + val.ToZinc() + "\n" +
				"\texpected: " + expectedVals[index].ToZinc())
		}
	}
}

func testTokenizerRead(t *testing.T, str string) ([]Token, []Val) {
	var tokenizer Tokenizer
	tokenizer.Init(strings.NewReader(str))

	var tokens []Token
	var vals []Val

	// TODO: Adjust to handle multiple vals/tokens
	for {
		nextToken, err := tokenizer.Next()
		if err != nil {
			t.Error(err)
		}
		if nextToken != tokenizer.token {
			t.Error("The same object doesn't equal itself")
		}
		if nextToken == EOF {
			break
		} else {
			tokens = append(tokens, tokenizer.token)
			vals = append(vals, tokenizer.val)
		}
	}

	return tokens, vals
}
