package interpreter

import "time"

func NewBuiltinEnvironments() *Environment {
	globals := NewEnvironment(nil)
	globals.define("clock", NewLoxBuiltinFunc(func(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
		return time.Now().Unix(), nil
	}, "clock", 0))
	return globals
}
