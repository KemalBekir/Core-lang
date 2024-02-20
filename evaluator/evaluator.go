package evaluator

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/object"
	"fmt"
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Evaluate(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evaluateStatements(node, env)

	case *ast.BlockStatement:
		return evaluateBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, env)

	case *ast.ReturnStatement:
		value := Evaluate(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}

	case *ast.VarStatement:
		value := Evaluate(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Evaluate(node.Right, env)
		if isError(right) {
			return right
		}
		return evaluatePrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Evaluate(node.Left, env)
		if isError(left) {
			return left
		}
		right := Evaluate(node.Right, env)
		if isError(right) {
			return right
		}
		return evaluateInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evaluateIfExpression(node, env)

	case *ast.Identifier:
		return evaluateIdentifier(node, env)

	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		return &object.Function{Parameters: parameters, Env: env, Body: body}

	}

	return nil
}

func evaluateStatements(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Evaluate(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
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
	// fmt.Printf("Operator: %s, Left: %v, Right: %v\n", operator, left, right)
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		result := evaluateIntegerInfixExpression(operator, left, right)
		// fmt.Printf("Intermediate result: %v\n", result)
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
	// fmt.Printf("Left value: %d, Right value: %d\n", leftValue, rightValue)
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

func evaluateBlockStatement(block *ast.BlockStatement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Evaluate(statement, environment)
		fmt.Printf("Evaluated: %#v\n", result)
		if result != nil {
			returnType := result.Type()
			if returnType == object.RETURN_VALUE_OBJ || returnType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evaluateIfExpression(
	ife *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := Evaluate(ife.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Evaluate(ife.Consequence, env)
	} else if ife.Alternative != nil {
		return Evaluate(ife.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evaluateIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if value, ok := env.Get(node.Value); ok {
		return value
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}
