package haystack

import (
	"io"
	"strconv"
	"strings"
	"unicode"
)

const runeEOF = -1

// Tokenizer implements a stream approach for reading Haystack formats such as Zinc, Trio, and Filters
type Tokenizer struct {
	in    io.RuneReader
	cur   rune
	peek  rune
	val   Val
	token Token
}

// InitString initializes a tokenizer on an in-memory string
func (tokenizer *Tokenizer) InitString(str string) {
	tokenizer.Init(strings.NewReader(str))
}

// Init initializes a tokenizer on a Reader
func (tokenizer *Tokenizer) Init(in io.RuneReader) {
	tokenizer.in = in
	tokenizer.cur = 0
	tokenizer.peek = 0
	tokenizer.val = NewNull()
	tokenizer.token = DEF

	tokenizer.consume()
	tokenizer.consume()
}

// Val returns the current tokenizer value
func (tokenizer *Tokenizer) Val() Val {
	return tokenizer.val
}

// Next reads the next token, storing the value in the Val, and Token fields
func (tokenizer *Tokenizer) Next() Token {
	// reset
	tokenizer.val = NewNull()

	// skip non-meaningful whitespace and comments
	//int startLine = line
	for {
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
		newToken = tokenizer.str()
	} else if tokenizer.cur == '@' {
		newToken = tokenizer.ref()
	} else if unicode.IsDigit(tokenizer.cur) {
		newToken = tokenizer.digits()
	} else if tokenizer.cur == '`' {
		newToken = tokenizer.uri()
	} else if tokenizer.cur == '-' && unicode.IsDigit(tokenizer.peek) {
		newToken = tokenizer.digits()
	} else {
		newToken = tokenizer.symbol()
	}

	tokenizer.token = newToken
	return tokenizer.token
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

func (tokenizer *Tokenizer) str() Token {
	tokenizer.consumeRune('"')
	buf := strings.Builder{}
	for {
		if tokenizer.cur == runeEOF { // runeEOF
			panic("Unexpected end of str")
		} else if tokenizer.cur == '\\' {
			esc := tokenizer.escape()
			buf.WriteRune(esc)
			// continue
		} else if tokenizer.cur == '"' {
			tokenizer.consumeRune('"')
			break
		} else {
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
		}
	}
	tokenizer.val = NewStr(buf.String())
	return STR
}

func (tokenizer *Tokenizer) uri() Token {
	tokenizer.consumeRune('`')
	buf := strings.Builder{}
	for {
		if tokenizer.cur == '`' {
			tokenizer.consumeRune('`')
			break
		} else if tokenizer.cur == runeEOF {
			panic("Unexpected end of uri: eof")
		} else if tokenizer.cur == '\n' {
			panic("Unexpected end of uri: newline")
		} else if tokenizer.cur == '\\' {
			if isUriEscapeIgnore(tokenizer.peek) {
				buf.WriteRune(tokenizer.cur)
				tokenizer.consume()
				buf.WriteRune(tokenizer.cur)
				tokenizer.consume()
			} else {
				char := tokenizer.escape()
				buf.WriteRune(char)
			}
		} else {
			buf.WriteRune(tokenizer.cur)
			tokenizer.consume()
		}
	}

	tokenizer.val = NewUri(buf.String())
	return URI
}

func (tokenizer *Tokenizer) escape() rune {
	tokenizer.consumeRune('\\')
	var result rune
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
			panic(codeErr)
		} else {
			result = int32(codeResult)
		}
	}

	if result == 0 {
		panic("Invalid escape sequence: " + string(tokenizer.cur))
	}
	tokenizer.consume()
	return result
}

func (tokenizer *Tokenizer) digits() Token {
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
		if err != nil {
			panic(err)
		}
		valFloat := float64(valInt)
		tokenizer.val = NewNumber(valFloat, "")
		return NUMBER
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

func (tokenizer *Tokenizer) dateTime(buf *strings.Builder) Token {
	// Format variable formats to: "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"

	// xxx timezone
	if tokenizer.cur != ' ' || !unicode.IsUpper(tokenizer.peek) {
		str := buf.String()
		if str[len(str)-1] == 'Z' {
			buf.WriteString(" UTC")
		} else {
			panic("Expecting timezone")
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
	if err != nil {
		panic(err)
	}

	tokenizer.val = dateTime
	return DATETIME
}

func (tokenizer *Tokenizer) date(str string) Token {
	date, err := NewDateFromIso(str)
	if err != nil {
		panic(err)
	}
	tokenizer.val = date
	return DATE
}

func (tokenizer *Tokenizer) time(str string, addSeconds bool) Token {
	if addSeconds {
		str = str + ":00"
	}
	time, err := NewTimeFromIso(str)
	if err != nil {
		panic(err)
	}
	tokenizer.val = time
	return TIME
}

func (tokenizer *Tokenizer) number(str string, unitIndex int) Token {
	if unitIndex == 0 {
		number, err := strconv.ParseFloat(str, 64)
		if err != nil {
			panic(err)
		}
		tokenizer.val = NewNumber(number, "")
		return NUMBER
	} else {
		numberStr := str[0:unitIndex]
		unit := str[unitIndex:]
		number, err := strconv.ParseFloat(numberStr, 64)
		if err != nil {
			panic(err)
		}
		tokenizer.val = NewNumber(number, unit)
		return NUMBER
	}
}

func (tokenizer *Tokenizer) symbol() Token {
	c := tokenizer.cur
	tokenizer.consume()
	if c == ',' {
		return COMMA
	} else if c == '/' {
		return SLASH
	} else if c == ':' {
		return COLON
	} else if c == ';' {
		return SEMICOLON
	} else if c == '[' {
		return LBRACKET
	} else if c == ']' {
		return RBRACKET
	} else if c == '{' {
		return LBRACE
	} else if c == '}' {
		return RBRACE
	} else if c == '(' {
		return LPAREN
	} else if c == ')' {
		return RPAREN
	} else if c == '<' {
		if tokenizer.cur == '<' {
			tokenizer.consumeRune('<')
			return LT2
		} else if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return LTEQ
		} else {
			return LT
		}
	} else if c == '>' {
		if tokenizer.cur == '>' {
			tokenizer.consumeRune('>')
			return GT2
		} else if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return GTEQ
		} else {
			return GT
		}
	} else if c == '-' {
		if tokenizer.cur == '>' {
			tokenizer.consumeRune('>')
			return ARROW
		} else {
			return MINUS
		}
	} else if c == '=' {
		if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return EQ
		}
		return ASSIGN
	} else if c == '!' {
		if tokenizer.cur == '=' {
			tokenizer.consumeRune('=')
			return NOTEQ
		}
		return BANG
	} else if c == runeEOF {
		return EOF
	}
	panic("Unexpected symbol: '" + string(c) + "'")
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

func (tokenizer *Tokenizer) consume() {
	var err error
	tokenizer.cur = tokenizer.peek
	tokenizer.peek, _, err = tokenizer.in.ReadRune()
	if err != nil { // If end-of-stream, indicate with val of runeEOF
		tokenizer.peek = runeEOF
	}
}

func (tokenizer *Tokenizer) consumeRune(expected rune) {
	if tokenizer.cur != expected {
		panic("Expected " + string(expected))
	}
	tokenizer.consume()
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
