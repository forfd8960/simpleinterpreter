package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/forfd8960/simpleinterpreter/eval"
	"github.com/forfd8960/simpleinterpreter/lexer"
	"github.com/forfd8960/simpleinterpreter/object"
	"github.com/forfd8960/simpleinterpreter/parser"
)

const PROMT = ">>%s"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMT, " ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		text := scanner.Text()
		tokens, err := lexer.TokensFromInput(text)
		if err != nil {
			fmt.Println("lexer err: ", err)
			continue
		}

		p := parser.NewParser(tokens)
		program, err := p.ParseProgram()
		if err != nil {
			fmt.Println("lexer err: ", err)
			continue
		}

		result, err := eval.Eval(program, env)
		if err != nil {
			fmt.Println("lexer err: ", err)
			continue
		}

		fmt.Printf("eval result: %v\n", result)
	}
}
