package io

import (
	"math"
	"testing"

	"github.com/NeedleInAJayStack/haystack"
)

func TestZincReader_empty(t *testing.T) {
	input := "ver:\"3.0\" tag:N\n" +
		"a nullmetatag:N, b markermetatag\n" +
		""

	gb := haystack.NewGridBuilder()
	gb.SetMeta(
		map[string]haystack.Val{
			"tag": haystack.NewNull(),
		},
	)
	gb.AddCol(
		"a",
		map[string]haystack.Val{
			"nullmetatag": haystack.NewNull(),
		},
	)
	gb.AddCol(
		"b",
		map[string]haystack.Val{
			"markermetatag": haystack.NewMarker(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_singleColEmpty(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"fooBar33\n" +
		"\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol(
		"fooBar33",
		map[string]haystack.Val{},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_singleCol(t *testing.T) {
	input := "ver:\"3.0\" tag foo:\"bar\"\n" +
		"xyz\n" +
		"\"val\"\n" +
		"\n"

	gb := haystack.NewGridBuilder()
	gb.SetMeta(
		map[string]haystack.Val{
			"tag": haystack.NewMarker(),
			"foo": haystack.NewStr("bar"),
		},
	)
	gb.AddCol(
		"xyz",
		map[string]haystack.Val{},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewStr("val"),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_singleColNull(t *testing.T) {
	input := "ver:\"3.0\"\n" +
		"val\n" +
		"N\n" +
		"\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol(
		"val",
		map[string]haystack.Val{},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNull(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_doubleCol(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"a,b\n" +
		"1,2\n" +
		"3,4\n" +
		"\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol(
		"a",
		map[string]haystack.Val{},
	)
	gb.AddCol(
		"b",
		map[string]haystack.Val{},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(1.0, ""),
			haystack.NewNumber(2.0, ""),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(3.0, ""),
			haystack.NewNumber(4.0, ""),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_large(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"a,    b,      c,      d\n" +
		"T,    F,      N,   -99\n" +
		"2.3,  -5e-10, 2.4e20, 123e-10\n" +
		"\"\",   \"a\",   \"\\\" \\\\ \\t \\n \\r\", \"\\uabcd\"\n" +
		"`path`, @12cbb082-0c02ae73, 4s, -2.5min\n" +
		"M,R,N,N\n" + // "M,R,Bin(image/png),Bin(image/png)\n" + // Don't support Bins yet
		"2009-12-31, 23:59:01, 01:02:03.123, 2009-02-03T04:05:06Z\n" +
		"INF, -INF, \"\", NaN\n" +
		"C(12,-34),C(0.123,-0.789),C(84.5,-77.45),C(-90,180)\n" +
		"\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol("a", map[string]haystack.Val{})
	gb.AddCol("b", map[string]haystack.Val{})
	gb.AddCol("c", map[string]haystack.Val{})
	gb.AddCol("d", map[string]haystack.Val{})
	gb.AddRow( // T,    F,      N,   -99
		[]haystack.Val{
			haystack.NewBool(true),
			haystack.NewBool(false),
			haystack.NewNull(),
			haystack.NewNumber(-99.0, ""),
		},
	)
	gb.AddRow( // 2.3,  -5e-10, 2.4e20, 123e-10
		[]haystack.Val{
			haystack.NewNumber(2.3, ""),
			haystack.NewNumber(-5e-10, ""),
			haystack.NewNumber(2.4e20, ""),
			haystack.NewNumber(123e-10, ""),
		},
	)
	gb.AddRow( // "",   "a",   "\" \\ \t \n \r", "\uabcd"
		[]haystack.Val{
			haystack.NewStr(""),
			haystack.NewStr("a"),
			haystack.NewStr("\" \\ \t \n \r"),
			haystack.NewStr("\uabcd"),
		},
	)
	gb.AddRow( // `path`, @12cbb082-0c02ae73, 4s, -2.5min
		[]haystack.Val{
			haystack.NewUri("path"),
			haystack.NewRef("12cbb082-0c02ae73", ""),
			haystack.NewNumber(4.0, "s"),
			haystack.NewNumber(-2.5, "min"),
		},
	)
	gb.AddRow( // M,R,N,N
		[]haystack.Val{
			haystack.NewMarker(),
			haystack.NewRemove(),
			haystack.NewNull(),
			haystack.NewNull(),
		},
	)
	date, _ := haystack.NewDateFromIso("2009-12-31")
	time1, _ := haystack.NewTimeFromIso("23:59:01")
	time2, _ := haystack.NewTimeFromIso("01:02:03.123")
	datetime, _ := haystack.NewDateTimeFromString("2009-02-03T04:05:06Z")
	gb.AddRow( // 2009-12-31, 23:59:01, 01:02:03.123, 2009-02-03T04:05:06Z
		[]haystack.Val{
			date,
			time1,
			time2,
			datetime,
		},
	)
	gb.AddRow( // INF, -INF, \"\", NaN
		[]haystack.Val{
			haystack.NewNumber(math.Inf(1), ""),
			haystack.NewNumber(math.Inf(-1), ""),
			haystack.NewStr(""),
			haystack.NewNumber(math.NaN(), ""),
		},
	)
	gb.AddRow( // C(12,-34),C(0.123,-0.789),C(84.5,-77.45),C(-90,180)
		[]haystack.Val{
			haystack.NewCoord(12.0, -34.0),
			haystack.NewCoord(0.123, -0.789),
			haystack.NewCoord(84.5, -77.45),
			haystack.NewCoord(-90, 180),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_escapes(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"foo\n" +
		"`foo$20bar`\n" +
		"`foo\\`bar`\n" +
		"`file \\#2`\n" +
		"\"$15\"\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol("foo", map[string]haystack.Val{})
	gb.AddRow(
		[]haystack.Val{
			haystack.NewUri("foo$20bar"),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewUri("foo`bar"),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewUri("file \\#2"),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewStr("$15"),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_numbers(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"a, b\n" +
		"-3.1kg,4kg\n" +
		"5%,3.2%\n" +
		"5kWh/ft\u00b2,-15kWh/m\u00b2\n" +
		"123e+12kJ/kg_dry,74\u0394\u00b0F\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol("a", map[string]haystack.Val{})
	gb.AddCol("b", map[string]haystack.Val{})
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(-3.1, "kg"),
			haystack.NewNumber(4.0, "kg"),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(5.0, "%"),
			haystack.NewNumber(3.2, "%"),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(5.0, "kWh/ft\u00b2"),
			haystack.NewNumber(-15.0, "kWh/m\u00b2"),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(123e+12, "kJ/kg_dry"),
			haystack.NewNumber(74.0, "\u0394\u00b0F"),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_nulls(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"a, b, c\n" +
		", 1, 2\n" +
		"3, , 5\n" +
		"6, 7_000,\n" +
		",,10\n" +
		",,\n" +
		"14,,\n" +
		"\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol("a", map[string]haystack.Val{})
	gb.AddCol("b", map[string]haystack.Val{})
	gb.AddCol("c", map[string]haystack.Val{})
	gb.AddRow( // , 1, 2
		[]haystack.Val{
			haystack.NewNull(),
			haystack.NewNumber(1.0, ""),
			haystack.NewNumber(2.0, ""),
		},
	)
	gb.AddRow( // 3, , 5
		[]haystack.Val{
			haystack.NewNumber(3.0, ""),
			haystack.NewNull(),
			haystack.NewNumber(5.0, ""),
		},
	)
	gb.AddRow( // 6, 7_000,
		[]haystack.Val{
			haystack.NewNumber(6.0, ""),
			haystack.NewNumber(7000.0, ""),
			haystack.NewNull(),
		},
	)
	gb.AddRow( // ,,10
		[]haystack.Val{
			haystack.NewNull(),
			haystack.NewNull(),
			haystack.NewNumber(10.0, ""),
		},
	)
	gb.AddRow( // ,,
		[]haystack.Val{
			haystack.NewNull(),
			haystack.NewNull(),
			haystack.NewNull(),
		},
	)
	gb.AddRow( // 14,,
		[]haystack.Val{
			haystack.NewNumber(14.0, ""),
			haystack.NewNull(),
			haystack.NewNull(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_datetimes(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"a,b\n" +
		"2010-03-01T23:55:00.013-05:00 GMT+5,2010-03-01T23:55:00.013+10:00 GMT-10\n"

	gb := haystack.NewGridBuilder()
	gb.AddCol("a", map[string]haystack.Val{})
	gb.AddCol("b", map[string]haystack.Val{})
	datetime1, _ := haystack.NewDateTimeFromString("2010-03-01T23:55:00.013-05:00 GMT+5")
	datetime2, _ := haystack.NewDateTimeFromString("2010-03-01T23:55:00.013+10:00 GMT-10")
	gb.AddRow(
		[]haystack.Val{
			datetime1,
			datetime2,
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}

func TestZincReader_meta(t *testing.T) {
	input := "ver:\"2.0\" a: 2009-02-03T04:05:06Z foo b: 2010-02-03T04:05:06Z UTC bar c: 2009-12-03T04:05:06Z London baz\n" +
		"a\n" +
		"3.814697265625E-6\n" +
		"2010-12-18T14:11:30.924Z\n" +
		"2010-12-18T14:11:30.925Z UTC\n" +
		"2010-12-18T14:11:30.925Z London\n" +
		"45$\n" +
		"33\u00a3\n" +
		"@12cbb08e-0c02ae73\n" +
		"7.15625E-4kWh/ft\u00b2\n" +
		"R\n" +
		"NA\n"

	gb := haystack.NewGridBuilder()
	metaDt1, _ := haystack.NewDateTimeFromString("2009-02-03T04:05:06Z")
	metaDt2, _ := haystack.NewDateTimeFromString("2010-02-03T04:05:06Z UTC")
	metaDt3, _ := haystack.NewDateTimeFromString("2009-12-03T04:05:06Z London")
	gb.SetMeta(
		map[string]haystack.Val{
			"a":   metaDt1,
			"foo": haystack.NewMarker(),
			"b":   metaDt2,
			"bar": haystack.NewMarker(),
			"c":   metaDt3,
			"baz": haystack.NewMarker(),
		},
	)
	gb.AddCol("a", map[string]haystack.Val{})
	gb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(3.814697265625e-6, ""),
		},
	)
	datetime1, _ := haystack.NewDateTimeFromString("2010-12-18T14:11:30.924Z")
	gb.AddRow(
		[]haystack.Val{
			datetime1,
		},
	)
	datetime2, _ := haystack.NewDateTimeFromString("2010-12-18T14:11:30.925Z UTC")
	gb.AddRow(
		[]haystack.Val{
			datetime2,
		},
	)
	datetime3, _ := haystack.NewDateTimeFromString("2010-12-18T14:11:30.925Z London")
	gb.AddRow(
		[]haystack.Val{
			datetime3,
		},
	)
	gb.AddRow( // 45$
		[]haystack.Val{
			haystack.NewNumber(45, "$"),
		},
	)
	gb.AddRow( // 33\u00a3
		[]haystack.Val{
			haystack.NewNumber(33, "\u00a3"),
		},
	)
	gb.AddRow( // @12cbb08e-0c02ae73
		[]haystack.Val{
			haystack.NewRef("12cbb08e-0c02ae73", ""),
		},
	)
	gb.AddRow( // 7.15625E-4kWh/ft\u00b2
		[]haystack.Val{
			haystack.NewNumber(7.15625e-4, "kWh/ft\u00b2"),
		},
	)
	gb.AddRow( // R
		[]haystack.Val{
			haystack.NewRemove(),
		},
	)
	gb.AddRow( // NA
		[]haystack.Val{
			haystack.NewNA(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}

func TestZincReader_nested(t *testing.T) {
	input := "ver:\"3.0\"\n" +
		"type, val\n" +
		"\"list\", [1, 2, 3]\n" +
		"\"dict\", {dis:\"Dict!\" foo}\n" +
		"\"grid\", <<\n" +
		"  ver:\"3.0\"\n" +
		"  a, b\n" +
		"  1, <<\n" +
		"    ver:\"3.0\"\n" +
		"    c, d\n" +
		"    5, 6\n" +
		"    >>\n" +
		"  3, 4\n" +
		"  >>\n" +
		"\"scalar\", \"simple string\""

	gb := haystack.NewGridBuilder()
	gb.AddCol("type", map[string]haystack.Val{})
	gb.AddCol("val", map[string]haystack.Val{})
	gb.AddRow(
		[]haystack.Val{
			haystack.NewStr("list"),
			haystack.NewList(
				[]haystack.Val{
					haystack.NewNumber(1, ""),
					haystack.NewNumber(2, ""),
					haystack.NewNumber(3, ""),
				},
			),
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewStr("dict"),
			haystack.NewDict(
				map[string]haystack.Val{
					"dis": haystack.NewStr("Dict!"),
					"foo": haystack.NewMarker(),
				},
			),
		},
	)
	var dblNestedGb haystack.GridBuilder
	dblNestedGb.AddCol("c", map[string]haystack.Val{})
	dblNestedGb.AddCol("d", map[string]haystack.Val{})
	dblNestedGb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(5, ""),
			haystack.NewNumber(6, ""),
		},
	)
	dblNestedGrid := dblNestedGb.ToGrid()
	var nestedGb haystack.GridBuilder
	nestedGb.AddCol("a", map[string]haystack.Val{})
	nestedGb.AddCol("b", map[string]haystack.Val{})
	nestedGb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(1, ""),
			dblNestedGrid,
		},
	)
	nestedGb.AddRow(
		[]haystack.Val{
			haystack.NewNumber(3, ""),
			haystack.NewNumber(4, ""),
		},
	)
	nestedGrid := nestedGb.ToGrid()
	gb.AddRow(
		[]haystack.Val{
			haystack.NewStr("grid"),
			nestedGrid,
		},
	)
	gb.AddRow(
		[]haystack.Val{
			haystack.NewStr("scalar"),
			haystack.NewStr("simple string"),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}

// UTILITIES

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'ToZinc' method
func testZincReaderGrid(t *testing.T, str string, expected haystack.Grid) {
	var reader ZincReader
	reader.InitString(str)

	val := reader.ReadVal()
	grid := val.(haystack.Grid)
	testGridEq(t, grid, expected)

	// write grid, read grid, and verify it equals the original
	writeStr := grid.ToZinc()
	var writtenReader ZincReader
	writtenReader.InitString(writeStr)
	writeReadVal := writtenReader.ReadVal()
	writeReadGrid := writeReadVal.(haystack.Grid)
	testGridEq(t, writeReadGrid, expected)
}

// Test whether the grids match based on a ToZinc call
func testGridEq(t *testing.T, actual haystack.Grid, expected haystack.Grid) {
	actualZinc := actual.ToZinc()
	expectedZinc := expected.ToZinc()

	if actualZinc != expectedZinc {
		t.Error("Grids don't match\nACTUAL:\n" + actualZinc + "\nEXPECTED:\n" + expectedZinc)
	}
}
