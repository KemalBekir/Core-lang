package parser

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/lexer"
	"fmt"
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

func TestInvalidVarStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{"var string b = 5;", "Type mismatch: cannot assign int to string variable"},
		{"var int a = \"abc\";", "Type mismatch: cannot assign string to int variable"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		par := New(lex)
		_ = par.ParseProgram()

		if len(par.Errors()) == 0 {
			t.Fatalf("Expected type mismatch error but got none")
		}

		isError := false
		for _, message := range par.Errors() {
			if message == tt.expectedMsg {
				isError = true
				break
			}
		}

		if !isError {
			t.Errorf("expected a type mismatch error with message '%s', but got: %v", tt.expectedMsg, par.Errors())
		}
	}
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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "7;"

	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements. got=%d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not *ast.IntegerLiteral. got=%T", statement.Expression)
	}
	if literal.Value != 7 {
		t.Errorf("literal.Value not %d. got=%d", 7, literal.Value)
	}
	if literal.TokenLiteral() != "7" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "7", literal.TokenLiteral())
	}
}

func testIntegerLiteral(t *testing.T, intL ast.Expression, value int64) bool {
	integer, ok := intL.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("intL not *ast.IntegerLiteral. got=%T", intL)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!4;", "!", 4},
		{"-9;", "-", 9},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement is not ast.PrefixExpression. got=%T", statement.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not '%s'. got=%s", tt.operator, expression.Operator)
		}
		if !testIntegerLiteral(t, expression.Right, tt.integerValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
