package ast

import (
	"Go-Tutorials/Core-lang/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (prg *Program) String() string {
	var out bytes.Buffer

	for _, str := range prg.Statements {
		out.WriteString(str.String())
	}

	return out.String()
}

type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var output bytes.Buffer

	output.WriteString(vs.TokenLiteral() + " ")
	output.WriteString(vs.Name.String())
	output.WriteString(" = ")

	if vs.Value != nil {
		output.WriteString(vs.Value.String())
	}

	output.WriteString(";")

	return output.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (idr *Identifier) expressionNode()      {}
func (idr *Identifier) TokenLiteral() string { return idr.Token.Literal }
func (idr *Identifier) String() string       { return idr.Value }

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var output bytes.Buffer

	output.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		output.WriteString(rs.ReturnValue.String())
	}

	output.WriteString(";")

	return output.String()
}

type ExpressionStatement struct {
	Token      token.Token // the 1st token of the expression
	Expression Expression
}

func (exs *ExpressionStatement) statementNode()       {}
func (exs *ExpressionStatement) TokenLiteral() string { return exs.Token.Literal }
func (exs *ExpressionStatement) String() string {
	if exs.Expression != nil {
		return exs.Expression.String()
	}
	return ""
}
