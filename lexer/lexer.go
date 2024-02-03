package lexer

import "Go-Tutorials/Core-lang/token"

type Lexer struct {
	input    string
	position int  // current position in input (point to current char)
	readPos  int  // current reading position in input (after current char)
	char     byte // current char under examination
}

func New(input string) *Lexer {
	lexInstance := &Lexer{input: input}
	lexInstance.readCharacter()
	return lexInstance
}

func (lex *Lexer) readCharacter() {
	if lex.readPos >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.readPos]
	}

	lex.position = lex.readPos
	lex.readPos += 1
}

func (lex *Lexer) NextToken() token.Token {
	var currentToken token.Token

	lex.ignoreWhitespace()

	switch lex.char {
	case '=':
		if lex.peekAheadCharacter() == '=' {
			char := lex.char
			lex.readCharacter()
			literal := string(char) + string(lex.char)
			currentToken = token.Token{Type: token.EQ, Literal: literal}
		} else {
			currentToken = newToken(token.ASSIGN_OP, lex.char)
		}
	case ';':
		currentToken = newToken(token.SEMICOLON, lex.char)
	case '(':
		currentToken = newToken(token.LEFT_PARANTHESIS, lex.char)
	case ')':
		currentToken = newToken(token.RIGHT_PARANTHESIS, lex.char)
	case ',':
		currentToken = newToken(token.COMMA, lex.char)
	case '+':
		currentToken = newToken(token.PLUS, lex.char)
	case '{':
		currentToken = newToken(token.LEFT_CURLY_BRACE, lex.char)
	case '}':
		currentToken = newToken(token.RIGHT_CURLY_BRACE, lex.char)
	case '"':
		currentToken.Type = token.STRING
		currentToken.Literal = lex.readString()

	case 0:
		currentToken.Literal = ""
		currentToken.Type = token.END
	default:
		if isAlphabetic(lex.char) {
			currentToken.Literal = lex.searchIdentifier()
			currentToken.Type = token.LookupIdentifier(currentToken.Literal)
			return currentToken
		} else if isNumber(lex.char) {
			currentToken.Type = token.INT
			currentToken.Literal = lex.checkNumber()
			return currentToken
		} else {
			currentToken = newToken(token.INVALID, lex.char)
		}
	}

	lex.readCharacter()
	return currentToken
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	if tokenType == "=" {
		return token.Token{Type: token.ASSIGN_OP, Literal: string(char)}
	} else if tokenType == ";" {
		return token.Token{Type: token.SEMICOLON, Literal: string(char)}
	} else {
		return token.Token{Type: tokenType, Literal: string(char)}
	}
}

func (lex *Lexer) searchIdentifier() string {
	position := lex.position
	for isAlphabetic(lex.char) || isNumber(lex.char) {
		lex.readCharacter()
	}

	return lex.input[position:lex.position]
}

func isAlphabetic(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (lex *Lexer) ignoreWhitespace() {
	for lex.char == ' ' || lex.char == '\t' || lex.char == '\n' || lex.char == '\r' {
		lex.readCharacter()
	}
}

func (lex *Lexer) checkNumber() string {
	position := lex.position
	for isNumber(lex.char) {
		lex.readCharacter()
	}

	return lex.input[position:lex.position]
}

func isNumber(char byte) bool {
	return '0' <= char && char <= '9'
}

func (lex *Lexer) peekAheadCharacter() byte {
	if lex.readPos >= len(lex.input) {
		return 0
	} else {
		return lex.input[lex.readPos]
	}
}

func (lex *Lexer) readString() string {
	position := lex.position + 1
	for {
		lex.readCharacter()
		if lex.char == '"' || lex.char == 0 {
			break
		}
	}

	return lex.input[position:lex.position]
}
