package haystack

type Token int

const (
	// Special tokens
	DEF Token = iota
	EOF

	// Literals
	literal_begin
	ID
	NUMBER
	STR
	REF
	URI
	DATE
	TIME
	DATETIME
	literal_end

	// Syntax
	syntax_begin
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
	syntax_end
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

func (tok Token) IsLiteral() bool {
	return literal_begin < tok && tok < literal_end
}

func (tok Token) IsSyntax() bool {
	return syntax_begin < tok && tok < syntax_end
}
