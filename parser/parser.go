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
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X OR !X
	CALL        // myFunction(x)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:               EQUALS,
	token.NOT_EQ:           EQUALS,
	token.LESS_THEN:        LESSGREATER,
	token.GREATER_THEN:     LESSGREATER,
	token.PLUS:             SUM,
	token.MINUS:            SUM,
	token.SLASH:            PRODUCT,
	token.ASTERISK:         PRODUCT,
	token.LEFT_PARANTHESIS: CALL,
	token.LEFT_BRACKET:     INDEX,
}

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
	par.registerPrefix(token.INT, par.parseIntegerLiteral)
	par.registerPrefix(token.BANG, par.parsePrefixExpression)
	par.registerPrefix(token.MINUS, par.parsePrefixExpression)
	par.registerPrefix(token.TRUE, par.parseBoolean)
	par.registerPrefix(token.FALSE, par.parseBoolean)
	par.registerPrefix(token.LEFT_PARANTHESIS, par.parseParenthesizedExpression)
	par.registerPrefix(token.IF, par.parseIfExpression)
	par.registerPrefix(token.FUNCTION, par.parseFunctionLiteral)
	par.registerPrefix(token.STRING, par.parseStringLiteral)

	par.infixParseFunction = make(map[token.TokenType]infixParseFunction)
	par.registerInfix(token.PLUS, par.parseInfixExpression)
	par.registerInfix(token.MINUS, par.parseInfixExpression)
	par.registerInfix(token.SLASH, par.parseInfixExpression)
	par.registerInfix(token.ASTERISK, par.parseInfixExpression)
	par.registerInfix(token.EQ, par.parseInfixExpression)
	par.registerInfix(token.NOT_EQ, par.parseInfixExpression)
	par.registerInfix(token.LESS_THEN, par.parseInfixExpression)
	par.registerInfix(token.GREATER_THEN, par.parseInfixExpression)
	par.registerInfix(token.SLASH, par.parseInfixExpression)

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

	if !par.expectNextType() {
		return nil
	}
	declaredType := par.currentToken.Literal

	if !par.ensureNext(token.IDENT) {
		return nil
	}
	statement.Name = &ast.Identifier{Token: par.currentToken, Value: par.currentToken.Literal}

	if !par.ensureNext(token.ASSIGN_OP) {
		return nil
	}

	par.nextToken()
	statement.Value = par.parseExpression(LOWEST)

	if statement.Value == nil {
		return nil
	}

	valueType := par.resolveExpressionType(statement.Value)
	if declaredType != valueType {
		errMsg := fmt.Sprintf("Type mismatch: cannot assign %s to %s variable", valueType, declaredType)
		par.errors = append(par.errors, errMsg)
	}

	if !par.ensureNext(token.SEMICOLON) {
		return nil
	}

	statement.Type = declaredType
	return statement
}

func (par *Parser) resolveExpressionType(expr ast.Expression) string {
	switch expr.(type) {
	case *ast.IntegerLiteral:
		return "int"
	case *ast.StringLiteral:
		return "string"
	default:
		return "uknown"
	}
}

func (par *Parser) currentTokenIs(tok token.TokenType) bool {
	return par.currentToken.Type == tok
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
	// fmt.Println("parseReturnStatement: Entering")
	statement := &ast.ReturnStatement{Token: par.currentToken}

	// fmt.Printf("parseReturnStatement: Current token: %+v\n", par.currentToken)

	par.nextToken()

	// fmt.Printf("parseReturnStatement: Token after nextToken: %+v\n", par.currentToken)

	statement.ReturnValue = par.parseExpression(LOWEST)

	// fmt.Printf("parseReturnStatement: Parsed return value: %+v\n", statement.ReturnValue)

	if par.peekedTokenIs(token.SEMICOLON) {
		// fmt.Println("parseReturnStatement: Found semicolon, consuming it.")
		par.nextToken()
	} else {
		// fmt.Println("parseReturnStatement: No semicolon found after return value.")
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

	for !par.peekedTokenIs(token.SEMICOLON) && precedence < par.getUpcomingPrecedence() {
		infix := par.infixParseFunction[par.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		par.nextToken()

		leftExpression = infix(leftExpression)
	}

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

func (par *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: par.currentToken, Value: par.currentToken.Literal}
}

func (par *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: par.currentToken, Value: par.currentTokenIs(token.TRUE)}
}

func (par *Parser) parseParenthesizedExpression() ast.Expression {
	par.nextToken()

	expression := par.parseExpression(LOWEST)

	if !par.ensureNext(token.RIGHT_PARANTHESIS) {
		return nil
	}

	return expression
}

func (par *Parser) parseIfExpression() ast.Expression {
	expression := *&ast.IfExpression{Token: par.currentToken}

	if !par.ensureNext(token.LEFT_PARANTHESIS) {
		return nil
	}

	par.nextToken()
	expression.Condition = par.parseExpression(LOWEST)

	if !par.ensureNext(token.RIGHT_PARANTHESIS) {
		return nil
	}

	if !par.ensureNext(token.LEFT_CURLY_BRACE) {
		return nil
	}

	expression.Consequence = par.parseBlockStatement()

	if par.peekedTokenIs(token.ELSE) {
		par.nextToken()

		if !par.ensureNext(token.LEFT_CURLY_BRACE) {
			return nil
		}

		expression.Alternative = par.parseBlockStatement()
	}

	return &expression
}

func (par *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: par.currentToken}
	block.Statements = []ast.Statement{}

	par.nextToken()

	for !par.currentTokenIs(token.RIGHT_CURLY_BRACE) && !par.currentTokenIs(token.END) {
		statement := par.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		par.nextToken()
	}

	return block
}

func (par *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: par.currentToken}

	if !par.ensureNext(token.LEFT_PARANTHESIS) {
		return nil
	}

	literal.Parameters = par.parseFunctionParameters()

	if !par.ensureNext(token.LEFT_CURLY_BRACE) {
		return nil
	}

	literal.Body = par.parseBlockStatement()

	return literal
}

func (par *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if par.peekedTokenIs(token.RIGHT_PARANTHESIS) {
		par.nextToken()
		return identifiers
	}

	par.nextToken()

	ident := &ast.Identifier{Token: par.currentToken, Value: par.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for par.peekedTokenIs(token.COMMA) {
		par.nextToken()
		par.nextToken()
		ident := &ast.Identifier{Token: par.currentToken, Value: par.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !par.ensureNext(token.RIGHT_PARANTHESIS) {
		return nil
	}

	return identifiers
}
