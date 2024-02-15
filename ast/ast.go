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
	Type  string //storing the variable type
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var output bytes.Buffer

	output.WriteString("var ")
	output.WriteString(vs.Type) // Include the variable type
	output.WriteString(" ")
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

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (inl *IntegerLiteral) expressionNode()      {}
func (inl *IntegerLiteral) TokenLiteral() string { return inl.Token.Literal }
func (inl *IntegerLiteral) String() string       { return inl.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pre *PrefixExpression) expressionNode()      {}
func (pre *PrefixExpression) TokenLiteral() string { return pre.Token.Literal }
func (pre *PrefixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(pre.Operator)
	output.WriteString(pre.Right.String())
	output.WriteString(")")

	return output.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ine *InfixExpression) expressionNode()      {}
func (ine *InfixExpression) TokenLiteral() string { return ine.Token.Literal }
func (ine *InfixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(ine.Left.String())
	output.WriteString(" " + ine.Operator + " ")
	output.WriteString(ine.Right.String())
	output.WriteString(")")

	return output.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (strl *StringLiteral) expressionNode()      {}
func (strl *StringLiteral) TokenLiteral() string { return strl.Token.Literal }
func (strl *StringLiteral) String() string       { return strl.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (bln *Boolean) expressionNode()      {}
func (bln *Boolean) TokenLiteral() string { return bln.Token.Literal }
func (bln *Boolean) String() string       { return bln.Token.Literal }

type IfExpression struct {
	Token       token.Token // 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) expressionNode()      {}
func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }
func (ife *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ife.Condition.String())
	out.WriteString(" ")
	out.WriteString(ife.Consequence.String())

	if ife.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ife.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
