package main

import (
	"fmt"
	"os"

	"github.com/forfd8960/simpleinterpreter/repl"
)

func main() {
	fmt.Println("-------starting simple interpreter-------")
	fmt.Println("feel free to type expressions")
	repl.Start(os.Stdin, os.Stdout)
}
