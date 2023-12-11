package main

import (
	"fmt"
	"os"

	"github.com/forfd8960/simpleinterpreter/repl"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("-------starting simple interpreter-------")
		fmt.Println("feel free to type expressions")
		repl.Start(os.Stdin, os.Stdout)
	} else if len(args) > 1 {
		file := os.Args[1]
		repl.RunScript(file)
	}
}
