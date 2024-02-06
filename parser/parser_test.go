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
	var int foobar = 123456;
	`
	// 	input := `
	// var x = 5;
	// var y = 10;
	// var foobar = 838383;
	// `

	lex := lexer.New(input)
	par := New(lex)

	program := par.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 5 {
		t.Fatalf("program.Statements does not contain 5 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		// {"b", "int"},
		// {"a", "string"},
		// {"foobar", "int"},
		{"x", "int"},
		{"y", "int"},
		{"foobar", "int"},
	}

	for i, tt := range tests {
		statement := program.Statements[i].(*ast.VarStatement) // Assuming direct cast for simplicity
		if !testVarStatement(t, statement, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, expectedName string, expectedType string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStatement, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStatement.Name.Value != expectedName {
		t.Errorf("varStatement.Name.Value not '%s', got=%s", expectedName, varStatement.Name.Value)
		return false
	}

	if varStatement.Name.TokenLiteral() != expectedName {
		t.Errorf("varStatement.Name.TokenLiteral() not '%s', got=%s", expectedName, varStatement.Name.TokenLiteral())
		return false
	}

	// Check for variable type
	if varStatement.Type != expectedType {
		t.Errorf("varStatement.Type not '%s', got=%s", expectedType, varStatement.Type)
		return false
	}

	return true
}
