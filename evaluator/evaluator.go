package evaluator

import (
	"Go-Tutorials/Core-lang/ast"
	"Go-Tutorials/Core-lang/object"
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
