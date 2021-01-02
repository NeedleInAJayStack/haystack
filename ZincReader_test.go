package haystack

import (
	"math"
	"testing"
)

func TestZincReader_empty(t *testing.T) {
	input := "ver:\"3.0\" tag:N\n" +
		"a nullmetatag:N, b markermetatag\n" +
		""

	gb := NewGridBuilder()
	gb.SetMeta(
		map[string]Val{
			"tag": NewNull(),
		},
	)
	gb.AddCol(
		"a",
		map[string]Val{
			"nullmetatag": NewNull(),
		},
	)
	gb.AddCol(
		"b",
		map[string]Val{
			"markermetatag": NewMarker(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_singleColEmpty(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"fooBar33\n" +
		"\n"

	gb := NewGridBuilder()
	gb.AddCol(
		"fooBar33",
		map[string]Val{},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_singleCol(t *testing.T) {
	input := "ver:\"3.0\" tag foo:\"bar\"\n" +
		"xyz\n" +
		"\"val\"\n" +
		"\n"

	gb := NewGridBuilder()
	gb.SetMeta(
		map[string]Val{
			"tag": NewMarker(),
			"foo": NewStr("bar"),
		},
	)
	gb.AddCol(
		"xyz",
		map[string]Val{},
	)
	gb.AddRow(
		[]Val{
			NewStr("val"),
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

	gb := NewGridBuilder()
	gb.AddCol(
		"val",
		map[string]Val{},
	)
	gb.AddRow(
		[]Val{
			NewNull(),
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

	gb := NewGridBuilder()
	gb.AddCol(
		"a",
		map[string]Val{},
	)
	gb.AddCol(
		"b",
		map[string]Val{},
	)
	gb.AddRow(
		[]Val{
			NewNumber(1.0, ""),
			NewNumber(2.0, ""),
		},
	)
	gb.AddRow(
		[]Val{
			NewNumber(3.0, ""),
			NewNumber(4.0, ""),
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

	gb := NewGridBuilder()
	gb.AddCol("a", map[string]Val{})
	gb.AddCol("b", map[string]Val{})
	gb.AddCol("c", map[string]Val{})
	gb.AddCol("d", map[string]Val{})
	gb.AddRow( // T,    F,      N,   -99
		[]Val{
			NewBool(true),
			NewBool(false),
			NewNull(),
			NewNumber(-99.0, ""),
		},
	)
	gb.AddRow( // 2.3,  -5e-10, 2.4e20, 123e-10
		[]Val{
			NewNumber(2.3, ""),
			NewNumber(-5e-10, ""),
			NewNumber(2.4e20, ""),
			NewNumber(123e-10, ""),
		},
	)
	gb.AddRow( // "",   "a",   "\" \\ \t \n \r", "\uabcd"
		[]Val{
			NewStr(""),
			NewStr("a"),
			NewStr("\" \\ \t \n \r"),
			NewStr("\uabcd"),
		},
	)
	gb.AddRow( // `path`, @12cbb082-0c02ae73, 4s, -2.5min
		[]Val{
			NewUri("path"),
			NewRef("12cbb082-0c02ae73", ""),
			NewNumber(4.0, "s"),
			NewNumber(-2.5, "min"),
		},
	)
	gb.AddRow( // M,R,N,N
		[]Val{
			NewMarker(),
			NewRemove(),
			NewNull(),
			NewNull(),
		},
	)
	date, _ := NewDateFromIso("2009-12-31")
	time1, _ := NewTimeFromIso("23:59:01")
	time2, _ := NewTimeFromIso("01:02:03.123")
	datetime, _ := NewDateTimeFromString("2009-02-03T04:05:06Z")
	gb.AddRow( // 2009-12-31, 23:59:01, 01:02:03.123, 2009-02-03T04:05:06Z
		[]Val{
			date,
			time1,
			time2,
			datetime,
		},
	)
	gb.AddRow( // INF, -INF, \"\", NaN
		[]Val{
			NewNumber(math.Inf(1), ""),
			NewNumber(math.Inf(-1), ""),
			NewStr(""),
			NewNumber(math.NaN(), ""),
		},
	)
	gb.AddRow( // C(12,-34),C(0.123,-0.789),C(84.5,-77.45),C(-90,180)
		[]Val{
			NewCoord(12.0, -34.0),
			NewCoord(0.123, -0.789),
			NewCoord(84.5, -77.45),
			NewCoord(-90, 180),
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

	gb := NewGridBuilder()
	gb.AddCol("foo", map[string]Val{})
	gb.AddRow(
		[]Val{
			NewUri("foo$20bar"),
		},
	)
	gb.AddRow(
		[]Val{
			NewUri("foo`bar"),
		},
	)
	gb.AddRow(
		[]Val{
			NewUri("file \\#2"),
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("$15"),
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

	gb := NewGridBuilder()
	gb.AddCol("a", map[string]Val{})
	gb.AddCol("b", map[string]Val{})
	gb.AddRow(
		[]Val{
			NewNumber(-3.1, "kg"),
			NewNumber(4.0, "kg"),
		},
	)
	gb.AddRow(
		[]Val{
			NewNumber(5.0, "%"),
			NewNumber(3.2, "%"),
		},
	)
	gb.AddRow(
		[]Val{
			NewNumber(5.0, "kWh/ft\u00b2"),
			NewNumber(-15.0, "kWh/m\u00b2"),
		},
	)
	gb.AddRow(
		[]Val{
			NewNumber(123e+12, "kJ/kg_dry"),
			NewNumber(74.0, "\u0394\u00b0F"),
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

	gb := NewGridBuilder()
	gb.AddCol("a", map[string]Val{})
	gb.AddCol("b", map[string]Val{})
	gb.AddCol("c", map[string]Val{})
	gb.AddRow( // , 1, 2
		[]Val{
			NewNull(),
			NewNumber(1.0, ""),
			NewNumber(2.0, ""),
		},
	)
	gb.AddRow( // 3, , 5
		[]Val{
			NewNumber(3.0, ""),
			NewNull(),
			NewNumber(5.0, ""),
		},
	)
	gb.AddRow( // 6, 7_000,
		[]Val{
			NewNumber(6.0, ""),
			NewNumber(7000.0, ""),
			NewNull(),
		},
	)
	gb.AddRow( // ,,10
		[]Val{
			NewNull(),
			NewNull(),
			NewNumber(10.0, ""),
		},
	)
	gb.AddRow( // ,,
		[]Val{
			NewNull(),
			NewNull(),
			NewNull(),
		},
	)
	gb.AddRow( // 14,,
		[]Val{
			NewNumber(14.0, ""),
			NewNull(),
			NewNull(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}
func TestZincReader_datetimes(t *testing.T) {
	input := "ver:\"2.0\"\n" +
		"a,b\n" +
		"2010-03-01T23:55:00.013-05:00 GMT+5,2010-03-01T23:55:00.013+10:00 GMT-10\n"

	gb := NewGridBuilder()
	gb.AddCol("a", map[string]Val{})
	gb.AddCol("b", map[string]Val{})
	datetime1, _ := NewDateTimeFromString("2010-03-01T23:55:00.013-05:00 GMT+5")
	datetime2, _ := NewDateTimeFromString("2010-03-01T23:55:00.013+10:00 GMT-10")
	gb.AddRow(
		[]Val{
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

	gb := NewGridBuilder()
	metaDt1, _ := NewDateTimeFromString("2009-02-03T04:05:06Z")
	metaDt2, _ := NewDateTimeFromString("2010-02-03T04:05:06Z UTC")
	metaDt3, _ := NewDateTimeFromString("2009-12-03T04:05:06Z London")
	gb.SetMeta(
		map[string]Val{
			"a":   metaDt1,
			"foo": NewMarker(),
			"b":   metaDt2,
			"bar": NewMarker(),
			"c":   metaDt3,
			"baz": NewMarker(),
		},
	)
	gb.AddCol("a", map[string]Val{})
	gb.AddRow(
		[]Val{
			NewNumber(3.814697265625e-6, ""),
		},
	)
	datetime1, _ := NewDateTimeFromString("2010-12-18T14:11:30.924Z")
	gb.AddRow(
		[]Val{
			datetime1,
		},
	)
	datetime2, _ := NewDateTimeFromString("2010-12-18T14:11:30.925Z UTC")
	gb.AddRow(
		[]Val{
			datetime2,
		},
	)
	datetime3, _ := NewDateTimeFromString("2010-12-18T14:11:30.925Z London")
	gb.AddRow(
		[]Val{
			datetime3,
		},
	)
	gb.AddRow( // 45$
		[]Val{
			NewNumber(45, "$"),
		},
	)
	gb.AddRow( // 33\u00a3
		[]Val{
			NewNumber(33, "\u00a3"),
		},
	)
	gb.AddRow( // @12cbb08e-0c02ae73
		[]Val{
			NewRef("12cbb08e-0c02ae73", ""),
		},
	)
	gb.AddRow( // 7.15625E-4kWh/ft\u00b2
		[]Val{
			NewNumber(7.15625e-4, "kWh/ft\u00b2"),
		},
	)
	gb.AddRow( // R
		[]Val{
			NewRemove(),
		},
	)
	gb.AddRow( // NA
		[]Val{
			NewNA(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}

// TODO improve support for writing nested grids
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

	gb := NewGridBuilder()
	gb.AddCol("type", map[string]Val{})
	gb.AddCol("val", map[string]Val{})
	gb.AddRow(
		[]Val{
			NewStr("list"),
			NewList(
				[]Val{
					NewNumber(1, ""),
					NewNumber(2, ""),
					NewNumber(3, ""),
				},
			),
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("dict"),
			NewDict(
				map[string]Val{
					"dis": NewStr("Dict!"),
					"foo": NewMarker(),
				},
			),
		},
	)
	var dblNestedGb GridBuilder
	dblNestedGb.AddCol("c", map[string]Val{})
	dblNestedGb.AddCol("d", map[string]Val{})
	dblNestedGb.AddRow(
		[]Val{
			NewNumber(5, ""),
			NewNumber(6, ""),
		},
	)
	dblNestedGrid := dblNestedGb.ToGrid()
	var nestedGb GridBuilder
	nestedGb.AddCol("a", map[string]Val{})
	nestedGb.AddCol("b", map[string]Val{})
	nestedGb.AddRow(
		[]Val{
			NewNumber(1, ""),
			dblNestedGrid,
		},
	)
	nestedGb.AddRow(
		[]Val{
			NewNumber(3, ""),
			NewNumber(4, ""),
		},
	)
	nestedGrid := nestedGb.ToGrid()
	gb.AddRow(
		[]Val{
			NewStr("grid"),
			nestedGrid,
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("scalar"),
			NewStr("simple string"),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
}

// UTILITIES

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'ToZinc' method
func testZincReaderGrid(t *testing.T, str string, expected *Grid) {
	var reader ZincReader
	reader.InitString(str)

	val := reader.ReadVal()
	grid := val.(*Grid)
	testGridEq(t, grid, expected)

	// write grid, read grid, and verify it equals the original
	writeStr := grid.ToZinc()
	var writtenReader ZincReader
	writtenReader.InitString(writeStr)
	writeReadVal := writtenReader.ReadVal()
	writeReadGrid := writeReadVal.(*Grid)
	testGridEq(t, writeReadGrid, expected)
}

// Test whether the grids match based on a ToZinc call
func testGridEq(t *testing.T, actual *Grid, expected *Grid) {
	actualZinc := actual.ToZinc()
	expectedZinc := expected.ToZinc()

	if actualZinc != expectedZinc {
		t.Error("Grids don't match\nACTUAL:\n" + actualZinc + "\nEXPECTED:\n" + expectedZinc)
	}
}
