package haystack

type Token struct {
	symbol  string
	literal bool
}

// Constructors

func tokenSyntax(symbol string) Token {
	return Token{
		symbol:  symbol,
		literal: false,
	}
}

func tokenLiteral(symbol string) Token {
	return Token{
		symbol:  symbol,
		literal: true,
	}
}

// Instances

// End of file
func tokenEof() Token {
	return tokenSyntax("eof")
}

// Literals

func tokenId() Token {
	return tokenSyntax("identifier")
}

func tokenNumber() Token {
	return tokenLiteral("Number")
}

func tokenStr() Token {
	return tokenLiteral("Str")
}

func tokenRef() Token {
	return tokenLiteral("Ref")
}

func tokenUri() Token {
	return tokenLiteral("Uri")
}

func tokenDate() Token {
	return tokenLiteral("Date")
}

func tokenTime() Token {
	return tokenLiteral("Time")
}

func tokenDateTime() Token {
	return tokenLiteral("DateTime")
}

// Syntax characters

func tokenColon() Token {
	return tokenSyntax(":")
}

func tokenComma() Token {
	return tokenSyntax(",")
}

func tokenSemiColon() Token {
	return tokenSyntax(",")
}

func tokenMinus() Token {
	return tokenSyntax("-")
}

func tokenEq() Token {
	return tokenSyntax("==")
}

func tokenNotEq() Token {
	return tokenSyntax("!=")
}

func tokenLt() Token {
	return tokenSyntax("<")
}

func tokenLt2() Token {
	return tokenSyntax("<<")
}

func tokenLtEq() Token {
	return tokenSyntax("<=")
}

func tokenGt() Token {
	return tokenSyntax(">")
}

func tokenGt2() Token {
	return tokenSyntax(">>")
}

func tokenGtEq() Token {
	return tokenSyntax(">=")
}

func tokenLbracket() Token {
	return tokenSyntax("[")
}

func tokenRbracket() Token {
	return tokenSyntax("]")
}

func tokenLbrace() Token {
	return tokenSyntax("{")
}

func tokenRbrace() Token {
	return tokenSyntax("}")
}

func tokenLparen() Token {
	return tokenSyntax("(")
}

func tokenRparen() Token {
	return tokenSyntax(")")
}

func tokenArrow() Token {
	return tokenSyntax("->")
}

func tokenSlash() Token {
	return tokenSyntax("/")
}

func tokenAssign() Token {
	return tokenSyntax("=")
}

func tokenBang() Token {
	return tokenSyntax("!")
}

func tokenNl() Token {
	return tokenSyntax("nl")
}

// Methods

func (token1 Token) equals(token2 Token) bool {
	return token1.symbol == token2.symbol
}
