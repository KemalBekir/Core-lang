package parser

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/lexer"
	"Go-Tutorials/Core-lang/token"
	"fmt"
	"strconv"
)

var precedences = map[token.TokenType]int{
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LESS_THEN:    LESSGREATER,
	token.GREATER_THEN: LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X OR !X
	CALL        // myFunction(x)
	INDEX       // array[index]
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
	par.registerPrefix(token.IDENT, par.parseIdentifier)
	par.registerPrefix(token.BANG, par.parsePrefixExpression)
	par.registerPrefix(token.MINUS, par.parsePrefixExpression)

	par.infixParseFunction = make(map[token.TokenType]infixParseFunction)
	par.registerInfix(token.PLUS, par.parseInfixExpression)
	par.registerInfix(token.MINUS, par.parseInfixExpression)
	par.registerInfix(token.SLASH, par.parseInfixExpression)
	par.registerInfix(token.EQ, par.parseInfixExpression)
	par.registerInfix(token.NOT_EQ, par.parseInfixExpression)
	par.registerInfix(token.LESS_THEN, par.parseInfixExpression)
	par.registerInfix(token.GREATER_THEN, par.parseInfixExpression)

	par.nextToken()
	par.nextToken()

	return par
}

func (par *Parser) nextToken() {
	par.currentToken = par.peekToken
	par.peekToken = par.lex.NextToken()
}

// func (p *Parser) ParseProgram() *ast.Program {
// 	program := &ast.Program{}
// 	program.Statements = []ast.Statement{}

// 	for p.currentToken.Type != token.END {
// 		stmt := p.parseStatement()
// 		if stmt != nil {
// 			fmt.Printf("Parsed statement: %T\n", stmt) // Debug print
// 			program.Statements = append(program.Statements, stmt)
// 		}
// 		p.nextToken()
// 	}

// 	return program
// }

// func (par *Parser) ParseProgram() *ast.Program {
// 	program := &ast.Program{}
// 	program.Statements = []ast.Statement{}

// 	for par.currentToken.Type != token.END {
// 		statement := par.parseStatement()
// 		if statement != nil {
// 			program.Statements = append(program.Statements, statement)
// 		}
// 		par.nextToken()
// 	}

// 	return program
// }

// Old ParseProgram
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

	if !par.ensureNext(token.VAR) {
		return nil
	}

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

	// Parse the expression for the value
	par.nextToken()
	statement.Value = par.parseExpression(LOWEST)

	if statement.Value == nil {
		return nil
	}

	// Ensure the statement ends with a semicolon
	if !par.ensureNext(token.SEMICOLON) {
		return nil
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

func (par *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFunction) {
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
		par.singnalPrefixParseFnNotFound(par.currentToken.Type)
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

func (par *Parser) singnalPrefixParseFnNotFound(tok token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s found", tok)
	par.errors = append(par.errors, message)
}

func (par *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    par.currentToken,
		Operator: par.currentToken.Literal,
	}

	par.nextToken()

	expression.Right = par.parseExpression(PREFIX)

	return expression
}

func (par *Parser) getUpcomingPrecedence() int {
	if par, ok := precedences[par.peekToken.Type]; ok {
		return par
	}

	return LOWEST
}

func (par *Parser) currentPrecedence() int {
	if par, ok := precedences[par.currentToken.Type]; ok {
		return par
	}

	return LOWEST
}

func (par *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    par.currentToken,
		Operator: par.currentToken.Literal,
		Left:     left,
	}

	precedence := par.currentPrecedence()
	par.nextToken()
	expression.Right = par.parseExpression(precedence)

	return expression
}
