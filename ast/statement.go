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
