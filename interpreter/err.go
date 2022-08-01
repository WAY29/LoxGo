package interpreter

import (
	"fmt"

	"github.com/WAY29/LoxGo/lexer"
)

type RuntimeError struct {
	extraMsg string
	token    *lexer.Token
}

func NewRuntimeError(token *lexer.Token, format string, a ...interface{}) *RuntimeError {
	return &RuntimeError{
		extraMsg: "Runtime Error: " + fmt.Sprintf(format, a...),
	}
}

func (e *RuntimeError) Error() string {
	var (
		where string
		token = e.token
	)

	if token == nil {
		return fmt.Sprintf("Runtime Error: %s", e.extraMsg)
	}

	if token.GetType() == lexer.EOF {
		where = " at end"
	} else {
		where = fmt.Sprintf("at '%s'", token.GetValue())
	}

	return fmt.Sprintf("Runtime Error in line %d: Error %s: %s", token.GetLine(), where, e.extraMsg)
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
	return fmt.Sprintf("Runtime error: Convert error: can't convert %v[%T] to %s%s.", e.value, e.value, e.typeString, e.extraMsg)
}
