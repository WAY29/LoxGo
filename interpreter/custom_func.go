package interpreter

import (
	"fmt"

	"github.com/WAY29/LoxGo/parser"
)

type LoxCustomFunc struct { // Impl LoxCallable
	closure     *Environment
	declaration *parser.Function
	name        string
}

func NewLoxCustomFunc(declaration *parser.Function, closure *Environment) *LoxCustomFunc {
	var funcName string
	if declaration.Name != nil {
		funcName = declaration.Name.GetValue()
	}
	return &LoxCustomFunc{
		closure:     closure,
		declaration: declaration,
		name:        funcName,
	}
}

func (f *LoxCustomFunc) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	for i, param := range f.declaration.Params {
		f.closure.define(param.GetValue(), arguments[i])
	}

	environ := NewEnvironment(f.closure)
	if block, ok := f.declaration.Body.(*parser.Block); !ok {
		return nil, NewRuntimeError("Func %s body invalid.", f.declaration.Name.GetValue())
	} else {
		result, err := interpreter.executeBlock(block, environ)
		return result, err
	}
}

func (f *LoxCustomFunc) arity() int {
	return len(f.declaration.Params)
}

func (f *LoxCustomFunc) String() string {
	return fmt.Sprintf("<fn %s>", f.name)
}

func (f *LoxCustomFunc) Name() string {
	return f.name
}
