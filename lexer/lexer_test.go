package lexer

import (
	"Go-Tutorials/Core-lang/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	var string b = "abc";
	var int a = 5;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.STRING_TYPE, "string"},
		{token.IDENT, "b"},
		{token.ASSIGN_OP, "="},
		{token.STRING, "abc"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.INT_TYPE, "int"},
		{token.IDENT, "a"},
		{token.ASSIGN_OP, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.END, ""},
	}

	lex := New(input)

	for i, tt := range tests {
		currentToken := lex.NextToken()

		if currentToken.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong, expected=%q, got=%q",
				i, tt.expectedType, currentToken.Type)
		}

		if currentToken.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong, expected=%q, got%q",
				i, tt.expectedLiteral, currentToken.Type)
		}
	}
}
