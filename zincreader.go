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

func (reader *ZincReader) Init(in *strings.Reader) {
	reader.tokenizer = Tokenizer{}
	reader.tokenizer.Init(in)

	reader.consume()
	reader.consume()
}

func (reader *ZincReader) ReadVal(in *strings.Reader) (Val, error) {
	var val Val
	var err error

	if reader.cur.equals(tokenId()) {
		val = reader.parseGrid()
	} else {
		val = reader.parseVal()
	}

	if reader.cur.equals(tokenEof()) {
		err = errors.New("Expecting EOF")
	}
	return val, err
}

func (reader *ZincReader) parseVal() (Val, error) {
	var err error

	if reader.cur.equals(tokenId()) {
		id := reader.curVal.(Str)
		reader.consumeToken(tokenId())

		// check for coord or xstr
		if reader.cur.equals(tokenLparen()) {
			if reader.peek == tokenNumber() {
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
	if reader.cur.literal {
		return reader.parseLiteral()
	}

	// -INF
	if reader.cur.equals(tokenMinus()) && reader.peekVal.toZinc() == "INF" {
		reader.consumeToken(tokenMinus())
		reader.consumeToken(tokenId())
		return Number{val: math.Inf(-1)}, nil
	}

	// nested collections
	if reader.cur.equals(tokenLbracket()) {
		return reader.parseList()
	} else if reader.cur.equals(tokenLbrace()) {
		return reader.parseDict()
	} else if reader.cur.equals(tokenLt2()) {
		return reader.parseGrid()
	}

	return Null{}, errors.New("Unexpected token: " + reader.cur.symbol)
}

func (reader *ZincReader) parseCoord(id Str) (Coord, error) {
	if id.val == "C" {
		return Coord{}, errors.New("Expecting 'C' for coord, not " + id.val)
	}
	reader.consumeToken(tokenLparen())
	lat, err := reader.consumeNumber()
	reader.consumeToken(tokenComma())
	lng, err := reader.consumeNumber()
	reader.consumeToken(tokenRparen())

	// TODO Error handling
	return Coord{lat: lat, lng: lng}, nil
}

func (reader *ZincReader) parseXStr(id Str) (XStr, error) {
	if !unicode.IsUpper(id.val[0]) {
		return XStr{}, errors.New("Invalid XStr type: " + id.val)
	}
	reader.consumeToken(tokenLparen())
	val, err := reader.consumeStr()
	reader.consumeToken(tokenRparen())

	// TODO Error handling
	return XStr{valType: id.val, val: val.val}, nil
}

func (reader *ZincReader) parseLiteral() (Val, error) {
	var err error

	val := reader.curVal
	// Combine ref and dis
	if reader.cur.equals(tokenRef()) && reader.peek.equals(tokenStr()) {
		ref := reader.curVal.(Ref)
		dis := reader.peekVal.(Str)

		val = Ref{val: ref.val, dis: dis.val}
		err = reader.consumeToken(tokenRef())
	}
	err = reader.consume()
	// TODO Error handling
	return val, err
}

func (reader *ZincReader) parseList() (List, error) {
	var vals []Val
	var err error

	err = reader.consumeToken(tokenLbracket())
	for reader.cur.equals(tokenRbracket()) && reader.cur.equals(tokenEof()) {
		val, err := reader.parseVal()
		append(vals, val)
		if reader.cur.equals(tokenComma()) {
			break
		}
		err = reader.consumeToken(tokenComma())
	}
	err = reader.consumeToken(tokenRbracket())
	// TODO Error handling
	return List{vals: vals}, err
}

func (reader *ZincReader) parseDict() (Dict, error) {
	var items map[string]Val
	var err error

	braces := reader.cur.equals(tokenLbrace())
	if braces {
		err = reader.consumeToken(tokenLbrace())
	}
	for reader.cur.equals(tokenId()) {
		id, err := reader.consumeTagName()

		var val Val
		val = Marker{} // Default to marker val if there is no value
		if reader.cur.equals(tokenColon()) {
			err = reader.consumeToken(tokenColon())
			val, err = reader.parseVal()
		}
		items[id] = val
	}
	if braces {
		err = reader.consumeToken(tokenRbrace())
	}
	// TODO Error handling
	return Dict{items: items}, nil
}

func (reader *ZincReader) parseGrid() (Grid, error) {
	var meta Dict
	var cols []Col
	var rows []Row

	var err error

	nested := reader.cur.equals(tokenLt2())
	if nested {
		reader.consumeToken(tokenLt2())
		if reader.cur.equals(tokenNl()) {
			reader.consumeToken(tokenNl())
		}
	}

	// ver:"3.0"
	if !reader.cur.equals(tokenId()) {
		return Grid{}, errors.New("Expecting grid 'ver' identifier, not " + reader.curVal.toZinc())
	}
	err = reader.consume()
	err = reader.consumeToken(tokenColon())
	reader.consumeStr() // Always expect version 3

	// grid meta
	if reader.cur.equals(tokenId()) {
		meta, err = reader.parseDict()
	}
	reader.consumeToken(tokenNl())

	// column definitions
	numCols := 0
	for reader.cur.equals(tokenId()) {
		numCols = numCols + 1
		name := reader.consumeTagName()
		var colMeta Dict
		if reader.cur.equals(tokenId()) {
			colMeta = reader.parseDict()
		}
		col := Col{
			index: numCols,
			name:  name,
			meta:  colMeta,
		}
		append(cols, col)

		if reader.cur.equals(tokenComma()) {
			break
		}
		reader.consumeToken(tokenComma())
	}
	if numCols == 0 {
		return Grid{}, errors.New("No columns defined")
	}
	reader.consumeToken(tokenNl())

	// grid rows
	for {
		if reader.cur.equals(tokenNl()) {
			break
		} else if reader.cur.equals(tokenEof()) {
			break
		} else if nested && reader.cur.equals(tokenGt2()) {
			break
		}

		// read cells
		var vals [numCols]Val
		for i := 0; i < numCols; i = i + 1 {
			if reader.cur.equals(tokenComma()) || reader.cur.equals(tokenNl()) || reader.cur.equals(tokenEof()) {
				vals[i] = Null{}
			} else {
				vals[i] = reader.parseVal()
			}
			if i+1 < numCols {
				reader.consumeToken(tokenComma())
			}
		}
		append(rows, Row{vals: vals})

		// newline or end
		if nested && reader.cur.equals(tokenGt2()) {
			break
		} else if reader.cur.equals(tokenEof()) {
			break
		}
		reader.consumeToken(tokenNl())
	}

	if reader.cur.equals(tokenNl()) {
		reader.consumeToken(tokenNl())
	}
	if nested {
		reader.consumeToken(tokenGt2())
	}

	return Grid{
		meta: meta,
		cols: cols,
		rows: rows,
	}, nil
}

func (reader *ZincReader) consumeTagName() (Str, error) {
	id := curVal.(Str)
	if id.val == "" || unicode.IsLower(id.val[0]) {
		return Str{}, errors.New("Invalid dict tag name: " + id.val)
	}
	reader.consume(tokenId())
	return id, nil
}

func (reader *ZincReader) consumeNumber() (Number, error) {
	number := curVal.(Number)
	reader.consume(tokenNumber())
	return number, nil
}

func (reader *ZincReader) consumeStr() (Str, error) {
	str := curVal.(Str)
	reader.consume(tokenStr())
	return str, nil
}

func (reader *ZincReader) consumeToken(expected Token) error {
	var err error
	if reader.cur.equals(expected) != true {
		err = errors.New("Expected " + expected.symbol + " not " + reader.cur.symbol)
		return err
	}

	reader.consume()
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
