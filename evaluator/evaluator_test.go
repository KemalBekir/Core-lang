package evaluator

import (
	"Go-Tutorials/Core-lang/lexer"
	"Go-Tutorials/Core-lang/object"
	"Go-Tutorials/Core-lang/parser"
	"fmt"
	"testing"
)

func TestEvaluateIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"3", 3},
		{"33", 33},
		{"-3", -3},
		{"-13", -13},
		{"30 + 30 + 30 + 10 - 80", 20},
		{"5 * 5 * 5 * 1 * 1", 125},
		{"-40 + 60 - 10", 10},
		{"3 * 3 + 60", 69},
		{"3 + 3 * 22", 69},
		{"20 - 2 * 10", 0},
		{"100 / 10 * 10 - 31", 69},
		{"3 * (17 + 6)", 69},
		{"3 * 3 * 3 + 42", 69},
		{"3 * (3 * 3) + 42", 69},
		{"(10 + 3 * 30 + 20 / 2) * 2 + -120", 100},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEvaluate(input string) object.Object {
	lex := lexer.New(input)
	par := parser.New(lex)
	program := par.ParseProgram()
	environment := object.NewEnvironment()

	evaluated := Evaluate(program, environment)

	// If the evaluated object is a return value, extract the inner object
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return evaluated
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Object has wrong value. Got %d, want %d",
			result.Value, expected)
		return false
	}

	return true
}

func TestEvaluateBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not Boolean. Got - %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Object has wrong value. Got %t, Want %t",
			result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!8", false},
		{"!!false", false},
		{"!!true", true},
		{"!!9", true},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 3", 3},
		{"return 33; 6;", 33},
		{"return 3 * 23; 6;", 69},
		{"3; return 3 * 23; 6;", 69},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		fmt.Printf("Input: %s\n", tt.input)
		fmt.Printf("Evaluated: %#v\n", evaluated)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
