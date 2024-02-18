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
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 3;", 3},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("programStataments does not contain 1 statement. Got %d",
				len(program.Statements))
		}

		statement := program.Statements[0]
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("statement not *ast.ReturnStatement. Got %T", statement)
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Fatalf("returnStatement.TokenLiteral not 'return', got %q",
				returnStatement.TokenLiteral())
		}

		if testLiteralExpression(t, returnStatement.ReturnValue, tt.expectedValue) {
			return
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

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
	}

	for _, tt := range infixTests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got %T",
				program.Statements[0])
		}

		// fmt.Println("Input:", tt.input)
		// fmt.Println("Expected:", tt.leftValue, tt.operator, tt.rightValue)
		// fmt.Println("Actual:", statement.Expression)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. Got %d\n",
				1, len(program.Statements))
		}

		if !testInfixExpression(t, statement.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. Got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. Got %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. Got %q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
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
		// {
		// 	"a + add(b * c) + d",
		// 	"((a + add((b * c))) + d)",
		// },
		// {
		// 	"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
		// 	"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		// },
		// {
		// 	"add(a + b + c * d / f + g)",
		// 	"add((((a + b) + ((c * d) / f)) + g))",
		// },
		// {
		// 	"a * [1, 2, 3, 4][b * c] * d",
		// 	"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		// },
		// {
		// 	"add(a * b[2], b[1], 2 * [1, 2][1])",
		// 	"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		// },
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	identifier, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Exp not *ast.Identifier. Got %T", exp)
		return false
	}

	if identifier.Value != value {
		t.Errorf("Identifier.Value not %s. Got %s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral not %s. Got %s", value,
			identifier.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. Got %T", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. Got %t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t, Got %s",
			value, boolean.TokenLiteral())
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. Got %d\n",
			1, len(program.Statements))
	}

	statament, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatament. Got %T",
			program.Statements[0])
	}

	expression, ok := statament.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Expression is not ast.IfExpression. Got %T",
			statament.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statament. Got %d\n",
			len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. Got %T",
			expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative.Statements was not nil. Got %+v", expression.Alternative)
	}
}

func TestParsingFunctionLiteral(t *testing.T) {
	input := `function(a, b) { a + b; }`

	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. Got %d\n",
			1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Stataments[0] is not ast.ExpressionStatement. Got %T",
			program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not ast.FunctionLiteral. Got %T",
			statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("Function literal parameters wrong. Require 2, got %d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "a")
	testLiteralExpression(t, function.Parameters[1], "b")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in function.Body.Statements, but found %d.\n",
			len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected function body statament to be ast.ExpressionStatement. Got %T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "a", "+", "b")

}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input              string
		expectedParameters []string
	}{
		{input: "function() {};", expectedParameters: []string{}},
		{input: "function(a) {};", expectedParameters: []string{"a"}},
		{input: "function(a, b, c) {};", expectedParameters: []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParameters) {
			t.Errorf("Expected parameter count: %d; found: %d.\n",
				len(tt.expectedParameters), len(function.Parameters))
		}

		for i, identifier := range tt.expectedParameters {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}
	}
}
