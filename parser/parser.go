package parser

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/lexer"
	"Go-Tutorials/Core-lang/token"
)

type Parser struct {
	lex *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
}

func New(lex *lexer.Lexer) *Parser {
	par := &Parser{lex: lex}

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
		return par.praseVarStatement()
	default:
		return nil
	}
}

func (par *Parser) praseVarStatement() *ast.VarStatement {
	statement := &ast.VarStatement{Token: par.currentToken}

	if !par.ensureNext(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: par.currentToken, Value: par.currentToken.Literal}

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
		return false
	}
}
