package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/WAY29/LoxGo/interpreter"
	"github.com/WAY29/LoxGo/lexer"
	"github.com/WAY29/LoxGo/parser"
)

type Lox struct {
	interpreter *interpreter.Interpreter
}

func NewLox() *Lox {
	return &Lox{
		interpreter: interpreter.NewInterpreter(),
	}
}

func (lox *Lox) RunFile(file string) {
	fp, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	lox.Eval(fp)
}

func (lox *Lox) RunPrompt() {

	var (
		line    string
		results []interface{}
	)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("lox > ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line = scanner.Text()
		// 末尾补充分号
		if !strings.HasSuffix(line, ";") && len(line) > 0 {
			line += ";"
		}

		results = lox.Eval(strings.NewReader(line))
		if len(results) < 1 {
			continue
		}
		for _, result := range results {
			if result == nil {
				continue
			}
			fmt.Printf("%v [%T]\n", result, result)
		}
	}
}

func (lox *Lox) Eval(r io.Reader) []interface{} {
	// error handle
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			if _, ok := r.(*interpreter.ConvertError); ok {
				fmt.Printf("[ERROR] %s\n", r)
			} else if _, ok := r.(*interpreter.RuntimeError); ok {
				fmt.Printf("[ERROR] %s\n", r)
			} else if _, ok := r.(*parser.ParseError); ok {
				fmt.Printf("[ERROR] %s\n", r)
			} else {
				panic(r)
			}
		}
	}()

	l := lexer.NewLexer(r)
	l.ScanTokens()
	if l.GetError() != nil {
		panic(l.GetError())
	}

	p := parser.NewParaser(l.GetTokens())
	statements := p.Parse()
	if statements == nil {
		return nil
	}
	resolver := interpreter.NewResolver(lox.interpreter)
	resolver.ResolveStmts(statements)

	results := lox.interpreter.Interpret(statements)
	return results
}
