package lexer

var keywords = map[string]TokenType{
	"if":    TOKEN_KEYWORD,
	"else":  TOKEN_KEYWORD,
	"while": TOKEN_KEYWORD,
	"print": TOKEN_KEYWORD,
	"input": TOKEN_KEYWORD,
	"true":  TOKEN_KEYWORD,
	"false": TOKEN_KEYWORD,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}
