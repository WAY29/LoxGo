package interpreter

import (
	"reflect"

	"github.com/WAY29/LoxGo/parser"
)

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	}

	return reflect.DeepEqual(a, b)
}

func interfaceToFloat64(a interface{}) (float64, bool) {
	if v, ok := a.(float64); ok {
		return v, ok
	} else if v, ok := a.(int); ok {
		return float64(v), ok
	}
	return 0, false
}

func isTruthy(a interface{}) bool {
	if a == nil {
		return false
	}
	if v, ok := a.(bool); ok {
		return v
	}

	return true
}

func recursiveBlockStop(block *parser.Block) {
	if block == nil {
		return
	}
	for {
		block.Stop = true
		if block.Parent == nil {
			break
		}
		block = block.Parent
	}
}

func recursiveWhileStop(while *parser.While) {
	if while == nil {
		return
	}
	for {
		while.Stop = true
		if while.Parent == nil {
			break
		}
		while = while.Parent
	}
}

// func interfaceToBool(a interface{}) (bool, bool) {
// 	if a == nil {
// 		return false, true
// 	}
// 	if v, ok := a.(bool); ok {
// 		return v, ok
// 	}

// 	return false, false
// }

// func interfaceToString(a interface{}) (string, bool) {
// 	if v, ok := a.(string); ok {
// 		return v, ok
// 	}
// 	return fmt.Sprintf("%v", a), true
// }
