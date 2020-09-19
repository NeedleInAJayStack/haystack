package haystack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const runeEOF = -1

// Tokenizer implements a stream approach for reading Haystack formats such as Zinc, Trio, and Filters
type Tokenizer struct {
	in    *strings.Reader
	cur   rune
	peek  rune
	val   Val
	token Token
}

// Init initializes a tokenizer on a Reader
func (tokenizer *Tokenizer) Init(in *strings.Reader) {
	tokenizer.in = in
	tokenizer.cur = 0
	tokenizer.peek = 0
	tokenizer.val = NewNull()
	tokenizer.token = DEF

	tokenizer.consume()
	tokenizer.consume()
}

// Next reads the next haystack token, storing the value in the Val, and Token fields
func (tokenizer *Tokenizer) Next() (Token, error) {
	// reset
	tokenizer.val = Null{}

	// skip non-meaningful whitespace and comments
	//int startLine = line
	for true {
		// treat space, tab, non-breaking space as whitespace
		if tokenizer.cur == ' ' || tokenizer.cur == '\t' || tokenizer.cur == 0xa0 {
			tokenizer.consume()
			continue
		}

		// comments
		if tokenizer.cur == '/' {
			if tokenizer.peek == '/' {
				tokenizer.skipCommentsSingleLine()
				continue
			}
			if tokenizer.peek == '*' {
				tokenizer.skipCommentsMultiLine()
				continue
			}
		}

		break
	}

	var err error
	newToken := DEF
	if tokenizer.cur == runeEOF {
		newToken = EOF
	} else if tokenizer.cur == '\n' || tokenizer.cur == '\r' { // newlines
		if tokenizer.cur == '\r' && tokenizer.peek == '\n' {
			tokenizer.consumeRune('\r')
		}
		tokenizer.consume()
		//++line
		newToken = NL
	} else if isIdStart(tokenizer.cur) { // handle various starting chars
		newToken = tokenizer.id()
	} else if tokenizer.cur == '"' {
		newToken, err = tokenizer.str()
	} else if tokenizer.cur == '@' {
		newToken = tokenizer.ref()
	} else if unicode.IsDigit(tokenizer.cur) {
		newToken, err = tokenizer.digits()
	} else if tokenizer.cur == '`' {
		newToken, err = tokenizer.uri()
	} else if tokenizer.cur == '-' && unicode.IsDigit(tokenizer.peek) {
		newToken, err = tokenizer.digits()
	} else {
		newToken, err = tokenizer.symbol()
	}

	tokenizer.token = newToken
	return tokenizer.token, err
}

// Token production methods

func (tokenizer *Tokenizer) id() Token {
	buf := strings.Builder{}
	for isIdPart(tokenizer.cur) {
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
	}
	tokenizer.val = NewId(buf.String())
	return ID
}

func (tokenizer *Tokenizer) ref() Token {
	tokenizer.consumeRune('@')
	buf := strings.Builder{}
	for isRefPart(tokenizer.cur) {
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
	}
	tokenizer.val = NewRef(buf.String(), "")
	return REF
}

func (tokenizer *Tokenizer) str() (Token, error) {
	tokenizer.consumeRune('"')
	buf := strings.Builder{}
	for {
		if tokenizer.cur == runeEOF { // runeEOF
			return EOF, errors.New("Unexpected end of str")
		} else if tokenizer.cur == '"' {
			tokenizer.consumeRune('"')
			break
		} else if tokenizer.cur == '\\' {
			esc, escErr := tokenizer.escape()
			if escErr != nil {
				return STR, escErr
			}
			buf.WriteRune(esc)
			// continue
		}
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
	}
	tokenizer.val = NewStr(buf.String())
	return STR, nil
}

func (tokenizer *Tokenizer) uri() (Token, error) {
	tokenizer.consumeRune('`')
	buf := strings.Builder{}
	for true {
		if tokenizer.cur == '`' {
			tokenizer.consumeRune('`')
			break
		} else if tokenizer.cur == runeEOF {
			return EOF, errors.New("Unexpected end of uri: eof")
		} else if tokenizer.cur == '\n' {
			return URI, errors.New("Unexpected end of uri: newline")
		} else if tokenizer.cur == '\\' {
			if isUriEscapeIgnore(tokenizer.peek) {
				buf.WriteRune(tokenizer.cur)
				tokenizer.consume()
				buf.WriteRune(tokenizer.cur)
				tokenizer.consume()
			} else {
				char, err := tokenizer.escape()
				if err != nil {
					return URI, err
				}
				buf.WriteRune(char)
			}
		} else {
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
		}
	}

	tokenizer.val = NewUri(buf.String())
	return URI, nil
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
	} else if tokenizer.cur == 'u' { // check for \uxxxx (utf16 literal format)
		buf := strings.Builder{}
		tokenizer.consumeRune('u')
		buf.WriteRune(tokenizer.cur) // Get the next 4 characters
		tokenizer.consume()
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
		buf.WriteRune(tokenizer.cur)
		tokenizer.consume()
		buf.WriteRune(tokenizer.cur)

		codeResult, codeErr := strconv.ParseInt(buf.String(), 16, 0) // Force to base-16 for utf16 format
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
	if tokenizer.cur == '0' && tokenizer.peek == 'x' { // hex integer (no unit allowed)
		buf := strings.Builder{}
		buf.WriteRune(tokenizer.cur)
		tokenizer.consumeRune('0')
		buf.WriteRune(tokenizer.cur)
		tokenizer.consumeRune('x')
		for isHex(tokenizer.cur) || tokenizer.cur == '_' {
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
		}
		valInt, err := strconv.ParseInt(buf.String(), 0, 0) // ParseInt accepts hex format
		valFloat := float64(valInt)
		tokenizer.val = NewNumber(valFloat, "")
		return NUMBER, err
	}
	// consume all things that might be part of this number token
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
			} else if (exponential || colonCount >= 1) && tokenizer.cur == '+' {
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
		return tokenizer.dateTime(&buf)
	} else {
		return tokenizer.number(buf.String(), unitIndex)
	}
}

func (tokenizer *Tokenizer) dateTime(buf *strings.Builder) (Token, error) {
	// Format variable formats to: "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"

	// xxx timezone
	if tokenizer.cur != ' ' || !unicode.IsUpper(tokenizer.peek) {
		str := buf.String()
		if str[len(str)-1] == 'Z' {
			buf.WriteString(" UTC")
		} else {
			return DATETIME, errors.New("Expecting timezone")
		}
	} else {
		tokenizer.consume()
		buf.WriteRune(' ')
		for isIdPart(tokenizer.cur) {
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
		}

		// handle GMT+xx or GMT-xx
		if (tokenizer.cur == '+' || tokenizer.cur == '-') && strings.HasSuffix(buf.String(), "GMT") {
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
			for unicode.IsDigit(tokenizer.cur) {
				buf.WriteRune(tokenizer.cur)
				tokenizer.consume()
			}
		}
	}

	dateTime, err := NewDateTimeFromString(buf.String())
	tokenizer.val = dateTime
	return DATETIME, err
}

func (tokenizer *Tokenizer) date(str string) (Token, error) {
	date, err := NewDateFromString(str)
	tokenizer.val = date
	return DATE, err
}

func (tokenizer *Tokenizer) time(str string, addSeconds bool) (Token, error) {
	if addSeconds {
		str = str + ":00"
	}
	time, err := NewTimeFromString(str)
	tokenizer.val = time
	return TIME, err
}

func (tokenizer *Tokenizer) number(str string, unitIndex int) (Token, error) {
	if unitIndex == 0 {
		number, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return NUMBER, err
		}
		tokenizer.val = NewNumber(number, "")
		return NUMBER, err
	} else {
		numberStr := str[0:unitIndex]
		unit := str[unitIndex:len(str)]
		number, err := strconv.ParseFloat(numberStr, 64)
		if err != nil {
			return NUMBER, err
		}
		tokenizer.val = NewNumber(number, unit)
		return NUMBER, err
	}
}

func (tokenizer *Tokenizer) symbol() (Token, error) {
	c := tokenizer.cur
	tokenizer.consume()
	if c == ',' {
		return COMMA, nil
	} else if c == '/' {
		return SLASH, nil
	} else if c == ':' {
		return COLON, nil
	} else if c == ';' {
		return SEMICOLON, nil
	} else if c == '[' {
		return LBRACKET, nil
	} else if c == ']' {
		return RBRACKET, nil
	} else if c == '{' {
		return LBRACE, nil
	} else if c == '}' {
		return RBRACE, nil
	} else if c == '(' {
		return LPAREN, nil
	} else if c == ')' {
		return RPAREN, nil
	} else if c == '<' {
		if tokenizer.cur == '<' {
			tokenizer.consumeRune('<')
			return LT2, nil
		} else if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return LTEQ, nil
		} else {
			return LT, nil
		}
	} else if c == '>' {
		if tokenizer.cur == '>' {
			tokenizer.consumeRune('>')
			return GT2, nil
		} else if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return GTEQ, nil
		} else {
			return GT, nil
		}
	} else if c == '-' {
		if tokenizer.cur == '>' {
			tokenizer.consumeRune('>')
			return ARROW, nil
		} else {
			return MINUS, nil
		}
	} else if c == '=' {
		if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return EQ, nil
		}
		return ASSIGN, nil
	} else if c == '!' {
		if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return NOTEQ, nil
		}
		return BANG, nil
	} else if c == runeEOF {
		return EOF, nil
	}
	return EOF, errors.New("Unexpected symbol: '" + string(c) + "'") // Not sure what other than runeEOF to return
}

// Comments

func (tokenizer *Tokenizer) skipCommentsSingleLine() {
	tokenizer.consumeRune('/')
	tokenizer.consumeRune('/')
	for tokenizer.cur != '\n' && tokenizer.cur != runeEOF {
		tokenizer.consume()
	}
}

func (tokenizer *Tokenizer) skipCommentsMultiLine() {
	tokenizer.consumeRune('/')
	tokenizer.consumeRune('*')
	depth := 1
	for tokenizer.cur != runeEOF {
		if tokenizer.cur == '*' && tokenizer.peek == '/' {
			tokenizer.consumeRune('*')
			tokenizer.consumeRune('/')
			depth--
			if depth <= 0 {
				break
			}
		}
		if tokenizer.cur == '/' && tokenizer.peek == '*' {
			tokenizer.consumeRune('/')
			tokenizer.consumeRune('*')
			depth++
			continue
		}
		// if tokenizer.cur == '\n' {
		// 	++line
		// }
		tokenizer.consume()
	}
}

// Consume methods

func (tokenizer *Tokenizer) consume() error {
	var err error
	tokenizer.cur = tokenizer.peek
	tokenizer.peek, _, err = tokenizer.in.ReadRune()
	if err != nil { // If end-of-stream, indicate with val of runeEOF
		tokenizer.peek = runeEOF
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

func isRefPart(char rune) bool {
	if isIdPart(char) {
		return true
	} else if char == '-' || char == ':' || char == '.' || char == '~' {
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
