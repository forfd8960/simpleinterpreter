package ast

import "github.com/forfd8960/simpleinterpreter/tokens"

type ReturnStmt struct {
	Keyword *tokens.Token
	Value   Expression
}

func NewReturnStmt(kw *tokens.Token, value Expression) *ReturnStmt {
	return &ReturnStmt{
		Keyword: kw,
		Value:   value,
	}
}

func (rt *ReturnStmt) StmtNode() {}
func (rt *ReturnStmt) TokenLiteral() string {
	return "return"
}

type Block struct {
	Statements []Stmt
}

func (blk *Block) StmtNode() {}
func (blk *Block) TokenLiteral() string {
	return "block"
}

func NewBlockStmt(stmts []Stmt) *Block {
	return &Block{Statements: stmts}
}

type ExpressionStmt struct {
	Expr Expression
}

func NewExpressionStmt(e Expression) *ExpressionStmt {
	return &ExpressionStmt{Expr: e}
}

func (espst *ExpressionStmt) StmtNode() {}
func (espst *ExpressionStmt) TokenLiteral() string {
	return "expression_block"
}

type Function struct {
	Name       *tokens.Token
	Parameters []*tokens.Token
	Body       []Stmt
}

func (fn *Function) StmtNode() {}
func (fn *Function) TokenLiteral() string {
	return "function"
}

func NewFunctionStmt(name *tokens.Token, params []*tokens.Token, body []Stmt) *Function {
	return &Function{Name: name, Parameters: params, Body: body}
}
