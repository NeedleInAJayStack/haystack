package haystack

import (
	"errors"
	"strings"
)

type Tokenizer struct {
	in   strings.Reader
	cur  rune
	peek rune
	val  Val
}

func (tokenizer *Tokenizer) consume() error {
	var err error
	tokenizer.cur = tokenizer.peek
	tokenizer.peek, _, err = tokenizer.in.ReadRune()
	// TODO adjust error to handle eof
	return err
}

func (tokenizer *Tokenizer) consumeChar(expected rune) error {
	if tokenizer.cur != expected {
		return errors.New("Expected " + string(expected))
	}
	tokenizer.consume()
	return nil
}

func (tokenizer *Tokenizer) ref() Token {
	tokenizer.consumeChar('@')
	buf := strings.Builder{}
	for isIdPart(tokenizer.cur) {
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
	}
	tokenizer.val = &Ref{val: buf.String()}
	return refToken()
}

func isDigit(cur rune) bool {
	return '0' <= cur && cur <= '9'
}

func isIdStart(cur rune) bool {
	if 'a' <= cur && cur <= 'z' {
		return true
	} else if 'A' <= cur && cur <= 'Z' {
		return true
	} else {
		return false
	}
}

func isIdPart(cur rune) bool {
	if isIdStart(cur) {
		return true
	} else if isDigit(cur) {
		return true
	} else if cur == '_' {
		return true
	} else {
		return false
	}
}
