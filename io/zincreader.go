package io

import (
	"math"
	"strings"
	"unicode"

	"gitlab.com/NeedleInAJayStack/haystack"
)

type ZincReader struct {
	tokenizer Tokenizer

	cur    Token
	curVal haystack.Val
	// curLine int

	peek    Token
	peekVal haystack.Val
	// peekLine int
}

func (reader *ZincReader) InitString(str string) {
	reader.Init(strings.NewReader(str))
}

func (reader *ZincReader) Init(in *strings.Reader) {
	reader.tokenizer = Tokenizer{}
	reader.tokenizer.Init(in)

	reader.consume()
	reader.consume()
}

func (reader *ZincReader) ReadVal() haystack.Val {
	var val haystack.Val

	if reader.cur == ID {
		val = reader.parseGrid()
	} else {
		val = reader.parseVal()
	}

	if reader.cur != EOF {
		panic("Expecting EOF, not " + reader.cur.String())
	}
	return val
}

func (reader *ZincReader) parseVal() haystack.Val {
	if reader.cur == ID {
		id := reader.curVal.(haystack.Id).String()
		reader.consumeToken(ID)

		// check for coord or xstr
		if reader.cur == LPAREN {
			if reader.peek == NUMBER {
				return reader.parseCoord(id)
			} else {
				return reader.parseXStr(id)
			}
		}

		// check for keyword
		if id == "T" {
			return haystack.NewBool(true)
		} else if id == "F" {
			return haystack.NewBool(false)
		} else if id == "N" {
			return haystack.NewNull()
		} else if id == "M" {
			return haystack.NewMarker()
		} else if id == "NA" {
			return haystack.NewNA()
		} else if id == "R" {
			return haystack.NewRemove()
		} else if id == "NaN" {
			return haystack.NewNumber(math.NaN(), "")
		} else if id == "INF" {
			return haystack.NewNumber(math.Inf(1), "")
		} else {
			panic("Unexpected identifier: " + id)
		}
	}

	// literals
	if reader.cur.IsLiteral() {
		return reader.parseLiteral()
	}

	// -INF
	if reader.cur == MINUS && reader.peekVal.ToZinc() == "INF" {
		reader.consumeToken(MINUS)
		reader.consumeToken(ID)
		return haystack.NewNumber(math.Inf(-1), "")
	}

	// nested collections
	if reader.cur == LBRACKET {
		return reader.parseList()
	} else if reader.cur == LBRACE {
		return reader.parseDict()
	} else if reader.cur == LT2 {
		return reader.parseGrid()
	}

	panic("Unexpected token: " + reader.cur.String())
}

func (reader *ZincReader) parseCoord(id string) haystack.Coord {
	if id != "C" {
		panic("Expecting 'C' for coord, not " + id)
	}

	var lat haystack.Number
	var lng haystack.Number
	reader.consumeToken(LPAREN)
	lat = reader.consumeNumber()
	reader.consumeToken(COMMA)
	lng = reader.consumeNumber()
	reader.consumeToken(RPAREN)

	return haystack.NewCoord(lat.Float(), lng.Float())
}

func (reader *ZincReader) parseXStr(id string) haystack.Val {
	if !unicode.IsUpper([]rune(id)[0]) {
		panic("Invalid XStr type: " + id)
	}
	if id == "Bin" { // I think Bins are obselete
		reader.consumeToken(LPAREN)
		mime := reader.consumeStr()
		reader.consumeToken(RPAREN)

		return haystack.NewBin(mime.String())
	} else {
		reader.consumeToken(LPAREN)
		val := reader.consumeStr()
		reader.consumeToken(RPAREN)

		return haystack.NewXStr(id, val.String())
	}
}

func (reader *ZincReader) parseLiteral() haystack.Val {
	val := reader.curVal
	// Combine ref and dis
	if reader.cur == REF && reader.peek == STR {
		ref := reader.curVal.(haystack.Ref)
		dis := reader.peekVal.(haystack.Str)

		val = haystack.NewRef(ref.Id(), dis.String())
		reader.consumeToken(REF)
	}
	reader.consume()
	return val
}

func (reader *ZincReader) parseList() haystack.List {
	var vals []haystack.Val

	reader.consumeToken(LBRACKET)
	for reader.cur != RBRACKET && reader.cur != EOF {
		var val haystack.Val
		val = reader.parseVal()
		vals = append(vals, val)
		if reader.cur != COMMA {
			break
		}
		reader.consumeToken(COMMA)
	}

	reader.consumeToken(RBRACKET)

	return haystack.NewList(vals)
}

func (reader *ZincReader) parseDict() haystack.Dict {
	items := make(map[string]haystack.Val)

	braces := reader.cur == LBRACE
	if braces {
		reader.consumeToken(LBRACE)
	}
	for reader.cur == ID {
		var id string
		var val haystack.Val

		id = reader.consumeTagName()

		val = haystack.NewMarker() // Default to marker val if there is no value
		if reader.cur == COLON {
			reader.consumeToken(COLON)
			val = reader.parseVal()
		}
		items[id] = val
	}
	if braces {
		reader.consumeToken(RBRACE)
	}

	return haystack.NewDict(items)
}

func (reader *ZincReader) parseGrid() haystack.Grid {
	var gb haystack.GridBuilder

	nested := reader.cur == LT2
	if nested {
		reader.consumeToken(LT2)
		if reader.cur == NL {
			reader.consumeToken(NL)
		}
	}

	// ver:"3.0"
	if reader.cur != ID {
		panic("Expecting grid 'ver' identifier, not " + reader.curVal.ToZinc())
	}
	reader.consume()
	reader.consumeToken(COLON)
	ver := reader.consumeStr()
	checkVersion(ver.String())

	// grid meta
	if reader.cur == ID {
		gb.SetMetaDict(reader.parseDict())
	}
	reader.consumeToken(NL)

	// column definitions
	numCols := 0
	for reader.cur == ID {
		numCols = numCols + 1
		name := reader.consumeTagName()

		var colMeta haystack.Dict
		if reader.cur == ID {
			colMeta = reader.parseDict()
		}
		gb.AddColDict(name, colMeta)

		if reader.cur != COMMA {
			break
		}
		reader.consumeToken(COMMA)
	}
	if numCols == 0 {
		panic("No columns defined")
	}
	reader.consumeToken(NL)

	// grid rows
	for {
		if reader.cur == NL {
			break
		} else if reader.cur == EOF {
			break
		} else if nested && reader.cur == GT2 {
			break
		}

		// read cells
		var vals []haystack.Val
		for i := 0; i < numCols; i = i + 1 {
			if reader.cur == COMMA || reader.cur == NL || reader.cur == EOF {
				vals = append(vals, haystack.NewNull())
			} else {
				vals = append(vals, reader.parseVal())
			}
			if i+1 < numCols {
				reader.consumeToken(COMMA)
			}
		}
		gb.AddRow(vals)

		// newline or end
		if nested && reader.cur == GT2 {
			break
		} else if reader.cur == EOF {
			break
		}
		reader.consumeToken(NL)
	}

	if reader.cur == NL {
		reader.consumeToken(NL)
	}
	if nested {
		reader.consumeToken(GT2)
	}

	return gb.ToGrid()
}

func (reader *ZincReader) consumeTagName() string {
	id := reader.curVal.(haystack.Id)
	val := id.String()
	if val == "" || unicode.IsUpper([]rune(val)[0]) {
		panic("Invalid dict tag name: " + val)
	}
	reader.consumeToken(ID)
	return val
}

func checkVersion(str string) {
	if str != "2.0" && str != "3.0" {
		panic("Unsupported version: " + str)
	}
}

func (reader *ZincReader) consumeNumber() haystack.Number {
	number := reader.curVal.(haystack.Number)
	reader.consumeToken(NUMBER)
	return number
}

func (reader *ZincReader) consumeStr() haystack.Str {
	str := reader.curVal.(haystack.Str)
	reader.consumeToken(STR)
	return str
}

func (reader *ZincReader) consumeToken(expected Token) {
	if reader.cur != expected {
		panic("Expected " + expected.String() + " not " + reader.cur.String())
	}
	reader.consume()
}

func (reader *ZincReader) consume() {
	newToken := reader.tokenizer.Next()

	reader.cur = reader.peek
	reader.curVal = reader.peekVal
	// reader.curLine = reader.peekLine

	reader.peek = newToken
	reader.peekVal = reader.tokenizer.Val()
	// reader.peekLine = reader.tokenizer.line
}
