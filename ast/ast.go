package ast

import "Go-Tutorials/Core-lang/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (pro *Program) TokenLiteral() string {
	if len(pro.Statements) > 0 {
		return pro.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (idr *Identifier) expressionNode()      {}
func (idr *Identifier) TokenLiteral() string { return idr.Token.Literal }
