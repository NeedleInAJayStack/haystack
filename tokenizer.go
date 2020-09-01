package haystack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// Stream based tokenizer for Haystack formats such as Zinc and Filters
type Tokenizer struct {
	in   strings.Reader
	cur  rune
	peek rune
	val  Val
}

// Consume methods

func (tokenizer *Tokenizer) consume() error {
	var err error
	tokenizer.cur = tokenizer.peek
	tokenizer.peek, _, err = tokenizer.in.ReadRune()
	// TODO adjust error to handle eof
	return err
}

func (tokenizer *Tokenizer) consumeRune(expected rune) error {
	if tokenizer.cur != expected {
		return errors.New("Expected " + string(expected))
	}
	tokenizer.consume()
	return nil
}

// Token methods

func (tokenizer *Tokenizer) id() Token {
	buf := strings.Builder{}
	for isIdPart(tokenizer.cur) {
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
	}
	tokenizer.val = &Str{val: buf.String()}
	return tokenId()
}

func (tokenizer *Tokenizer) ref() Token {
	tokenizer.consumeRune('@')
	buf := strings.Builder{}
	for isIdPart(tokenizer.cur) {
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
	}
	tokenizer.val = &Ref{val: buf.String()}
	return tokenRef()
}

func (tokenizer *Tokenizer) digits() (Token, error) {
	if tokenizer.cur == '0' && tokenizer.peek == 'x' { // hex number (no unit allowed)
		tokenizer.consumeRune('0')
		tokenizer.consumeRune('x')
		buf := strings.Builder{}
		for isHex(tokenizer.cur) || tokenizer.cur == '_' {
			if isHex(tokenizer.cur) {
				buf.WriteRune(tokenizer.cur)
			}
			tokenizer.consume()
		}
		float, err := strconv.ParseFloat(buf.String(), 64)
		tokenizer.val = &Number{val: float}
		return tokenNumber(), err
	} else { // consume all things that might be part of this number token
		buf := strings.Builder{}
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()

		colonCount := 0
		dashCount := 0
		exponential := false
		unitIndex := 0 // Determines unit location in the token

		for {
			if !unicode.IsDigit(tokenizer.cur) {
				if exponential && isSign(tokenizer.cur) {
					// Just fall through
				} else if tokenizer.cur == '-' {
					dashCount = dashCount + 1
				} else if tokenizer.cur == ':' && unicode.IsDigit(tokenizer.peek) {
					colonCount = colonCount + 1
				} else if exponential || (colonCount >= 1 && tokenizer.cur == '+') {
					// Just fall through
				} else if tokenizer.cur == '.' {
					if !unicode.IsDigit(tokenizer.peek) { // Break numbers at the following decimal
						break
					}
					// Keep reading if the demical is followed by digits
				} else if unicode.ToLower(tokenizer.cur) == 'e' && (isSign(tokenizer.peek) || unicode.IsDigit(tokenizer.peek)) {
					exponential = true
				} else if isUnit(tokenizer.cur) {
					if unitIndex == 0 {
						unitIndex = buf.Len()
					}
				} else if tokenizer.cur == '_' {
					if unitIndex == 0 {
						if unicode.IsDigit(tokenizer.peek) { // Skip underscores grouping digits
							tokenizer.consume()
						} else { // If not a digit, it's a custom unit
							unitIndex = buf.Len()
						}
					}
				} else {
					break
				}
			}
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
		}

		if dashCount == 2 && colonCount == 0 {
			return tokenizer.date(buf.String())
		} else if dashCount == 0 && colonCount >= 1 {
			return tokenizer.time(buf.String(), colonCount == 1)
		} else if dashCount >= 2 {
			// TODO Implement
			// return tokenizer.dateTime(buf.String())
			return tokenDateTime(), nil
		} else {
			// TODO Implement
			// return tokenizer.number(buf.String(), unitIndex)
			return tokenNumber(), nil
		}
	}
}

func (tokenizer *Tokenizer) date(str string) (Token, error) {
	date, err := dateFromZinc(str)
	tokenizer.val = &date
	return tokenDate(), err
}

func (tokenizer *Tokenizer) time(str string, addSeconds bool) (Token, error) {
	if addSeconds {
		str = str + ":00"
	}
	time, err := timeFromZinc(str)
	tokenizer.val = &time
	return tokenTime(), err
}

// Rune detection methods. These add onto those in unicode

func isSign(char rune) bool {
	return char == '-' || char == '+'
}

func isUnit(char rune) bool {
	return unicode.IsLetter(char) || char == '%' || char == '$' || char == '/' || char > 128
}

func isHex(char rune) bool {
	char = unicode.ToLower(char)
	if 'a' <= char && char <= 'f' {
		return true
	} else if unicode.IsDigit(char) {
		return true
	} else {
		return false
	}
}

func isIdStart(char rune) bool {
	if 'a' <= char && char <= 'z' {
		return true
	} else if 'A' <= char && char <= 'Z' {
		return true
	} else {
		return false
	}
}

func isIdPart(char rune) bool {
	if isIdStart(char) {
		return true
	} else if unicode.IsDigit(char) {
		return true
	} else if char == '_' {
		return true
	} else {
		return false
	}
}
