package interpreter

import "reflect"

type Variable struct {
	Value interface{}
	Type  reflect.Kind
}
