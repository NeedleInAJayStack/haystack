package haystack

import (
	"testing"
)

func TestTokenizer_empty(t *testing.T) {
	testTokenizerMulti(t, "", []Token{}, []Val{})
}
func TestTokenizer_testId(t *testing.T) {
	testTokenizerSingle(t, "x", ID, NewId("x"))
	testTokenizerSingle(t, "fooBar", ID, NewId("fooBar"))
	testTokenizerSingle(t, "fooBar1999x", ID, NewId("fooBar1999x"))
	testTokenizerSingle(t, "foo_23", ID, NewId("foo_23"))
	testTokenizerSingle(t, "Foo", ID, NewId("Foo"))
}
func TestTokenizer_testInts(t *testing.T) {
	testTokenizerSingle(t, "5", NUMBER, NewNumber(5, ""))
	testTokenizerSingle(t, "0x1234_abcd", NUMBER, NewNumber(0x1234_abcd, ""))
}
func TestTokenizer_testFloats(t *testing.T) {
	testTokenizerSingle(t, "5.0", NUMBER, NewNumber(5.0, ""))
	testTokenizerSingle(t, "5.42", NUMBER, NewNumber(5.42, ""))
	testTokenizerSingle(t, "123.2e32", NUMBER, NewNumber(123.2e32, ""))
	testTokenizerSingle(t, "123.2e+32", NUMBER, NewNumber(123.2e+32, ""))
	testTokenizerSingle(t, "2_123.2e+32", NUMBER, NewNumber(2_123.2e+32, ""))
	testTokenizerSingle(t, "4.2e-7", NUMBER, NewNumber(4.2e-7, ""))
}
func TestTokenizer_testNumberWithUnits(t *testing.T) {
	testTokenizerSingle(t, "-40ms", NUMBER, NewNumber(-40, "ms"))
	testTokenizerSingle(t, "1sec", NUMBER, NewNumber(1, "sec"))
	testTokenizerSingle(t, "2.5day", NUMBER, NewNumber(2.5, "day"))
	testTokenizerSingle(t, "12%", NUMBER, NewNumber(12, "%"))
	testTokenizerSingle(t, "987_foo", NUMBER, NewNumber(987, "_foo"))
	testTokenizerSingle(t, "-1.2m/s", NUMBER, NewNumber(-1.2, "m/s"))
	testTokenizerSingle(t, "12kWh/ft\u00B2", NUMBER, NewNumber(12, "kWh/ft\u00B2"))
	testTokenizerSingle(t, "3_000.5J/kg_dry", NUMBER, NewNumber(3000.5, "J/kg_dry"))
}
func TestTokenizer_testStr(t *testing.T) {
	testTokenizerSingle(t, "\"\"", STR, NewStr(""))
	testTokenizerSingle(t, "\"x y\"", STR, NewStr("x y"))
	testTokenizerSingle(t, "\"x\\\"y\"", STR, NewStr("x\"y"))
	testTokenizerSingle(t, "\"_\\u012f \\n \\t\\\" \\\\_\"", STR, NewStr("_\u012f \n \t\" \\_"))
}
func TestTokenizer_testDate(t *testing.T) {
	testTokenizerSingle(t, "2016-06-06", DATE, NewDate(2016, 6, 6))
}
func TestTokenizer_testTime(t *testing.T) {
	testTokenizerSingle(t, "8:30", TIME, NewTime(8, 30, 0, 0))
	testTokenizerSingle(t, "20:15", TIME, NewTime(20, 15, 0, 0))
	testTokenizerSingle(t, "00:00", TIME, NewTime(0, 0, 0, 0))
	testTokenizerSingle(t, "00:00:00", TIME, NewTime(0, 0, 0, 0))
	testTokenizerSingle(t, "01:02:03", TIME, NewTime(1, 2, 3, 0))
	testTokenizerSingle(t, "01:02:03", TIME, NewTime(1, 2, 3, 0))
	testTokenizerSingle(t, "23:59:59", TIME, NewTime(23, 59, 59, 0))
	testTokenizerSingle(t, "12:00:12.9", TIME, NewTime(12, 00, 12, 900))
	testTokenizerSingle(t, "12:00:12.9", TIME, NewTime(12, 00, 12, 900))
	testTokenizerSingle(t, "12:00:12.9", TIME, NewTime(12, 00, 12, 900))
	testTokenizerSingle(t, "12:00:12.99", TIME, NewTime(12, 00, 12, 990))
	testTokenizerSingle(t, "12:00:12.999", TIME, NewTime(12, 00, 12, 999))
	testTokenizerSingle(t, "12:00:12.000", TIME, NewTime(12, 00, 12, 0))
	testTokenizerSingle(t, "12:00:12.001", TIME, NewTime(12, 00, 12, 1))
}
func TestTokenizer_testDateTime(t *testing.T) {
	testTokenizerSingle(t, "2016-01-13T09:51:33-05:00 New_York", DATETIME,
		NewDateTime(2016, 1, 13, 9, 51, 33, 0, -18000, "New_York"),
	)
	testTokenizerSingle(t, "2016-01-13T09:51:33.353-05:00 New_York", DATETIME,
		NewDateTime(2016, 1, 13, 9, 51, 33, 353, -18000, "New_York"),
	)
	testTokenizerSingle(t, "2010-12-18T14:11:30.924Z", DATETIME,
		NewDateTime(2010, 12, 18, 14, 11, 30, 924, 0, "UTC"),
	)
	testTokenizerSingle(t, "2010-12-18T14:11:30.924Z UTC", DATETIME,
		NewDateTime(2010, 12, 18, 14, 11, 30, 924, 0, "UTC"),
	)
	// TODO: extract tzOffset from timezone name (go has no tz lookup)
	// testTokenizerSingle(t, "2010-12-18T14:11:30.924Z London", DATETIME,
	//	 NewDateTime(2010, 12, 18, 14, 11, 30, 924, 0, "London"),
	// )
	testTokenizerSingle(t, "2010-03-01T23:55:00.013-05:00 GMT+5", DATETIME,
		NewDateTime(2010, 3, 1, 23, 55, 00, 13, -18000, "GMT+5"),
	)
	testTokenizerSingle(t, "2010-03-01T23:55:00.013+10:00 GMT-10 ", DATETIME,
		NewDateTime(2010, 3, 1, 23, 55, 00, 13, 36000, "GMT-10"),
	)
}
func TestTokenizer_testRef(t *testing.T) {
	testTokenizerSingle(t, "@125b780e-0684e169", REF, NewRef("125b780e-0684e169", ""))
	testTokenizerSingle(t, "@demo:_:-.~", REF, NewRef("demo:_:-.~", ""))
}
func TestTokenizer_testUri(t *testing.T) {
	testTokenizerSingle(t, "`http://foo/`", URI, NewUri("http://foo/"))
	testTokenizerSingle(t, "`_ \\n \\\\ \\`_`", URI, NewUri("_ \n \\\\ `_"))
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
			NewId("a"),
			NewNull(),
			NewId("b"),
			NewNull(),
			NewId("c"),
			NewNull(),
			NewId("d"),
			NewNull(),
			NewNull(),
			NewId("e"),
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
				"\tACTUAL:   " + token.String() + "\n" +
				"\tEXPECTED: " + expectedTokens[index].String())
		}
	}

	if len(vals) != len(expectedVals) {
		t.Error(str + " - Actual and expected value list lengths don't match")
	}
	for index, val := range vals {
		if val.ToZinc() != expectedVals[index].ToZinc() {
			t.Error(str + " - Val doesn't match expected\n" +
				"\tACTUAL:   " + val.ToZinc() + "\n" +
				"\tEXPECTED: " + expectedVals[index].ToZinc())
		}
	}
}

func testTokenizerRead(t *testing.T, str string) ([]Token, []Val) {
	var tokenizer Tokenizer
	tokenizer.InitString(str)

	var tokens []Token
	var vals []Val

	for {
		nextToken := tokenizer.Next()
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
