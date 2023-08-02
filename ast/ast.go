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

func NewIdentifier(name *tokens.Token) *Identifier {
	return &Identifier{Token: name, Name: name.Literal}
}

func (ident *Identifier) ExprNode() {}
func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}

type LetStmt struct {
	Ident    *Identifier
	InitExpr Expression
}

func (ls *LetStmt) StmtNode() {}
func (ls *LetStmt) TokenLiteral() string {
	return tokens.LET.String()
}

func NewLetStmt(ident *tokens.Token, expr Expression) *LetStmt {
	return &LetStmt{
		Ident: &Identifier{
			Token: ident,
			Name:  ident.Literal,
		},
		InitExpr: expr,
	}
}

type Assign struct {
	Name  *tokens.Token
	Value Expression
}

func (as *Assign) ExprNode() {}
func (as *Assign) TokenLiteral() string {
	return as.Name.Literal
}

func NewAssign(name *tokens.Token, value Expression) *Assign {
	return &Assign{
		Name:  name,
		Value: value,
	}
}

type Logical struct {
	Left     Expression
	Right    Expression
	Operator *tokens.Token
}

func NewLogical(left, right Expression, op *tokens.Token) *Logical {
	return &Logical{Left: left, Right: right, Operator: op}
}
func (lg *Logical) ExprNode() {}
func (lg *Logical) TokenLiteral() string {
	return lg.Operator.Literal
}
