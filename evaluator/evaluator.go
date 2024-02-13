package evaluator

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Evaluate(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evaluateStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Evaluate(node.Right)
		return evaluatePrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Evaluate(node.Left)
		right := Evaluate(node.Right)
		return evaluateInfixExpression(node.Operator, left, right)
	}

	return nil
}

func evaluateStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Evaluate(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(right)
	case "-":
		return evaluateMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evaluateBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}

}

func evaluateMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evaluateInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	fmt.Printf("Operator: %s, Left: %v, Right: %v\n", operator, left, right)
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		result := evaluateIntegerInfixExpression(operator, left, right)
		fmt.Printf("Intermediate result: %v\n", result)
		return result
	default:
		return NULL
	}
}

func evaluateIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	fmt.Printf("Left value: %d, Right value: %d\n", leftValue, rightValue)
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	default:
		return NULL
	}
}
