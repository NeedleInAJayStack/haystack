package haystack

type Token struct {
	symbol  string
	literal bool
}

func refToken() Token {
	return Token{
		symbol:  "Ref",
		literal: true,
	}
}
