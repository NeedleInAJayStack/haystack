package haystack

import (
	"math"
	"strings"
	"unicode"
)

type ZincReader struct {
	tokenizer Tokenizer

	cur    Token
	curVal Val
	// curLine int

	peek    Token
	peekVal Val
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

func (reader *ZincReader) ReadVal() Val {
	var val Val

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

func (reader *ZincReader) parseVal() Val {
	if reader.cur == ID {
		id := reader.curVal.(Id).val
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
			return NewBool(true)
		} else if id == "F" {
			return NewBool(false)
		} else if id == "N" {
			return NewNull()
		} else if id == "M" {
			return NewMarker()
		} else if id == "NA" {
			return NewNA()
		} else if id == "R" {
			return NewRemove()
		} else if id == "NaN" {
			return NewNumber(math.NaN(), "")
		} else if id == "INF" {
			return NewNumber(math.Inf(1), "")
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
		return NewNumber(math.Inf(-1), "")
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

func (reader *ZincReader) parseCoord(id string) Coord {
	if id != "C" {
		panic("Expecting 'C' for coord, not " + id)
	}

	var lat Number
	var lng Number
	reader.consumeToken(LPAREN)
	lat = reader.consumeNumber()
	reader.consumeToken(COMMA)
	lng = reader.consumeNumber()
	reader.consumeToken(RPAREN)

	return NewCoord(lat.val, lng.val)
}

func (reader *ZincReader) parseXStr(id string) XStr {
	if !unicode.IsUpper([]rune(id)[0]) {
		panic("Invalid XStr type: " + id)
	}

	var val Str
	reader.consumeToken(LPAREN)
	val = reader.consumeStr()
	reader.consumeToken(RPAREN)

	return NewXStr(id, val.val)
}

func (reader *ZincReader) parseLiteral() Val {
	val := reader.curVal
	// Combine ref and dis
	if reader.cur == REF && reader.peek == STR {
		ref := reader.curVal.(Ref)
		dis := reader.peekVal.(Str)

		val = NewRef(ref.val, dis.val)
		reader.consumeToken(REF)
	}
	reader.consume()
	return val
}

func (reader *ZincReader) parseList() List {
	var vals []Val

	reader.consumeToken(LBRACKET)
	for reader.cur != RBRACKET && reader.cur != EOF {
		var val Val
		val = reader.parseVal()
		vals = append(vals, val)
		if reader.cur == COMMA {
			break
		}
		reader.consumeToken(COMMA)
	}

	reader.consumeToken(RBRACKET)

	return NewList(vals)
}

func (reader *ZincReader) parseDict() Dict {
	items := make(map[string]Val)

	braces := reader.cur == LBRACE
	if braces {
		reader.consumeToken(LBRACE)
	}
	for reader.cur == ID {
		var id string
		var val Val

		id = reader.consumeTagName()

		val = Marker{} // Default to marker val if there is no value
		if reader.cur == COLON {
			reader.consumeToken(COLON)
			val = reader.parseVal()
		}
		items[id] = val
	}
	if braces {
		reader.consumeToken(RBRACE)
	}

	return Dict{items: items}
}

func (reader *ZincReader) parseGrid() Grid {
	var meta Dict
	var cols []Col
	var rows []Row

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
	reader.consumeStr() // Always expect version 3
	// TODO Check for version

	// grid meta
	if reader.cur == ID {
		meta = reader.parseDict()
	}
	reader.consumeToken(NL)

	// column definitions
	numCols := 0
	for reader.cur == ID {
		numCols = numCols + 1
		name := reader.consumeTagName()

		var colMeta Dict
		if reader.cur == ID {
			colMeta = reader.parseDict()
		}
		col := Col{
			index: numCols,
			name:  name,
			meta:  colMeta,
		}
		cols = append(cols, col)

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
		vals := make(map[string]Val)
		for i := 0; i < numCols; i = i + 1 {
			col := cols[i]
			if reader.cur == COMMA || reader.cur == NL || reader.cur == EOF {
				vals[col.Name()] = Null{}
			} else {
				vals[col.Name()] = reader.parseVal()
			}
			if i+1 < numCols {
				reader.consumeToken(COMMA)
			}
		}
		rows = append(rows, Row{items: vals})

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

	return Grid{
		meta: meta,
		cols: cols,
		rows: rows,
	}
}

func (reader *ZincReader) consumeTagName() string {
	id := reader.curVal.(Id)
	val := id.val
	if val == "" || unicode.IsUpper([]rune(val)[0]) {
		panic("Invalid dict tag name: " + val)
	}
	reader.consumeToken(ID)
	return val
}

func (reader *ZincReader) consumeNumber() Number {
	number := reader.curVal.(Number)
	reader.consumeToken(NUMBER)
	return number
}

func (reader *ZincReader) consumeStr() Str {
	str := reader.curVal.(Str)
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
	reader.peekVal = reader.tokenizer.val
	// reader.peekLine = reader.tokenizer.line
}
