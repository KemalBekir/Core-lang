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

	lex := lexer.New(input)
	// for tok := lex.NextToken(); tok.Type != token.END; tok = lex.NextToken() {
	// 	fmt.Printf("%+v\n", tok)
	// }

	par := New(lex)

	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"b", "int"},
		{"a", "string"},
		{"foobar", "int"},
	}

	for i, tt := range tests {
		statement := program.Statements[i].(*ast.VarStatement)
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

	if varStatement.Type != expectedType {
		t.Errorf("varStatement.Type not '%s', got=%s", expectedType, varStatement.Type)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, par *Parser) {
	errors := par.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
return 3;
return 15;
return 123321;
	`

	lex := lexer.New(input)
	par := New(lex)

	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements expected 3 statements, got=%d",
			len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement. got=%T", statement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return', got %q",
				returnStatement.TokenLiteral())
		}
	}
}
