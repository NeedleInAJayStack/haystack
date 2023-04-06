package io

// Token categorizes the contents of a Zinc string into specific literals or syntax objects. See https://project-haystack.org/doc/Zinc
type Token int

const (
	// Special tokens

	// DEF is the default value token
	DEF Token = iota
	// EOF indicates that the end of the input has been reached
	EOF

	// Literals

	literalBegin
	// ID indicates non-string literal text, typically variable names, etc.
	ID
	// NUMBER indicates a Number literal
	NUMBER
	// STR indicates a Str literal
	STR
	// SYMBOL indicates a Symbol literal
	SYMBOL
	// REF indicates a Ref literal
	REF
	// URI indicates a Uri literal
	URI
	// DATE indicates a Date literal
	DATE
	// TIME indicates a Time literal
	TIME
	// DATETIME indicates a DateTime literal
	DATETIME
	literalEnd

	// Syntax

	syntaxBegin
	// COLON indicates ':'
	COLON
	// COMMA indicates ','
	COMMA
	// SEMICOLON indicates ';'
	SEMICOLON
	// MINUS indicates '-'
	MINUS
	// EQ indicates '=='
	EQ
	// NOTEQ indicates '!='
	NOTEQ
	// LT indicates '<'
	LT
	// LT2 indicates '<<'
	LT2
	// LTEQ indicates '<='
	LTEQ
	// GT indicates '>'
	GT
	// GT2 inciates '>>'
	GT2
	// GTEQ indicates '>='
	GTEQ
	// LBRACKET indicates '['
	LBRACKET
	// RBRACKET indicates ']'
	RBRACKET
	// LBRACE indicates '{'
	LBRACE
	// RBRACE indicates '}'
	RBRACE
	// LPAREN indicates '('
	LPAREN
	// RPAREN indicates ')'
	RPAREN
	// ARROW indicates '->'
	ARROW
	// SLASH indicates '/'
	SLASH
	// ASSIGN indicates '='
	ASSIGN
	// BANG indicates '!'
	BANG
	// NL indicates the newline character
	NL
	syntaxEnd
)

var tokens = [...]string{
	DEF: "",
	EOF: "eof",

	ID:       "id",
	NUMBER:   "Number",
	STR:      "Str",
	SYMBOL:   "Symbol",
	REF:      "Ref",
	URI:      "Uri",
	DATE:     "Date",
	TIME:     "Time",
	DATETIME: "DateTime",

	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",
	MINUS:     "-",
	EQ:        "==",
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
