package main

import (
	"fmt"
	"os"

	"github.com/WAY29/LoxGo/lox"
)

//go:generate go run ./tools/ast/generator.go ./parser
func main() {
	x := lox.NewLox()

	argsLen := len(os.Args)
	if argsLen > 2 {
		fmt.Println("Usage LoxGo [script]")
		os.Exit(1)
	} else if argsLen == 2 {
		x.RunFile(os.Args[1])
	} else {
		x.RunPrompt()
	}
}
