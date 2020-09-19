package haystack

import (
	"errors"
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

func (reader *ZincReader) ReadVal() (Val, error) {
	var val Val
	var err error

	if reader.cur == ID {
		val, err = reader.parseGrid()
	} else {
		val, err = reader.parseVal()
	}
	if err != nil {
		return Null{}, err
	}

	if reader.cur == EOF {
		err = errors.New("Expecting EOF")
	}
	return val, err
}

func (reader *ZincReader) parseVal() (Val, error) {
	var err error

	if reader.cur == ID {
		id := reader.curVal.(Str)
		err = reader.consumeToken(ID)
		if err != nil {
			return Null{}, err
		}

		// check for coord or xstr
		if reader.cur == LPAREN {
			if reader.peek == NUMBER {
				return reader.parseCoord(id)
			} else {
				return reader.parseXStr(id)
			}
		}

		// check for keyword
		if id.val == "T" {
			return Bool{val: true}, nil
		} else if id.val == "F" {
			return Bool{val: false}, nil
		} else if id.val == "N" {
			return Null{}, nil
		} else if id.val == "M" {
			return Marker{}, nil
		} else if id.val == "NA" {
			return NA{}, nil
		} else if id.val == "R" {
			return Remove{}, nil
		} else if id.val == "NaN" {
			return Number{val: math.NaN()}, nil
		} else if id.val == "INF" {
			return Number{val: math.Inf(1)}, nil
		} else {
			return Null{}, errors.New("Unexpected identifier: " + id.val)
		}
	}

	// literals
	if reader.cur.IsLiteral() {
		return reader.parseLiteral()
	}

	// -INF
	if reader.cur == MINUS && reader.peekVal.toZinc() == "INF" {
		err = reader.consumeToken(MINUS)
		if err != nil {
			return Null{}, err
		}
		err = reader.consumeToken(ID)
		if err != nil {
			return Null{}, err
		}
		return Number{val: math.Inf(-1)}, nil
	}

	// nested collections
	if reader.cur == LBRACKET {
		return reader.parseList()
	} else if reader.cur == LBRACE {
		return reader.parseDict()
	} else if reader.cur == LT2 {
		return reader.parseGrid()
	}

	return Null{}, errors.New("Unexpected token: " + reader.cur.String())
}

func (reader *ZincReader) parseCoord(id Str) (Coord, error) {
	if id.val == "C" {
		return Coord{}, errors.New("Expecting 'C' for coord, not " + id.val)
	}

	var err error
	var lat Number
	var lng Number
	err = reader.consumeToken(LPAREN)
	if err != nil {
		return Coord{}, err
	}
	lat, err = reader.consumeNumber()
	if err != nil {
		return Coord{}, err
	}
	err = reader.consumeToken(COMMA)
	if err != nil {
		return Coord{}, err
	}
	lng, err = reader.consumeNumber()
	if err != nil {
		return Coord{}, err
	}
	err = reader.consumeToken(RPAREN)
	if err != nil { // I hate go error handling so much
		return Coord{}, err
	}

	return Coord{lat: lat.val, lng: lng.val}, err
}

func (reader *ZincReader) parseXStr(id Str) (XStr, error) {
	if !unicode.IsUpper([]rune(id.val)[0]) {
		return XStr{}, errors.New("Invalid XStr type: " + id.val)
	}

	var err error
	var val Str
	err = reader.consumeToken(LPAREN)
	if err != nil {
		return XStr{}, err
	}
	val, err = reader.consumeStr()
	if err != nil {
		return XStr{}, err
	}
	err = reader.consumeToken(RPAREN)
	if err != nil {
		return XStr{}, err
	}

	return XStr{valType: id.val, val: val.val}, err
}

func (reader *ZincReader) parseLiteral() (Val, error) {
	var err error

	val := reader.curVal
	// Combine ref and dis
	if reader.cur == REF && reader.peek == STR {
		ref := reader.curVal.(Ref)
		dis := reader.peekVal.(Str)

		val = Ref{val: ref.val, dis: dis.val}
		err = reader.consumeToken(REF)
		if err != nil {
			return Null{}, err
		}
	}
	err = reader.consume()
	return val, err
}

func (reader *ZincReader) parseList() (List, error) {
	var vals []Val
	var err error

	err = reader.consumeToken(LBRACKET)
	if err != nil {
		return List{}, err
	}
	for reader.cur != RBRACKET && reader.cur != EOF {
		var val Val
		val, err = reader.parseVal()
		if err != nil {
			break
		}
		vals = append(vals, val)
		if reader.cur == COMMA {
			break
		}
		err = reader.consumeToken(COMMA)
		if err != nil {
			break
		}
	}
	if err != nil {
		return List{}, err
	}

	err = reader.consumeToken(RBRACKET)

	return List{vals: vals}, err
}

func (reader *ZincReader) parseDict() (Dict, error) {
	var items map[string]Val
	var err error

	braces := reader.cur == LBRACE
	if braces {
		err = reader.consumeToken(LBRACE)
		if err != nil {
			return Dict{}, err
		}
	}
	for reader.cur == ID {
		var id Str
		var val Val

		id, err = reader.consumeTagName()
		if err != nil {
			break
		}

		val = Marker{} // Default to marker val if there is no value
		if reader.cur == COLON {
			err = reader.consumeToken(COLON)
			if err != nil {
				break
			}
			val, err = reader.parseVal()
			if err != nil {
				break
			}
		}
		items[id.val] = val
	}
	if err != nil {
		return Dict{}, err
	}

	if braces {
		err = reader.consumeToken(RBRACE)
		if err != nil {
			return Dict{}, err
		}
	}

	return Dict{items: items}, err
}

func (reader *ZincReader) parseGrid() (Grid, error) {
	var meta Dict
	var cols []Col
	var rows []Row

	var err error

	nested := reader.cur == LT2
	if nested {
		err = reader.consumeToken(LT2)
		if err != nil {
			return Grid{}, err
		}
		if reader.cur == NL {
			err = reader.consumeToken(NL)
			if err != nil {
				return Grid{}, err
			}
		}
	}

	// ver:"3.0"
	if reader.cur != ID {
		return Grid{}, errors.New("Expecting grid 'ver' identifier, not " + reader.curVal.toZinc())
	}
	err = reader.consume()
	if err != nil {
		return Grid{}, err
	}
	err = reader.consumeToken(COLON)
	if err != nil {
		return Grid{}, err
	}
	_, err = reader.consumeStr() // Always expect version 3
	if err != nil {
		return Grid{}, err
	}

	// grid meta
	if reader.cur == ID {
		meta, err = reader.parseDict()
		if err != nil {
			return Grid{}, err
		}
	}
	err = reader.consumeToken(NL)
	if err != nil {
		return Grid{}, err
	}

	// column definitions
	numCols := 0
	for reader.cur == ID {
		numCols = numCols + 1
		var name Str
		name, err = reader.consumeTagName()
		if err != nil {
			break
		}

		var colMeta Dict
		if reader.cur == ID {
			colMeta, err = reader.parseDict()
			if err != nil {
				break
			}
		}
		col := Col{
			index: numCols,
			name:  name.val,
			meta:  colMeta,
		}
		cols = append(cols, col)

		if reader.cur == COMMA {
			break
		}
		err = reader.consumeToken(COMMA)
		if err != nil {
			break
		}
	}
	if err != nil {
		return Grid{}, err
	}
	if numCols == 0 {
		return Grid{}, errors.New("No columns defined")
	}
	err = reader.consumeToken(NL)
	if err != nil {
		return Grid{}, err
	}

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
		vals := make([]Val, numCols)
		for i := 0; i < numCols; i = i + 1 {
			if reader.cur == COMMA || reader.cur == NL || reader.cur == EOF {
				vals[i] = Null{}
			} else {
				vals[i], err = reader.parseVal()
				if err != nil {
					break
				}
			}
			if i+1 < numCols {
				err = reader.consumeToken(COMMA)
				if err != nil {
					break
				}
			}
		}
		if err != nil {
			break
		}
		rows = append(rows, Row{vals: vals})

		// newline or end
		if nested && reader.cur == GT2 {
			break
		} else if reader.cur == EOF {
			break
		}
		err = reader.consumeToken(NL)
		if err != nil {
			break
		}
	}
	if err != nil {
		return Grid{}, err
	}

	if reader.cur == NL {
		err = reader.consumeToken(NL)
		if err != nil {
			return Grid{}, err
		}
	}
	if nested {
		err = reader.consumeToken(GT2)
		if err != nil {
			return Grid{}, err
		}
	}

	return Grid{
		meta: meta,
		cols: cols,
		rows: rows,
	}, err
}

func (reader *ZincReader) consumeTagName() (Str, error) {
	id := reader.curVal.(Str)
	if id.val == "" || unicode.IsLower([]rune(id.val)[0]) {
		return Str{}, errors.New("Invalid dict tag name: " + id.val)
	}
	err := reader.consumeToken(ID)
	return id, err
}

func (reader *ZincReader) consumeNumber() (Number, error) {
	number := reader.curVal.(Number)
	err := reader.consumeToken(NUMBER)
	return number, err
}

func (reader *ZincReader) consumeStr() (Str, error) {
	str := reader.curVal.(Str)
	err := reader.consumeToken(STR)
	return str, err
}

func (reader *ZincReader) consumeToken(expected Token) error {
	var err error
	if reader.cur != expected {
		err = errors.New("Expected " + expected.String() + " not " + reader.cur.String())
		return err
	}

	err = reader.consume()
	return err
}

func (reader *ZincReader) consume() error {
	newToken, err := reader.tokenizer.Next()
	if err != nil {
		return err
	}

	reader.cur = reader.peek
	reader.curVal = reader.peekVal
	// reader.curLine = reader.peekLine

	reader.peek = newToken
	reader.peekVal = reader.tokenizer.val
	// reader.peekLine = reader.tokenizer.line

	return err
}
