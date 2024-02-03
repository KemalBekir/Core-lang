package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	INVALID = "INVALID"
	END     = "END"

	EQ     = "=="
	NOT_EQ = "!="

	// Identifiers and literals
	IDENT  = "IDENT" // variable names,
	INT    = "INT"   // integer numbers
	STRING = "STRING"

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
	FUNCTION    = "FUNCTION"
	VAR         = "VAR"
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	RETURN      = "RETURN"
	INT_TYPE    = "INT_TYPE"
	STRING_TYPE = "STRING_TYPE"
)

var tokenDictionary = map[string]TokenType{
	"function": FUNCTION,
	"var":      VAR,
	"int":      INT_TYPE,
	"string":   STRING_TYPE,
	"true":     TRUE,
	"false":    FALSE,
	"return":   RETURN,
}

func LookupIdentifier(identifier string) TokenType {
	if currentToken, ok := tokenDictionary[identifier]; ok {
		return currentToken
	}
	return IDENT
}
