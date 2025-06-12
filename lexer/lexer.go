package lexer

import "fmt"

type TokenType int

func (tt TokenType) String() string {
	return tokenTypeToString[tt]
}

var tokenTypeToString = map[TokenType]string{
	TOKEN_IDENT:        "IDENT",
	TOKEN_NUMBER:       "NUMBER",
	TOKEN_PLUS:         "PLUS",
	TOKEN_MINUS:        "MINUS",
	TOKEN_MULTIPLY:     "MULTIPLY",
	TOKEN_DIVIDE:       "DIVIDE",
	TOKEN_ASSIGN:       "ASSIGN",
	TOKEN_SEMICOLON:    "SEMICOLON",
	TOKEN_LPAREN:       "LPAREN",
	TOKEN_RPAREN:       "RPAREN",
	TOKEN_LBRACE:       "LBRACE",
	TOKEN_RBRACE:       "RBRACE",
	TOKEN_LESS:         "LESS",
	TOKEN_GREATER:      "GREATER",
	TOKEN_EQUAL:        "EQUAL",
	TOKEN_NOT_EQUAL:    "NOT_EQUAL",
	TOKEN_LESS_EQUAL:   "LESS_EQUAL",
	TOKEN_GREATER_EQUAL: "GREATER_EQUAL",
	TOKEN_KEYWORD:      "KEYWORD",
	TOKEN_EOF:          "EOF",
	TOKEN_ILLEGAL:      "ILLEGAL",
}

const (
	TOKEN_IDENT TokenType = iota
	TOKEN_NUMBER
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULTIPLY
	TOKEN_DIVIDE
	TOKEN_ASSIGN
	TOKEN_SEMICOLON
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_LESS
	TOKEN_GREATER
	TOKEN_EQUAL
	TOKEN_NOT_EQUAL
	TOKEN_LESS_EQUAL
	TOKEN_GREATER_EQUAL
	TOKEN_KEYWORD
	TOKEN_EOF
	TOKEN_ILLEGAL
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type LexerError struct {
	Line    int
	Column  int
	Message string
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("词法错误：第%d行第%d列: %s", e.Line, e.Column, e.Message)
}

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      rune
	line    int
	column  int
	errors  []*LexerError
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
		errors: make([]*LexerError, 0),
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.readPos])
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
}

func (l *Lexer) NextToken() Token {
	var tok Token
	tok.Line = l.line
	tok.Column = l.column

	l.skipWhitespace()

	switch l.ch {
	case '+':
		tok = newToken(TOKEN_PLUS, l.ch)
	case '-':
		if isDigit(l.peekChar()) {
			l.readChar()
			tok.Type = TOKEN_NUMBER
			tok.Literal = "-" + l.readNumber()
			return tok
		}
		tok = newToken(TOKEN_MINUS, l.ch)
	case '*':
		tok = newToken(TOKEN_MULTIPLY, l.ch)
	case '/':
		tok = newToken(TOKEN_DIVIDE, l.ch)
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_EQUAL, Literal: "=="}
		} else {
			tok = newToken(TOKEN_ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_NOT_EQUAL, Literal: "!="}
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_LESS_EQUAL, Literal: "<="}
		} else {
			tok = newToken(TOKEN_LESS, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_GREATER_EQUAL, Literal: ">="}
		} else {
			tok = newToken(TOKEN_GREATER, l.ch)
		}
	case ';':
		tok = newToken(TOKEN_SEMICOLON, l.ch)
	case '(':
		tok = newToken(TOKEN_LPAREN, l.ch)
	case ')':
		tok = newToken(TOKEN_RPAREN, l.ch)
	case '{':
		tok = newToken(TOKEN_LBRACE, l.ch)
	case '}':
		tok = newToken(TOKEN_RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TOKEN_NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			l.errors = append(l.errors, &LexerError{
				Line:    l.line,
				Column:  l.column,
				Message: fmt.Sprintf("非法字符: %c", l.ch),
			})
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
	// 跳过注释
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.readChar() // 跳过换行符
		l.skipWhitespace() // 继续跳过空白字符
	}
}

func (l *Lexer) HasErrors() bool {
	return len(l.errors) > 0
}

func (l *Lexer) GetErrors() []*LexerError {
	return l.errors
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func newToken(tokenType TokenType, ch rune) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPos])
}
