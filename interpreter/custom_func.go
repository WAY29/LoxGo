package interpreter

import (
	"fmt"

	"github.com/WAY29/LoxGo/parser"
)

type LoxCustomFunc struct { // Impl LoxCallable
	parentEnvironment *Environment
	declaration       *parser.Function
	name              string
	isInitializer     bool
}

func NewLoxCustomFunc(declaration *parser.Function, parentEnviron *Environment, isInitializer bool) *LoxCustomFunc {
	var funcName string
	if declaration.Name != nil {
		funcName = declaration.Name.GetValue()
	}
	return &LoxCustomFunc{
		parentEnvironment: parentEnviron,
		declaration:       declaration,
		name:              funcName,
		isInitializer:     isInitializer,
	}
}

func (f *LoxCustomFunc) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environ := NewEnvironment(f.parentEnvironment)
	for i, param := range f.declaration.Params {
		environ.define(param.GetValue(), arguments[i])
	}

	closure := NewEnvironment(environ)
	if block, ok := f.declaration.Body.(*parser.Block); !ok {
		return nil, NewRuntimeError(nil, "Func %s body invalid.", f.declaration.Name.GetValue())
	} else {
		result, err := interpreter.executeBlock(block, closure)
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

func (f *LoxCustomFunc) bind(instance *LoxInstance) *LoxCustomFunc {
	environ := NewEnvironment(f.parentEnvironment)
	environ.define("this", instance)
	return NewLoxCustomFunc(f.declaration, environ, f.isInitializer)
}
