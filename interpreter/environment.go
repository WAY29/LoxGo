package interpreter

import (
	"reflect"

	"github.com/WAY29/LoxGo/lexer"
	"github.com/WAY29/LoxGo/parser"
)

type Environment struct {
	enclosing *Environment
	values    map[string]Variable
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]Variable),
	}
}

func (e *Environment) get(name *lexer.Token) interface{} {
	if v, ok := e.values[name.GetValue()]; ok {
		if v.Value == nil {
			panic(NewRuntimeError(name, "Access empty variable '%s'.", name))
		}
		return v.Value
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	panic(NewRuntimeError(name, "Undefined variable '%s'.", name))
}
func (e *Environment) ancestor(distance int) *Environment {
	var environ *Environment = e
	for i := 0; i < distance; i++ {
		environ = environ.enclosing
	}
	return environ
}

func (e *Environment) getAt(distance int, name string) interface{} {
	environ := e.ancestor(distance)
	if v, ok := environ.values[name]; ok {
		if v.Value == nil {
			panic(NewRuntimeError(nil, "Access empty variable '%s'.", name))
		}
		return v.Value
	}
	panic(NewRuntimeError(nil, "Undefined variable '%s'.", name))
}

func (e *Environment) getWithBool(name string) (interface{}, bool) {
	if v, ok := e.values[name]; ok {
		return v.Value, true
	}
	if e.enclosing != nil {
		return e.enclosing.getWithBool(name)
	}
	return nil, false
}

func (e *Environment) set(name string, value interface{}) {
	if value == nil {
		e.values[name] = Variable{
			Value: value,
			Type:  reflect.Ptr,
		}
	} else {
		rv := reflect.TypeOf(value)
		e.values[name] = Variable{
			Value: value,
			Type:  rv.Kind(),
		}
	}

}

func (e *Environment) define(name string, value interface{}) {
	e.set(name, value)
}

func (e *Environment) assign(name string, value interface{}) {
	if _, ok := e.values[name]; ok {
		e.set(name, value)
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}
}

func (e *Environment) assignAt(distance int, name *lexer.Token, value parser.Expr) {
	e.ancestor(distance).set(name.GetValue(), value)
}
