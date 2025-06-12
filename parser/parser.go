package parser

import (
	"compiler/lexer"
	"fmt"
)

type Parser struct {
	lex       *lexer.Lexer
	lookahead lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lex: l}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.lookahead = p.lex.NextToken()
}

func (p *Parser) Parse() (*AST, error) {
	ast := &AST{}
	for p.lookahead.Type != lexer.TOKEN_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		ast.Statements = append(ast.Statements, stmt)
	}
	// 语义分析：检查所有变量使用是否已定义
	if err := semanticCheck(ast); err != nil {
		return nil, err
	}
	return ast, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	switch p.lookahead.Type {
	case lexer.TOKEN_IDENT:
		return p.parseAssignment()
	case lexer.TOKEN_KEYWORD:
		switch p.lookahead.Literal {
		case "print":
			return p.parsePrint()
		case "input":
			return p.parseInput()
		case "if":
			return p.parseIf()
		case "while":
			return p.parseWhile()
		}
	}
	return nil, p.newError("未知语句")
}

func (p *Parser) parseAssignment() (Statement, error) {
	ident := p.lookahead.Literal
	p.nextToken()
	if p.lookahead.Type != lexer.TOKEN_ASSIGN {
		return nil, p.newError("赋值语句缺少 '='")
	}
	p.nextToken()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("DEBUG: In parseAssignment, before semicolon check. Lookahead: Type=%s, Literal=\"%s\", Line=%d, Column=%d\n", p.lookahead.Type, p.lookahead.Literal, p.lookahead.Line, p.lookahead.Column)
	if p.lookahead.Type != lexer.TOKEN_SEMICOLON {
		return nil, p.newError("赋值语句缺少分号")
	}
	p.nextToken()
	return &Assignment{Ident: ident, Value: expr}, nil
}

func (p *Parser) parsePrint() (Statement, error) {
	p.nextToken()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("DEBUG: In parsePrint, before semicolon check. Lookahead: Type=%s, Literal=\"%s\", Line=%d, Column=%d\n", p.lookahead.Type, p.lookahead.Literal, p.lookahead.Line, p.lookahead.Column)
	if p.lookahead.Type != lexer.TOKEN_SEMICOLON {
		return nil, p.newError("print语句缺少分号")
	}
	p.nextToken()
	return &PrintStatement{Expr: expr}, nil
}

func (p *Parser) parseInput() (Statement, error) {
	p.nextToken()
	if p.lookahead.Type != lexer.TOKEN_IDENT {
		return nil, p.newError("input语句需要变量名")
	}
	ident := p.lookahead.Literal
	p.nextToken()
	if p.lookahead.Type != lexer.TOKEN_SEMICOLON {
		return nil, p.newError("input语句缺少分号")
	}
	p.nextToken()
	return &InputStatement{Ident: ident}, nil
}

func (p *Parser) parseIf() (Statement, error) {
	p.nextToken()
	if p.lookahead.Type != lexer.TOKEN_LPAREN {
		return nil, p.newError("if语句缺少左括号")
	}
	p.nextToken()
	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.lookahead.Type != lexer.TOKEN_RPAREN {
		return nil, p.newError("if语句缺少右括号")
	}
	p.nextToken()
	
	if p.lookahead.Type != lexer.TOKEN_LBRACE {
		return nil, p.newError("if语句缺少左花括号")
	}
	p.nextToken()
	
	thenStmts := []Statement{}
	for p.lookahead.Type != lexer.TOKEN_RBRACE {
		if p.lookahead.Type == lexer.TOKEN_EOF {
			return nil, p.newError("if语句缺少右花括号")
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		thenStmts = append(thenStmts, stmt)
	}
	p.nextToken()
	
	elseStmts := []Statement{}
	if p.lookahead.Type == lexer.TOKEN_KEYWORD && p.lookahead.Literal == "else" {
		p.nextToken()
		if p.lookahead.Type != lexer.TOKEN_LBRACE {
			return nil, p.newError("else语句缺少左花括号")
		}
		p.nextToken()
		
		for p.lookahead.Type != lexer.TOKEN_RBRACE {
			if p.lookahead.Type == lexer.TOKEN_EOF {
				return nil, p.newError("else语句缺少右花括号")
			}
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			elseStmts = append(elseStmts, stmt)
		}
		p.nextToken()
	}
	
	return &IfStatement{
		Condition: condition,
		Then:      thenStmts,
		Else:      elseStmts,
	}, nil
}

func (p *Parser) parseWhile() (Statement, error) {
	p.nextToken()
	if p.lookahead.Type != lexer.TOKEN_LPAREN {
		return nil, p.newError("while语句缺少左括号")
	}
	p.nextToken()
	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.lookahead.Type != lexer.TOKEN_RPAREN {
		return nil, p.newError("while语句缺少右括号")
	}
	p.nextToken()
	
	if p.lookahead.Type != lexer.TOKEN_LBRACE {
		return nil, p.newError("while语句缺少左花括号")
	}
	p.nextToken()
	
	body := []Statement{}
	for p.lookahead.Type != lexer.TOKEN_RBRACE {
		if p.lookahead.Type == lexer.TOKEN_EOF {
			return nil, p.newError("while语句缺少右花括号")
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}
	p.nextToken()
	
	return &WhileStatement{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) parseExpr() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	
	for p.lookahead.Type == lexer.TOKEN_PLUS || p.lookahead.Type == lexer.TOKEN_MINUS {
		op := p.lookahead.Literal
		p.nextToken()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Op: op, Left: left, Right: right}
	}
	
	// 检查比较运算符
	if p.lookahead.Type == lexer.TOKEN_LESS || p.lookahead.Type == lexer.TOKEN_GREATER ||
	   p.lookahead.Type == lexer.TOKEN_EQUAL || p.lookahead.Type == lexer.TOKEN_NOT_EQUAL ||
	   p.lookahead.Type == lexer.TOKEN_LESS_EQUAL || p.lookahead.Type == lexer.TOKEN_GREATER_EQUAL {
		op := p.lookahead.Literal
		p.nextToken()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &ComparisonExpr{Op: op, Left: left, Right: right}, nil
	}
	
	return left, nil
}

func (p *Parser) parseTerm() (Expr, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for p.lookahead.Type == lexer.TOKEN_MULTIPLY || p.lookahead.Type == lexer.TOKEN_DIVIDE {
		op := p.lookahead.Literal
		p.nextToken()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseFactor() (Expr, error) {
	switch p.lookahead.Type {
	case lexer.TOKEN_IDENT:
		ident := p.lookahead.Literal
		p.nextToken()
		return &IdentExpr{Name: ident}, nil
	case lexer.TOKEN_NUMBER:
		val := p.lookahead.Literal
		p.nextToken()
		return &NumberExpr{Value: val}, nil
	case lexer.TOKEN_KEYWORD:
		if p.lookahead.Literal == "true" || p.lookahead.Literal == "false" {
			val := p.lookahead.Literal == "true"
			p.nextToken()
			return &BooleanExpr{Value: val}, nil
		}
		return nil, p.newError("非法的关键字")
	case lexer.TOKEN_LPAREN:
		p.nextToken()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.lookahead.Type != lexer.TOKEN_RPAREN {
			return nil, p.newError("缺少右括号")
		}
		p.nextToken()
		return expr, nil
	default:
		return nil, p.newError("非法表达式")
	}
}

func (p *Parser) newError(msg string) error {
	return fmt.Errorf("语法错误：第%d行第%d列: %s", p.lookahead.Line, p.lookahead.Column, msg)
}

// 语义分析：检查所有变量使用是否已定义
func semanticCheck(ast *AST) error {
	defined := make(map[string]bool)
	for _, stmt := range ast.Statements {
		switch s := stmt.(type) {
		case *Assignment:
			// 先检查右侧表达式
			if err := checkExprDefined(s.Value, defined); err != nil {
				return err
			}
			// 再标记左侧变量为已定义
			defined[s.Ident] = true
		case *PrintStatement:
			if err := checkExprDefined(s.Expr, defined); err != nil {
				return err
			}
		case *InputStatement:
			defined[s.Ident] = true
		case *IfStatement:
			if err := checkExprDefined(s.Condition, defined); err != nil {
				return err
			}
			for _, stmt := range s.Then {
				if err := checkStatementDefined(stmt, defined); err != nil {
					return err
				}
			}
			for _, stmt := range s.Else {
				if err := checkStatementDefined(stmt, defined); err != nil {
					return err
				}
			}
		case *WhileStatement:
			if err := checkExprDefined(s.Condition, defined); err != nil {
				return err
			}
			for _, stmt := range s.Body {
				if err := checkStatementDefined(stmt, defined); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func checkExprDefined(expr Expr, defined map[string]bool) error {
	switch e := expr.(type) {
	case *IdentExpr:
		if !defined[e.Name] {
			return fmt.Errorf("语义错误：变量 '%s' 未定义", e.Name)
		}
	case *BinaryExpr:
		if err := checkExprDefined(e.Left, defined); err != nil {
			return err
		}
		if err := checkExprDefined(e.Right, defined); err != nil {
			return err
		}
	case *ComparisonExpr:
		if err := checkExprDefined(e.Left, defined); err != nil {
			return err
		}
		if err := checkExprDefined(e.Right, defined); err != nil {
			return err
		}
	}
	return nil
}

func checkStatementDefined(stmt Statement, defined map[string]bool) error {
	switch s := stmt.(type) {
	case *Assignment:
		if err := checkExprDefined(s.Value, defined); err != nil {
			return err
		}
		defined[s.Ident] = true
	case *PrintStatement:
		if err := checkExprDefined(s.Expr, defined); err != nil {
			return err
		}
	case *InputStatement:
		defined[s.Ident] = true
	}
	return nil
}
