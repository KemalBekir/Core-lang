package parser

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
var int b = 5;
var string a = "abc";
var int fizz = 123456;
`

	lex := lexer.New(input)
	par := New(lex)

	program := par.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"b"},
		{"y"},
		{"fizz"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testVarStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStatment, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStatment.Name.Value != name {
		t.Errorf("varStatement.Name.Value not '%s', got=%s",
			name, varStatment.Name.TokenLiteral())
		return false
	}

	if varStatment.Name.TokenLiteral() != name {
		t.Errorf("varStatement.Name.TokenLiteral() not '%s', got=%s",
			name, varStatment.Name.TokenLiteral())
		return false
	}

	return true
}
