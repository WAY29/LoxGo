package interpreter

import (
	"fmt"
	"strconv"
	"time"
)

func _clock(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}

func _string(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return fmt.Sprintf("%v", arguments[0]), nil
}

func _int(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	value := arguments[0]
	switch v := value.(type) {
	case string:
		return strconv.Atoi(v)
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case int64:
		return v, nil
	case int:
		return v, nil
	}
	return nil, NewTypeConvertError(value, "int")
}

func _float(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	value := arguments[0]
	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case float32:
		return value, nil
	case float64:
		return value, nil
	}
	return nil, NewTypeConvertError(value, "float")
}

func _bool(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return isTruthy(arguments[0]), nil
}

func _type(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return fmt.Sprintf("%T", arguments[0]), nil
}

func NewBuiltinEnvironments() *Environment {
	globals := NewEnvironment(nil)
	globals.define("clock", NewLoxBuiltinFunc(_clock, "clock", 0))
	globals.define("type", NewLoxBuiltinFunc(_type, "type", 1))
	globals.define("int", NewLoxBuiltinFunc(_int, "int", 1))
	globals.define("float", NewLoxBuiltinFunc(_float, "float", 1))
	globals.define("bool", NewLoxBuiltinFunc(_bool, "bool", 1))
	globals.define("string", NewLoxBuiltinFunc(_string, "string", 1))
	return globals
}
