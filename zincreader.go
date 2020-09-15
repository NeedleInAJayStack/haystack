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
	lat := reader.consumeNumber()
	reader.consumeToken(tokenComma())
	lng := reader.consumeNumber()
	reader.consumeToken(tokenRparen())

	// TODO Error handling
	return Coord{lat: lat, lng: lng}, nil
}

func (reader *ZincReader) parseXStr(id Str) (XStr, error) {
	if !unicode.IsUpper(id.val[0]) {
		return XStr{}, errors.New("Invalid XStr type: " + id.val)
	}
	reader.consumeToken(tokenLparen())
	val := reader.consumeStr()
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
	return val, err
}

func (reader *ZincReader) parseList() (List, error) {
	var arr []Val
	var err error

	err = reader.consumeToken(tokenLbracket)
	for reader.cur.equals(tokenRbracket()) && reader.cur.equals(tokenEof()) {
		val, err := reader.parseVal()
		append(arr, val)
		if reader.cur.equals(tokenComma()) {
			break
		}
		reader.consumeToken(tokenComma())
	}
	reader.consumeToken(tokenRbracket())
	return List{vals: arr}
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
