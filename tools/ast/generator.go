package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		out string
		err error
	)
	if len(os.Args) == 2 {
		out, err = filepath.Abs(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		out = "../../parser/"
	}

	defineAst(out, "Expr", []string{
		"Ternary  : Condition Expr, ThenExpr Expr, ElseExpr Expr",
		"Assign   : Name *lexer.Token, Value Expr",
		"Binary   : Left Expr, Operator *lexer.Token, Right Expr",
		"Call     : Callee Expr, Paren *lexer.Token, Arguments []Expr, ReturnValue interface{}",
		"Grouping : Expression Expr",
		"Literal  : Value interface{}",
		"Logical  : Left Expr, Operator *lexer.Token, Right Expr",
		"Unary    : Operator *lexer.Token, Right Expr, Prefix bool",
		"Variable : Name *lexer.Token",
		"Lambda   : Token *lexer.Token, Function Stmt",
	})

	defineAst(out, "Stmt", []string{
		"Block      : Statements []Stmt, Stop bool, Parent *Block",
		"Expression : Expr Expr",
		"Function   : Name *lexer.Token, Params []*lexer.Token, Body Stmt",
		"If         : Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print      : Expr Expr",
		"Return     : Keyword *lexer.Token, Initializer Expr",
		"Var        : Names []*lexer.Token, Initializers []Expr",
		"While      : Condition Expr, Body Stmt, Stop bool, Parent *While",
		"Break      : Parent Stmt, Block Stmt, Token *lexer.Token",
		"Continue   : Block Stmt, Token *lexer.Token",
	})
}

func defineAst(out, base string, types []string) {
	var src string

	// visitor
	src += fmt.Sprintln("// Generated code, do not edit.")
	src += fmt.Sprintln("")
	src += fmt.Sprintln("package parser")
	src += defineVisitor(base, types)
	path := fmt.Sprintf("%s/%s_visitor.go", out, strings.ToLower(base))
	if err := saveFile(path, src); err != nil {
		panic(err)
	}

	// expr
	src = ""
	src += fmt.Sprintln("// Generated code, do not edit.")
	src += fmt.Sprintln("")
	src += fmt.Sprintln("package parser")
	src += fmt.Sprintln(`import "github.com/WAY29/LoxGo/lexer"`)
	src += defineBase(base, types)
	for _, t := range types {
		cls := strings.TrimRight(strings.Split(t, ":")[0], " ")
		fld := strings.TrimRight(strings.Split(t, ":")[1], " ")
		src += defineType(base, cls, fld)
	}
	path = fmt.Sprintf("%s/%s.go", out, strings.ToLower(base))
	if err := saveFile(path, src); err != nil {
		panic(err)
	}
}

func defineVisitor(base string, types []string) string {
	var src string

	// visitor interface
	src += fmt.Sprintln("")
	src += fmt.Sprintf("type %sVisitor interface {\n", base)
	for _, t := range types {
		cls := strings.TrimRight(strings.Split(t, ":")[0], " ")
		src += fmt.Sprintf("Visit%s%s(%s *%s) (interface{}, error)", cls, base, strings.ToLower(base), cls)
		src += fmt.Sprintln("")
	}
	src += fmt.Sprintln("}")
	return src
}

func defineBase(base string, types []string) string {
	var src string = ""

	// base interface
	src += fmt.Sprintf("type %s interface {", base)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("Accept(v %sVisitor) (interface{}, error)\n", base)
	src += fmt.Sprintln("}")

	return src
}

func defineType(base, cls, fld string) string {
	var src string

	src += fmt.Sprintln("")
	src += fmt.Sprintf("type %s struct {", cls)
	src += fmt.Sprintln("")

	// fields
	fs := strings.Split(fld, ",")
	for _, f := range fs {
		src += fmt.Sprintln(f)
	}
	src += fmt.Sprintln("}")

	// new func
	src += fmt.Sprintf("func New%s(", cls)
	params := []string{}
	for _, f := range fs {
		t := strings.Split(f, " ")[2]
		n := strings.ToLower(strings.Split(f, " ")[1])
		params = append(params, fmt.Sprintf("%s %s", n, t))
	}
	src += fmt.Sprintf(strings.Join(params, ","))
	src += fmt.Sprintf(") *%s {", cls)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("return &%s{", cls)
	args := []string{}
	for _, f := range fs {
		t := strings.ToLower(strings.Split(f, " ")[1])
		n := strings.Split(f, " ")[1]
		args = append(args, fmt.Sprintf("%s: %s", n, t))
	}
	src += fmt.Sprintf(strings.Join(args, ","))
	src += fmt.Sprintln("}")
	src += fmt.Sprintln("}")

	// accept func
	src += fmt.Sprintf("func (n *%s) Accept(v %sVisitor) (interface{}, error) {", cls, base)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("return v.Visit%s%s(n)", cls, base)
	src += fmt.Sprintf("}")
	src += fmt.Sprintln("")

	return src
}

func saveFile(path, src string) error {
	// gofmt
	buf, err := format.Source([]byte(src))
	if err != nil {
		return err
	}
	// save
	ioutil.WriteFile(path, buf, 0644)

	return nil
}
