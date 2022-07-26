package interpreter

import (
	"fmt"
)

type RuntimeError struct {
	Msg string
}

func NewRuntimeError(format string, a ...interface{}) *RuntimeError {
	return &RuntimeError{
		Msg: "Runtime Error: " + fmt.Sprintf(format, a...),
	}
}

func NewTypeConvertError(value interface{}, _type string) *RuntimeError {
	return &RuntimeError{
		Msg: "Runtime Error: Type conversion error: " + fmt.Sprintf("can't convert %v[%T] to %s", value, value, _type),
	}
}

func (e *RuntimeError) Error() string {
	return e.Msg
}

type ConvertError struct {
	value      interface{}
	typeString string
	extraMsg   string
}

func NewConvertError(value interface{}, typeString, msg string) *ConvertError {
	if len(msg) > 0 {
		msg = ": " + msg
	}
	return &ConvertError{
		value:      value,
		typeString: typeString,
		extraMsg:   msg,
	}
}

func (e *ConvertError) Error() string {
	if e.value == nil {
		return fmt.Sprintf("Runtime error: Convert Error: %s.", e.extraMsg)
	}
	return fmt.Sprintf("Runtime error: Convert error: can't convert %v to %s%s.", e.value, e.typeString, e.extraMsg)
}
