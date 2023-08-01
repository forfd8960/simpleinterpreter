package ast

import "github.com/forfd8960/simpleinterpreter/tokens"

type Node interface {
	TokenLiteral() string
}

type Stmt interface {
	Node
	StmtNode()
}

type Expression interface {
	Node
	ExprNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	}

	return ""
}

type Identifier struct {
	Token *tokens.Token
	Name  string
}

func (ident *Identifier) ExprNode() {}
func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}

type LetStmt struct {
	Token tokens.Token
	Ident *Identifier
	Value Expression
}

func (ls *LetStmt) StmtNode() {}
func (ls *LetStmt) TokenLiteral() string {
	return ls.Token.Literal
}
