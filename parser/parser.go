package parser

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/lexer"
	"Go-Tutorials/Core-lang/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	SUM // +
)

type Parser struct {
	lex    *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFunction map[token.TokenType]prefixParseFunction
	infixParseFunction  map[token.TokenType]infixParseFunction
}

type (
	prefixParseFunction func() ast.Expression
	infixParseFunction  func(ast.Expression) ast.Expression
)

func New(lex *lexer.Lexer) *Parser {
	par := &Parser{
		lex:    lex,
		errors: []string{},
	}

	par.prefixParseFunction = make(map[token.TokenType]prefixParseFunction)
	par.regiterPrefix(token.IDENT, par.parseIdentifier)

	par.nextToken()
	par.nextToken()

	return par
}

func (par *Parser) nextToken() {
	par.currentToken = par.peekToken
	par.peekToken = par.lex.NextToken()
}

func (par *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for par.currentToken.Type != token.END {
		statement := par.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		par.nextToken()
	}

	return program
}

func (par *Parser) parseStatement() ast.Statement {
	switch par.currentToken.Type {
	case token.VAR:
		return par.parseVarStatement()
	case token.RETURN:
		return par.parseReturnStatement()
	default:
		return par.parseExpressionStatement()
	}
}

func (par *Parser) parseVarStatement() *ast.VarStatement {
	statement := &ast.VarStatement{Token: par.currentToken}

	if !par.ensureNext(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: par.currentToken, Value: par.currentToken.Literal}

	if !par.expectNextType() {
		return nil
	}

	statement.Type = par.currentToken.Literal

	if !par.ensureNext(token.ASSIGN_OP) {
		return nil
	}

	for !par.currentTokenIs(token.SEMICOLON) {
		par.nextToken()
	}

	return statement
}

func (par *Parser) currentTokenIs(tok token.TokenType) bool {
	return par.peekToken.Type == tok
}

func (par *Parser) peekedTokenIs(tok token.TokenType) bool {
	return par.peekToken.Type == tok
}
func (par *Parser) ensureNext(tok token.TokenType) bool {
	if par.peekedTokenIs(tok) {
		par.nextToken()
		return true
	} else {
		par.peekUnexpectedError(tok)
		return false
	}
}

func (par *Parser) Errors() []string {
	return par.errors
}

func (par *Parser) peekUnexpectedError(tok token.TokenType) {
	msg := fmt.Sprintf("expected next token to be - %s, got - %s instead",
		tok, par.peekToken.Type)
	par.errors = append(par.errors, msg)
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: par.currentToken}

	par.nextToken()

	for !par.currentTokenIs(token.SEMICOLON) {
		par.nextToken()
	}

	return statement
}

func (par *Parser) regiterPrefix(tokenType token.TokenType, fn prefixParseFunction) {
	par.prefixParseFunction[tokenType] = fn
}

func (par *Parser) registerInfix(tokenType token.TokenType, fn infixParseFunction) {
	par.infixParseFunction[tokenType] = fn
}

func (par *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: par.currentToken}

	statement.Expression = par.parseExpression(LOWEST)

	if par.peekedTokenIs(token.SEMICOLON) {
		par.nextToken()
	}

	return statement
}

func (par *Parser) parseExpression(precedence int) ast.Expression {
	prefix := par.prefixParseFunction[par.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	return leftExpression
}

func (par *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: par.currentToken, Value: par.currentToken.Literal}
}

func (par *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: par.currentToken}

	value, err := strconv.ParseInt(par.currentToken.Literal, 0, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", par.currentToken.Literal)
		par.errors = append(par.errors, message)
		return nil
	}

	literal.Value = value

	return literal
}

func (par *Parser) expectNextType() bool {
	types := map[token.TokenType]bool{
		token.INT_TYPE:    true,
		token.STRING_TYPE: true,
	}
	if types[par.peekToken.Type] {
		par.nextToken()
		return true
	} else {
		par.peekUnexpectedError(token.INT_TYPE)
		return false
	}
}
