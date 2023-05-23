package io

import (
	"errors"
	"strings"
	"unicode"

	"github.com/NeedleInAJayStack/haystack"
)

// ZincReader reads Zinc strings into Haystack Vals
type ZincReader struct {
	tokenizer Tokenizer

	cur    Token
	curVal haystack.Val
	// curLine int

	peek    Token
	peekVal haystack.Val
	// peekLine int
}

// InitString initializes with a specific string
func (reader *ZincReader) InitString(str string) {
	reader.Init(strings.NewReader(str))
}

// Init initializes by wrapping the input reader
func (reader *ZincReader) Init(in *strings.Reader) {
	reader.tokenizer = Tokenizer{}
	reader.tokenizer.Init(in)

	reader.consume()
	reader.consume()
}

// ReadVal proceeds through the next haystack.Val and returns it
func (reader *ZincReader) ReadVal() (haystack.Val, error) {
	var val haystack.Val
	var err error

	if reader.cur == ID {
		val, err = reader.parseGrid()
	} else {
		val, err = reader.parseVal()
	}
	if err != nil {
		return haystack.NewNull(), err
	}

	if reader.cur != EOF {
		return haystack.NewNull(), errors.New("Expecting EOF, not " + reader.cur.String())
	}
	return val, err
}

func (reader *ZincReader) parseVal() (haystack.Val, error) {
	if reader.cur == ID {
		id := reader.curVal.(haystack.Id).String()
		err := reader.consumeToken(ID)
		if err != nil {
			return haystack.NewNull(), err
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
		if id == "T" {
			return haystack.NewBool(true), nil
		} else if id == "F" {
			return haystack.NewBool(false), nil
		} else if id == "N" {
			return haystack.NewNull(), nil
		} else if id == "M" {
			return haystack.NewMarker(), nil
		} else if id == "NA" {
			return haystack.NewNA(), nil
		} else if id == "R" {
			return haystack.NewRemove(), nil
		} else if id == "NaN" {
			return haystack.NaN(), nil
		} else if id == "INF" {
			return haystack.Inf(), nil
		} else {
			return haystack.NewNull(), errors.New("unexpected identifier: " + id)
		}
	}

	// literals
	if reader.cur.IsLiteral() {
		return reader.parseLiteral()
	}

	// -INF
	if reader.cur == MINUS && reader.peekVal.ToZinc() == "INF" {
		minusErr := reader.consumeToken(MINUS)
		if minusErr != nil {
			return haystack.NewNull(), minusErr
		}
		idErr := reader.consumeToken(ID)
		if idErr != nil {
			return haystack.NewNull(), idErr
		}
		return haystack.NegInf(), nil
	}

	// nested collections
	if reader.cur == LBRACKET {
		return reader.parseList()
	} else if reader.cur == LBRACE {
		return reader.parseDict()
	} else if reader.cur == LT2 {
		return reader.parseGrid()
	}

	return haystack.NewNull(), errors.New("Unexpected token: " + reader.cur.String())
}

func (reader *ZincReader) parseCoord(id string) (haystack.Coord, error) {
	if id != "C" {
		return haystack.NewCoord(0, 0), errors.New("Expecting 'C' for coord, not " + id)
	}

	var lat haystack.Number
	var lng haystack.Number
	var err error

	err = reader.consumeToken(LPAREN)
	if err != nil {
		return haystack.NewCoord(0, 0), err
	}

	lat, err = reader.consumeNumber()
	if err != nil {
		return haystack.NewCoord(0, 0), err
	}

	err = reader.consumeToken(COMMA)
	if err != nil {
		return haystack.NewCoord(0, 0), err
	}

	lng, err = reader.consumeNumber()
	if err != nil {
		return haystack.NewCoord(0, 0), err
	}

	err = reader.consumeToken(RPAREN)
	if err != nil {
		return haystack.NewCoord(0, 0), err
	}

	return haystack.NewCoord(lat.Float(), lng.Float()), nil
}

func (reader *ZincReader) parseXStr(id string) (haystack.Val, error) {
	if !unicode.IsUpper([]rune(id)[0]) {
		return haystack.NewNull(), errors.New("Invalid XStr type: " + id)
	}
	if id == "Bin" { // I think Bins are obselete
		var mime haystack.Str
		var err error
		err = reader.consumeToken(LPAREN)
		if err != nil {
			return haystack.NewNull(), err
		}

		mime, err = reader.consumeStr()
		if err != nil {
			return haystack.NewNull(), err
		}

		err = reader.consumeToken(RPAREN)
		if err != nil {
			return haystack.NewNull(), err
		}

		return haystack.NewBin(mime.String()), nil
	} else {
		var val haystack.Str
		var err error
		err = reader.consumeToken(LPAREN)
		if err != nil {
			return haystack.NewNull(), err
		}

		val, err = reader.consumeStr()
		if err != nil {
			return haystack.NewNull(), err
		}

		err = reader.consumeToken(RPAREN)
		if err != nil {
			return haystack.NewNull(), err
		}

		return haystack.NewXStr(id, val.String()), nil
	}
}

func (reader *ZincReader) parseLiteral() (haystack.Val, error) {
	val := reader.curVal
	// Combine ref and dis
	if reader.cur == REF && reader.peek == STR {
		ref := reader.curVal.(haystack.Ref)
		dis := reader.peekVal.(haystack.Str)

		val = haystack.NewRef(ref.Id(), dis.String())
		err := reader.consumeToken(REF)
		if err != nil {
			return haystack.NewNull(), err
		}
	}
	reader.consume()
	return val, nil
}

func (reader *ZincReader) parseList() (haystack.List, error) {
	var vals []haystack.Val

	lbracketErr := reader.consumeToken(LBRACKET)
	if lbracketErr != nil {
		return haystack.NewList([]haystack.Val{}), lbracketErr
	}

	for reader.cur != RBRACKET && reader.cur != EOF {
		val, valErr := reader.parseVal()
		if valErr != nil {
			return haystack.NewList([]haystack.Val{}), valErr
		}

		vals = append(vals, val)
		if reader.cur != COMMA {
			break
		}
		reader.consumeToken(COMMA)
	}

	rbracketErr := reader.consumeToken(RBRACKET)
	if rbracketErr != nil {
		return haystack.NewList([]haystack.Val{}), rbracketErr
	}

	return haystack.NewList(vals), nil
}

func (reader *ZincReader) parseDict() (haystack.Dict, error) {
	items := make(map[string]haystack.Val)

	braces := reader.cur == LBRACE
	if braces {
		err := reader.consumeToken(LBRACE)
		if err != nil {
			return haystack.NewDict(map[string]haystack.Val{}), err
		}
	}
	for reader.cur == ID {
		var id string
		var val haystack.Val
		var err error

		id, err = reader.consumeTagName()
		if err != nil {
			return haystack.NewDict(map[string]haystack.Val{}), err
		}

		val = haystack.NewMarker() // Default to marker val if there is no value
		if reader.cur == COLON {
			reader.consumeToken(COLON)
			val, err = reader.parseVal()
			if err != nil {
				return haystack.NewDict(map[string]haystack.Val{}), err
			}
		}
		items[id] = val
	}
	if braces {
		err := reader.consumeToken(RBRACE)
		if err != nil {
			return haystack.NewDict(map[string]haystack.Val{}), err
		}
	}

	return haystack.NewDict(items), nil
}

func (reader *ZincReader) parseGrid() (haystack.Grid, error) {
	var err error

	gb := haystack.NewGridBuilder()

	nested := reader.cur == LT2
	if nested {
		err := reader.consumeToken(LT2)
		if err != nil {
			return haystack.EmptyGrid(), err
		}

		if reader.cur == NL {
			reader.consumeToken(NL)
		}
	}

	// ver:"3.0"
	if reader.cur != ID {
		return haystack.EmptyGrid(), errors.New("Expecting grid 'ver' identifier, not " + reader.curVal.ToZinc())
	}
	err = reader.consume()
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	err = reader.consumeToken(COLON)
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	ver, verErr := reader.consumeStr()
	if verErr != nil {
		return haystack.EmptyGrid(), verErr
	}
	err = checkVersion(ver.String())
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	// grid meta
	if reader.cur == ID {
		dict, err := reader.parseDict()
		if err != nil {
			return haystack.EmptyGrid(), err
		}
		gb.SetMetaDict(dict)
	}
	err = reader.consumeToken(NL)
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	// column definitions
	numCols := 0
	for reader.cur == ID {
		numCols = numCols + 1
		name, err := reader.consumeTagName()
		if err != nil {
			return haystack.EmptyGrid(), err
		}

		colMeta := haystack.EmptyDict()
		if reader.cur == ID {
			colMeta, err = reader.parseDict()
			if err != nil {
				return haystack.EmptyGrid(), err
			}
		}
		gb.AddColDict(name, colMeta)

		if reader.cur != COMMA {
			break
		}
		err = reader.consumeToken(COMMA)
		if err != nil {
			return haystack.EmptyGrid(), err
		}
	}
	if numCols == 0 {
		return haystack.EmptyGrid(), errors.New("no columns defined")
	}
	err = reader.consumeToken(NL)
	if err != nil {
		return haystack.EmptyGrid(), err
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
		var vals []haystack.Val
		for i := 0; i < numCols; i = i + 1 {
			if reader.cur == COMMA || reader.cur == NL || reader.cur == EOF {
				vals = append(vals, haystack.NewNull())
			} else {
				val, err := reader.parseVal()
				if err != nil {
					return haystack.EmptyGrid(), err
				}
				vals = append(vals, val)
			}
			if i+1 < numCols {
				err = reader.consumeToken(COMMA)
				if err != nil {
					return haystack.EmptyGrid(), err
				}
			}
		}
		gb.AddRow(vals)

		// newline or end
		if nested && reader.cur == GT2 {
			break
		} else if reader.cur == EOF {
			break
		}
		err = reader.consumeToken(NL)
		if err != nil {
			return haystack.EmptyGrid(), err
		}
	}

	if reader.cur == NL {
		reader.consumeToken(NL)
	}
	if nested {
		err = reader.consumeToken(GT2)
		if err != nil {
			return haystack.EmptyGrid(), err
		}
	}

	return gb.ToGrid(), nil
}

func (reader *ZincReader) consumeTagName() (string, error) {
	id := reader.curVal.(haystack.Id)
	val := id.String()
	if val == "" || unicode.IsUpper([]rune(val)[0]) {
		return "", errors.New("Invalid dict tag name: " + val)
	}
	err := reader.consumeToken(ID)
	if err != nil {
		return "", err
	}
	return val, nil
}

func checkVersion(str string) error {
	if str != "2.0" && str != "3.0" {
		return errors.New("Unsupported version: " + str)
	}
	return nil
}

func (reader *ZincReader) consumeNumber() (haystack.Number, error) {
	number := reader.curVal.(haystack.Number)
	err := reader.consumeToken(NUMBER)
	if err != nil {
		return haystack.NewNumber(0, ""), err
	}
	return number, nil
}

func (reader *ZincReader) consumeStr() (haystack.Str, error) {
	str := reader.curVal.(haystack.Str)
	err := reader.consumeToken(STR)
	if err != nil {
		return haystack.NewStr(""), err
	}
	return str, nil
}

func (reader *ZincReader) consumeToken(expected Token) error {
	if reader.cur != expected {
		return errors.New("Expected " + expected.String() + " not " + reader.cur.String())
	}
	reader.consume()
	return nil
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
	reader.peekVal = reader.tokenizer.Val()
	// reader.peekLine = reader.tokenizer.line
	return nil
}
