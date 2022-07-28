package interpreter

import (
	"math"
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
	switch v := a.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	}

	return 0, false
}

func float642Int(v float64) (int, bool) {
	if v == math.Trunc(v) {
		return int(v), true
	}
	return 0, false
}

func interfaceToInt(a interface{}) (int, bool) {
	switch v := a.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
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
