package ast

import (
	"github.com/forfd8960/simpleinterpreter/tokens"
)

type Stringer interface {
	String() string
}

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

func NewIdentifier1(name string) *Identifier {
	return &Identifier{
		Name:  name,
		Token: tokens.NewToken(tokens.IDENT, name, name),
	}
}

func (ident *Identifier) StmtNode() {}
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

type ClassStmt struct {
	NameIdent *Identifier
	Methods   map[string]*Function
}

func (ls *ClassStmt) StmtNode() {}
func (ls *ClassStmt) TokenLiteral() string {
	return tokens.CLASS.String()
}

func NewClassStmt(className *tokens.Token, methods []*Function) *ClassStmt {
	cls := &ClassStmt{
		NameIdent: &Identifier{
			Token: className,
			Name:  className.Literal,
		},
		Methods: make(map[string]*Function),
	}

	for _, mth := range methods {
		cls.Methods[mth.Name.Literal] = mth
	}

	return cls
}

type ThisExpr struct {
	keyword *tokens.Token
}

func NewThisExpr(kw *tokens.Token) *ThisExpr {
	return &ThisExpr{
		keyword: kw,
	}
}

func (te *ThisExpr) StmtNode() {}
func (te *ThisExpr) ExprNode() {}
func (te *ThisExpr) TokenLiteral() string {
	return te.keyword.Literal
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

type Binary struct {
	Left     Expression
	Operator *tokens.Token
	Right    Expression
}

func NewBinary(left, right Expression, operator *tokens.Token) *Binary {
	return &Binary{Left: left, Operator: operator, Right: right}
}

func (bin *Binary) ExprNode() {}
func (bin *Binary) TokenLiteral() string {
	return bin.Operator.Literal
}

type Unary struct {
	Operator *tokens.Token
	Right    Expression
}

func NewUnary(operator *tokens.Token, right Expression) *Unary {
	return &Unary{Operator: operator, Right: right}
}
func (un *Unary) ExprNode() {}
func (un *Unary) TokenLiteral() string {
	return un.Operator.Literal
}

type Call struct {
	Callee    Expression
	Arguments []Expression
}

func NewCall(callee Expression, arguments []Expression) *Call {
	return &Call{
		Callee:    callee,
		Arguments: arguments,
	}
}

func (call *Call) ExprNode() {}
func (call *Call) TokenLiteral() string {
	return call.Callee.TokenLiteral()
}

type Get struct {
	Expr Expression
	Name *tokens.Token
}

func NewGet(expr Expression, name *tokens.Token) *Get {
	return &Get{
		Expr: expr,
		Name: name,
	}
}

func (get *Get) ExprNode() {}
func (get *Get) TokenLiteral() string {
	return get.Expr.TokenLiteral()
}

// Set classIsntance.Property = value
type Set struct {
	Expr  Expression    // class instance
	Name  *tokens.Token // property name
	Value Expression    // property value
}

func NewSet(expr Expression, name *tokens.Token, value Expression) *Set {
	return &Set{
		Expr:  expr,
		Name:  name,
		Value: value,
	}
}

func (set *Set) ExprNode() {}
func (set *Set) TokenLiteral() string {
	return set.Expr.TokenLiteral() + "." + set.Name.Literal + " = " + set.Value.TokenLiteral()
}

type Literal struct {
	Value *tokens.Token
}

func NewLiteral(v *tokens.Token) *Literal {
	return &Literal{Value: v}
}

func (lter *Literal) ExprNode() {}
func (lter *Literal) TokenLiteral() string {
	return lter.Value.Literal
}

type Grouping struct {
	Expr Expression
}

func NewGrouping(expr Expression) *Grouping {
	return &Grouping{Expr: expr}
}

func (gp *Grouping) ExprNode() {}
func (gp *Grouping) TokenLiteral() string {
	return "grouping"
}
