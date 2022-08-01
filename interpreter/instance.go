package interpreter

import "github.com/WAY29/LoxGo/lexer"

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (i *LoxInstance) String() string {
	return i.class.className + " instance"
}

func (i *LoxInstance) get(name *lexer.Token) (result interface{}) {
	var ok bool
	if result, ok = i.fields[name.GetValue()]; ok {
		return
	}
	method := i.class.findMethod(name.GetValue())
	if method != nil {
		return method.bind(i)
	}

	panic(NewRuntimeError(name, "Undefined property '%s'.", name.GetValue()))
}

func (i *LoxInstance) set(name *lexer.Token, value interface{}) {
	i.fields[name.GetValue()] = value
}
