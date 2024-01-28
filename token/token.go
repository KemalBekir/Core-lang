package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	INVALID = "INVALID"
	END     = "END"

	// Identifiers and literals
	NAME = "NAME" // variable names,
	INT  = "INT"  // integer numbers

	// Operators
	ASSIGN_OP = "="
	PLUS      = "+"

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"

	LEFT_PARANTHESIS  = "("
	RIGHT_PARANTHESIS = ")"
	LEFT_CURLY_BRACE  = "{"
	RIGHT_CURLY_BRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
