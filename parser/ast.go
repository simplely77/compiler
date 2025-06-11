package parser

type AST struct {
	Statements []Statement
}

type Statement interface {
	stmtNode()
}

type Assignment struct {
	Ident string
	Value Expr
}

func (a *Assignment) stmtNode() {}

type PrintStatement struct {
	Expr Expr
}

func (p *PrintStatement) stmtNode() {}

type InputStatement struct {
	Ident string
}

func (i *InputStatement) stmtNode() {}

type IfStatement struct {
	Condition Expr
	Then      []Statement
	Else      []Statement
}

func (i *IfStatement) stmtNode() {}

type WhileStatement struct {
	Condition Expr
	Body      []Statement
}

func (w *WhileStatement) stmtNode() {}

type Expr interface {
	exprNode()
}

type BinaryExpr struct {
	Op    string
	Left  Expr
	Right Expr
}

func (b *BinaryExpr) exprNode() {}

type IdentExpr struct {
	Name string
}

func (i *IdentExpr) exprNode() {}

type NumberExpr struct {
	Value string
}

func (n *NumberExpr) exprNode() {}

type BooleanExpr struct {
	Value bool
}

func (b *BooleanExpr) exprNode() {}

type ComparisonExpr struct {
	Op    string
	Left  Expr
	Right Expr
}

func (c *ComparisonExpr) exprNode() {}
