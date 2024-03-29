package ast

import (
	"Go-Tutorials/Core-lang/token"
	"bytes"
	"strings"
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
	var output bytes.Buffer

	output.WriteString("if")
	output.WriteString(ife.Condition.String())
	output.WriteString(" ")
	output.WriteString(ife.Consequence.String())

	if ife.Alternative != nil {
		output.WriteString("else ")
		output.WriteString(ife.Alternative.String())
	}

	return output.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var output bytes.Buffer

	for _, s := range bs.Statements {
		output.WriteString(s.String())
	}

	return output.String()
}

type FunctionLiteral struct {
	Token      token.Token // 'function' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fnl *FunctionLiteral) expressionNode()      {}
func (fnl *FunctionLiteral) TokenLiteral() string { return fnl.Token.Literal }
func (fnl *FunctionLiteral) String() string {
	var output bytes.Buffer

	parameters := []string{}
	for _, p := range fnl.Parameters {
		parameters = append(parameters, p.String())
	}

	output.WriteString(fnl.TokenLiteral())
	output.WriteString("(")
	output.WriteString(strings.Join(parameters, ", "))
	output.WriteString(") ")
	output.WriteString(fnl.Body.String())

	return output.String()
}

type CallExpression struct {
	Token     token.Token // '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var output bytes.Buffer

	arguments := []string{}
	for _, a := range ce.Arguments {
		arguments = append(arguments, a.String())
	}

	output.WriteString(ce.Function.String())
	output.WriteString("(")
	output.WriteString(strings.Join(arguments, ", "))
	output.WriteString(")")

	return output.String()
}

type ArrayLiteral struct {
	Token    token.Token // '[' token
	Elements []Expression
}

func (arl *ArrayLiteral) expressionNode()      {}
func (arl *ArrayLiteral) TokenLiteral() string { return arl.Token.Literal }
func (arl *ArrayLiteral) String() string {
	var output bytes.Buffer

	elements := []string{}

	for _, el := range arl.Elements {
		elements = append(elements, el.String())
	}

	output.WriteString("[")
	output.WriteString(strings.Join(elements, ", "))
	output.WriteString("]")

	return output.String()
}

type IndexExpression struct {
	Token token.Token // [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(ie.Left.String())
	output.WriteString("[")
	output.WriteString(ie.Index.String())
	output.WriteString("])")

	return output.String()
}

type HashLiteral struct {
	Token token.Token // '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var output bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	output.WriteString("{")
	output.WriteString(strings.Join(pairs, ", "))
	output.WriteString("}")

	return output.String()
}
