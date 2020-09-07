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
	cur  rune // -1 indicates end-of-stream
	peek rune // -1 indicates end-of-stream
	val  Val
}

// Consume methods

func (tokenizer *Tokenizer) consume() error {
	var err error
	tokenizer.cur = tokenizer.peek
	tokenizer.peek, _, err = tokenizer.in.ReadRune()
	if err != nil { // If end-of-stream, indicate with val of -1
		tokenizer.cur = -1
		tokenizer.peek = -1
	}
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

func (tokenizer *Tokenizer) str() (Token, error) {
	tokenizer.consumeRune('"')
	buf := strings.Builder{}
	for true {
		if tokenizer.cur == -1 {
			return tokenEof(), errors.New("Unexpected end of str")
		} else if tokenizer.cur == '"' {
			tokenizer.consumeRune('"')
			break
		} else if tokenizer.cur == '\\' {
			esc, escErr := tokenizer.escape()
			if escErr != nil {
				return tokenStr(), escErr
			}
			buf.WriteRune(esc)
			// continue
		}
	}
	tokenizer.val = &Str{val: buf.String()}
	return tokenStr(), nil
}

func (tokenizer *Tokenizer) uri() (Token, error) {
	tokenizer.consumeRune('`')
	buf := strings.Builder{}
	for true {
		if tokenizer.cur == '`' {
			tokenizer.consumeRune('`')
			break
		} else if tokenizer.cur == -1 {
			return tokenEof(), errors.New("Unexpected end of uri: eof")
		} else if tokenizer.cur == '\n' {
			return tokenUri(), errors.New("Unexpected end of uri: newline")
		} else if tokenizer.cur == '\\' {
			if isUriEscapeIgnore(tokenizer.peek) {
				buf.WriteRune(tokenizer.cur)
				tokenizer.consume()
				buf.WriteRune(tokenizer.cur)
				tokenizer.consume()
				break
			} else {
				char, err := tokenizer.escape()
				if err != nil {
					return tokenUri(), err
				}
				buf.WriteRune(char)
			}
		} else {
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
		}
	}

	tokenizer.val = &Uri{val: buf.String()}
	return tokenUri(), nil
}

func (tokenizer *Tokenizer) escape() (rune, error) {
	tokenizer.consumeRune('\\')
	var result rune
	var err error
	if tokenizer.cur == 'b' {
		result = '\b'
	} else if tokenizer.cur == 'f' {
		result = '\f'
	} else if tokenizer.cur == 'n' {
		result = '\n'
	} else if tokenizer.cur == 'r' {
		result = '\r'
	} else if tokenizer.cur == 't' {
		result = '\t'
	} else if tokenizer.cur == '"' {
		result = '"'
	} else if tokenizer.cur == '$' {
		result = '$'
	} else if tokenizer.cur == '\'' {
		result = '\''
	} else if tokenizer.cur == '`' {
		result = '`'
	} else if tokenizer.cur == '\\' {
		result = '\\'
	} else if tokenizer.cur == 'u' { // check for \uxxxx
		buf := strings.Builder{}
		tokenizer.consumeRune('u')
		buf.WriteRune(tokenizer.cur) // Get the next 4 characters
		tokenizer.consume()
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
		buf.WriteRune(tokenizer.cur)
		// Wait to consume until we return

		codeResult, codeErr := strconv.ParseInt(buf.String(), 0, 32) // ParseFloat accepts hex format
		if codeErr != nil {
			err = codeErr
		} else {
			result = int32(codeResult)
		}
	}

	if result == 0 && err == nil {
		err = errors.New("Invalid escape sequence: " + string(tokenizer.cur))
	}
	tokenizer.consume()
	return result, err
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
		float, err := strconv.ParseFloat(buf.String(), 64) // ParseFloat accepts hex format
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
			return tokenizer.number(buf.String(), unitIndex)
		}
	}
}

// TODO implement dateTimeFromZinc
// func (tokenizer *Tokenizer) dateTime(buf strings.Builder) (Token, error) {
// 	// Format variable formats to: "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"

// 	// xxx timezone
// 	if tokenizer.cur != ' ' || !unicode.IsUpper(tokenizer.peek) {
// 		str := buf.String()
// 		if str[len(str)-1] == 'Z' {
// 			buf.WriteString(" UTC")
// 		} else {
// 			return tokenDateTime(), errors.New("Expecting timezone")
// 		}
// 	} else {
// 		tokenizer.consume()
// 		buf.WriteRune(' ')
// 		for isIdPart(tokenizer.cur) {
// 			buf.WriteRune(tokenizer.cur)
// 			tokenizer.consume()
// 		}

// 		// handle GMT+xx or GMT-xx
// 		if (tokenizer.cur == '+' || tokenizer.cur == '-') && strings.HasSuffix(buf.String(), "GMT") {
// 			buf.WriteRune(tokenizer.cur)
// 			tokenizer.consume()
// 			for unicode.IsDigit(tokenizer.cur) {
// 				buf.WriteRune(tokenizer.cur)
// 				tokenizer.consume()
// 			}
// 		}
// 	}

// 	dateTime, err := dateTimeFromZinc(buf.String())
// 	tokenizer.val = &dateTime
// 	return tokenDateTime(), err
// }

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

func (tokenizer *Tokenizer) number(str string, unitIndex int) (Token, error) {
	if unitIndex == 0 {
		number, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return tokenNumber(), err
		} else {
			tokenizer.val = &Number{val: number}
			return tokenNumber(), err
		}
	} else {
		numberStr := str[0:unitIndex]
		unit := str[unitIndex:len(str)]
		number, err := strconv.ParseFloat(numberStr, 64)
		if err != nil {
			return tokenNumber(), err
		} else {
			tokenizer.val = &Number{val: number, unit: unit}
			return tokenNumber(), err
		}
	}
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

func isUriEscapeIgnore(char rune) bool {
	return char == ':' ||
		char == '/' ||
		char == '?' ||
		char == '#' ||
		char == '[' ||
		char == ']' ||
		char == '@' ||
		char == '\\' ||
		char == '&' ||
		char == '=' ||
		char == ';'
}
