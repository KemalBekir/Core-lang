package lexer

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
	if lex.position >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.readPos]
	}

	lex.position = lex.readPos
	lex.readPos += 1
}
