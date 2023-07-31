package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/forfd8960/simpleinterpreter/lexer"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

const PROMT = ">>%s"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMT, " ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		text := scanner.Text()
		lxer := lexer.NewLexer(text)

		var tok *tokens.Token
		var err error

		for {
			tok, err = lxer.NextToken()
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Printf("%+v\n", tok)

			if tok.TkType == tokens.EOF {
				break
			}
		}
	}
}
