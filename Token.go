package haystack

// Models a Zinc token. See https://project-haystack.org/doc/Zinc
type Token int

const (
	// Special tokens

	DEF Token = iota
	EOF

	// Literals
	literalBegin
	ID
	NUMBER
	STR
	REF
	URI
	DATE
	TIME
	DATETIME
	literalEnd

	// Syntax
	syntaxBegin
	COLON
	COMMA
	SEMICOLON
	MINUS
	EQ
	NOTEQ
	LT
	LT2
	LTEQ
	GT
	GT2
	GTEQ
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	LPAREN
	RPAREN
	ARROW
	SLASH
	ASSIGN
	BANG
	NL
	syntaxEnd
)

var tokens = [...]string{
	DEF: "",
	EOF: "eof",

	ID:       "id",
	NUMBER:   "Number",
	STR:      "Str",
	REF:      "Ref",
	URI:      "Uri",
	DATE:     "Date",
	TIME:     "Time",
	DATETIME: "DateTime",

	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",
	MINUS:     "-",
	EQ:        "=",
	NOTEQ:     "!=",
	LT:        "<",
	LT2:       "<<",
	LTEQ:      "<=",
	GT:        ">",
	GT2:       ">>",
	GTEQ:      ">=",
	LBRACKET:  "[",
	RBRACKET:  "]",
	LBRACE:    "{",
	RBRACE:    "}",
	LPAREN:    "(",
	RPAREN:    ")",
	ARROW:     "->",
	SLASH:     "/",
	ASSIGN:    "=",
	BANG:      "!",
	NL:        "nl",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	return s
}

// IsLiteral determines whether the token is an object literal
func (tok Token) IsLiteral() bool {
	return literalBegin < tok && tok < literalEnd
}

// IsSyntax determines whether the token is syntax punctuation
func (tok Token) IsSyntax() bool {
	return syntaxBegin < tok && tok < syntaxEnd
}
