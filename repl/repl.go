package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

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
			fmt.Println("parser err: ", err)
			continue
		}

		result, err := eval.Eval(program, env)
		if err != nil {
			fmt.Println("eval err: ", err)
			continue
		}

		fmt.Printf("eval result: %+v\n", result)
		if result != nil {
			fmt.Printf("eval result: %s\n", result.Inspect())
		}
	}
}

func RunScript(file string) error {
	bs, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	tokens, err := lexer.TokensFromInput(string(bs))
	if err != nil {
		fmt.Println("lexer err: ", err)
		return err
	}

	p := parser.NewParser(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Println("parser err: ", err)
		return err
	}

	env := object.NewEnvironment()
	result, err := eval.Eval(program, env)
	if err != nil {
		fmt.Println("eval err: ", err)
		return err
	}

	fmt.Printf("\neval result: %+v\n", result)
	if result != nil {
		fmt.Printf("\neval: %s\n", result.Inspect())
	}
	return nil
}
