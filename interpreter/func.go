package interpreter

import "fmt"

type CallableFunc = func(interpreter *Interpreter, arguments []interface{}) (interface{}, error)

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	String() string
	Name() string
}

type LoxBuiltinFunc struct { // Impl LoxCallable
	name       string
	argsNumber int
	callback   CallableFunc
}

func (f *LoxBuiltinFunc) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return f.callback(interpreter, arguments)
}

func (f *LoxBuiltinFunc) arity() int {
	return f.argsNumber
}

func (f *LoxBuiltinFunc) String() string {
	return fmt.Sprintf("<builtin-fn %s>", f.name)
}

func (f *LoxBuiltinFunc) Name() string {
	return f.name
}

func NewLoxBuiltinFunc(callback CallableFunc, name string, arity int) *LoxBuiltinFunc {
	return &LoxBuiltinFunc{
		name:       name,
		argsNumber: arity,
		callback:   callback,
	}
}
