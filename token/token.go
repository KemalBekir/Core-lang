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
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"

	LESS_THEN    = "<"
	GREATER_THEN = ">"

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"

	LEFT_PARANTHESIS  = "("
	RIGHT_PARANTHESIS = ")"
	LEFT_CURLY_BRACE  = "{"
	RIGHT_CURLY_BRACE = "}"
	LEFT_BRACKET      = "["
	RIGHT_BRACKET     = "]"

	// Keywords
	FUNCTION    = "FUNCTION"
	VAR         = "VAR"
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	RETURN      = "RETURN"
	IF          = "IF"
	ELSE        = "ELSE"
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
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
}

func LookupIdentifier(identifier string) TokenType {
	if currentToken, ok := tokenDictionary[identifier]; ok {
		return currentToken
	}
	return IDENT
}
