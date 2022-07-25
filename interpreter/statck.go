package interpreter

import (
	"fmt"

	"github.com/WAY29/LoxGo/parser"
)

type Stack struct {
	parent *Stack
	call   *parser.Call
	callee LoxCallable
}

func NewStack(parent *Stack, call *parser.Call, callee LoxCallable) *Stack {
	return &Stack{
		parent: parent,
		call:   call,
		callee: callee,
	}
}

func (s *Stack) String() string {
	value := s.callee.Name()
	stack := s
	for {
		if stack.parent != nil && stack.parent.callee != nil {
			stack = stack.parent
		} else {
			break
		}
		value += fmt.Sprintf(" -> %s", stack.callee.Name())
	}
	return value
}
