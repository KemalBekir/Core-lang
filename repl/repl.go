package repl

import (
	"Go-Tutorials/Core-lang/evaluator"
	"Go-Tutorials/Core-lang/lexer"
	"Go-Tutorials/Core-lang/object"
	"Go-Tutorials/Core-lang/parser"
	"bufio"
	"fmt"
	"io"
)

const CORE_LANG = `
_________                        .____                           
\_   ___ \  ___________   ____   |    |   _____    ____    ____  
/    \  \/ /  _ \_  __ \_/ __ \  |    |   \__  \  /    \  / ___\ 
\     \___(  <_> )  | \/\  ___/  |    |___ / __ \|   |  \/ /_/  >
 \______  /\____/|__|    \___  > |_______ (____  /___|  /\___  / 
        \/                   \/          \/    \/     \//_____/  
`

const PROMPT = "::> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	environment := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		par := parser.New(lex)

		program := par.ParseProgram()
		if len(par.Errors()) != 0 {
			printParseErrors(out, par.Errors())
			continue
		}

		evaluated := evaluator.Evaluate(program, environment)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, CORE_LANG)
	io.WriteString(out, "Opps! We ran in to some issue \n")
	io.WriteString(out, " parser errors:\n")
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
