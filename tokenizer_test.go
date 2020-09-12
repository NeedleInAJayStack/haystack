package haystack

import (
	"strings"
	"testing"
)

// func TestTokenizer_empty(t *testing.T) {
// 	verifyToks(t, "", nil, nil)
// }

func TestTokenizer_testId(t *testing.T) {
	testTokenizer(t, "x", tokenId(), Id{val: "x"})
	testTokenizer(t, "fooBar", tokenId(), Id{val: "fooBar"})
	testTokenizer(t, "fooBar1999x", tokenId(), Id{val: "fooBar1999x"})
	testTokenizer(t, "foo_23", tokenId(), Id{val: "foo_23"})
	testTokenizer(t, "Foo", tokenId(), Id{val: "Foo"})
}
func TestTokenizer_testInts(t *testing.T) {
	testTokenizer(t, "5", tokenNumber(), Number{val: 5})
	testTokenizer(t, "0x1234_abcd", tokenNumber(), Number{val: 0x1234_abcd})
}
func TestTokenizer_testFloats(t *testing.T) {
	testTokenizer(t, "5.0", tokenNumber(), Number{val: 5.0})
	testTokenizer(t, "5.42", tokenNumber(), Number{val: 5.42})
	testTokenizer(t, "123.2e32", tokenNumber(), Number{val: 123.2e32})
	testTokenizer(t, "123.2e+32", tokenNumber(), Number{val: 123.2e+32})
	testTokenizer(t, "2_123.2e+32", tokenNumber(), Number{val: 2_123.2e+32})
	testTokenizer(t, "4.2e-7", tokenNumber(), Number{val: 4.2e-7})
}
func TestTokenizer_testNumberWithUnits(t *testing.T) {
	testTokenizer(t, "-40ms", tokenNumber(), Number{val: -40, unit: "ms"})
	testTokenizer(t, "1sec", tokenNumber(), Number{val: 1, unit: "sec"})
	testTokenizer(t, "2.5day", tokenNumber(), Number{val: 2.5, unit: "day"})
	testTokenizer(t, "12%", tokenNumber(), Number{val: 12, unit: "%"})
	testTokenizer(t, "987_foo", tokenNumber(), Number{val: 987, unit: "_foo"})
	testTokenizer(t, "-1.2m/s", tokenNumber(), Number{val: -1.2, unit: "m/s"})
	testTokenizer(t, "12kWh/ft\u00B2", tokenNumber(), Number{val: 12, unit: "kWh/ft\u00B2"})
	testTokenizer(t, "3_000.5J/kg_dry", tokenNumber(), Number{val: 3000.5, unit: "J/kg_dry"})
}
func TestTokenizer_testStr(t *testing.T) {
	testTokenizer(t, "\"\"", tokenStr(), Str{val: ""})
	testTokenizer(t, "\"x y\"", tokenStr(), Str{val: "x y"})
	testTokenizer(t, "\"x\\\"y\"", tokenStr(), Str{val: "x\"y"})
	testTokenizer(t, "\"_\\u012f \\n \\t \\\\_\"", tokenStr(), Str{val: "_\u012f \n \t \\_"})
}
func TestTokenizer_testDate(t *testing.T) {
	testTokenizer(t, "2016-06-06", tokenDate(), Date{year: 2016, month: 6, day: 6})
}
func TestTokenizer_testTime(t *testing.T) {
	testTokenizer(t, "8:30", tokenTime(), Time{hour: 8, min: 30})
	testTokenizer(t, "20:15", tokenTime(), Time{hour: 20, min: 15})
	testTokenizer(t, "00:00", tokenTime(), Time{hour: 0, min: 0})
	testTokenizer(t, "00:00:00", tokenTime(), Time{hour: 0, min: 0, sec: 0})
	testTokenizer(t, "01:02:03", tokenTime(), Time{hour: 1, min: 2, sec: 3})
	testTokenizer(t, "01:02:03", tokenTime(), Time{hour: 1, min: 2, sec: 3})
	testTokenizer(t, "23:59:59", tokenTime(), Time{hour: 23, min: 59, sec: 59})
	testTokenizer(t, "12:00:12.9", tokenTime(), Time{hour: 12, min: 00, sec: 12, ms: 900})
	testTokenizer(t, "12:00:12.9", tokenTime(), Time{hour: 12, min: 00, sec: 12, ms: 900})
	testTokenizer(t, "12:00:12.9", tokenTime(), Time{hour: 12, min: 00, sec: 12, ms: 900})
	testTokenizer(t, "12:00:12.99", tokenTime(), Time{hour: 12, min: 00, sec: 12, ms: 990})
	testTokenizer(t, "12:00:12.999", tokenTime(), Time{hour: 12, min: 00, sec: 12, ms: 999})
	testTokenizer(t, "12:00:12.000", tokenTime(), Time{hour: 12, min: 00, sec: 12, ms: 0})
	testTokenizer(t, "12:00:12.001", tokenTime(), Time{hour: 12, min: 00, sec: 12, ms: 1})
}
func TestTokenizer_testDateTime(t *testing.T) {
	testTokenizer(t, "2016-01-13T09:51:33-05:00 New_York", tokenDateTime(),
		DateTime{
			date:     Date{year: 2016, month: 1, day: 13},
			time:     Time{hour: 9, min: 51, sec: 33},
			tz:       "New_York",
			tzOffset: -18000,
		},
	)
	testTokenizer(t, "2016-01-13T09:51:33.353-05:00 New_York", tokenDateTime(),
		DateTime{
			date:     Date{year: 2016, month: 1, day: 13},
			time:     Time{hour: 9, min: 51, sec: 33, ms: 353},
			tz:       "New_York",
			tzOffset: -18000,
		},
	)
	testTokenizer(t, "2010-12-18T14:11:30.924Z", tokenDateTime(),
		DateTime{
			date:     Date{year: 2010, month: 12, day: 18},
			time:     Time{hour: 14, min: 11, sec: 30, ms: 924},
			tz:       "UTC",
			tzOffset: 0,
		},
	)
	testTokenizer(t, "2010-12-18T14:11:30.924Z UTC", tokenDateTime(),
		DateTime{
			date:     Date{year: 2010, month: 12, day: 18},
			time:     Time{hour: 14, min: 11, sec: 30, ms: 924},
			tz:       "UTC",
			tzOffset: 0,
		},
	)
	// TODO: extract tzOffset from timezone name (go has no tz lookup)
	// testTokenizer(t, "2010-12-18T14:11:30.924Z London", tokenDateTime(),
	// 	DateTime{
	// 		date: Date{year: 2010, month: 12, day: 18},
	// 		time: Time{hour: 14, min: 11, sec: 30, ms: 924},
	// 		tz: "London",
	// 		tzOffset: 0,
	// 	},
	// )
	testTokenizer(t, "2010-03-01T23:55:00.013-05:00 GMT+5", tokenDateTime(),
		DateTime{
			date:     Date{year: 2010, month: 3, day: 1},
			time:     Time{hour: 23, min: 55, sec: 00, ms: 13},
			tz:       "GMT+5",
			tzOffset: -18000,
		},
	)
	testTokenizer(t, "2010-03-01T23:55:00.013+10:00 GMT-10 ", tokenDateTime(),
		DateTime{
			date:     Date{year: 2010, month: 3, day: 1},
			time:     Time{hour: 23, min: 55, sec: 00, ms: 13},
			tz:       "GMT-10",
			tzOffset: 36000,
		},
	)
}

// TODO testRef
// TODO testUri
// TODO testWhitespace

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'toZinc' method
func testTokenizer(t *testing.T, str string, expectedToken Token, expectedVal Val) {
	var token Token
	var val Val

	tokenizer := newTokenizer(strings.NewReader(str))

	// TODO: Adjust to handle multiple vals/tokens
	for {
		nextToken, err := tokenizer.Next()
		if err != nil {
			t.Error(err)
		}
		if nextToken != tokenizer.token {
			t.Error("The same object doesn't equal itself")
		}
		if nextToken == tokenEof() {
			break
		} else {
			token = tokenizer.token
			val = tokenizer.val
		}
	}

	if !token.equals(expectedToken) {
		t.Error(str + " - Tokens don't match:\n" +
			"\tactual:   " + token.symbol + "\n" +
			"\texpected: " + expectedToken.symbol)
	}
	if val.toZinc() != expectedVal.toZinc() {
		t.Error(str + " - Val doesn't match expected\n" +
			"\tactual:   " + val.toZinc() + "\n" +
			"\texpected: " + expectedVal.toZinc())
	}
}
