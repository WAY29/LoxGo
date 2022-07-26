package interpreter

import (
	"fmt"
	"math"
	"sync/atomic"

	"github.com/WAY29/LoxGo/lexer"
	"github.com/WAY29/LoxGo/parser"
)

var (
	ONE uint64 = 1
)

type Interpreter struct { // impl ExprVisitor, StmtVisitor
	globals     *Environment
	environment *Environment
	locals      map[parser.Expr]int

	block *parser.Block
	while *parser.While

	stack     *Stack
	stackSize uint64
}

func NewInterpreter() *Interpreter {
	globals := NewBuiltinEnvironments()
	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[parser.Expr]int),
		block:       nil,
		stack:       NewStack(nil, nil, nil),
	}
}

func (i *Interpreter) newEnvironmentState(newEnviron *Environment) func() {
	oldEnviron := i.environment
	i.environment = newEnviron

	return func() {
		i.environment = oldEnviron
	}
}

func (i *Interpreter) newBlockState(newBlock *parser.Block) func() {
	oldBlock := i.block
	i.block = newBlock

	return func() {
		i.block = oldBlock
	}
}

func (i *Interpreter) newWhileState(newWhile *parser.While) func() {
	oldWhile := i.while
	i.while = newWhile

	return func() {
		i.while = oldWhile
	}
}

func (i *Interpreter) newStackState(call *parser.Call, callee LoxCallable) func() {
	oldStack := i.stack
	newStack := NewStack(i.stack, call, callee)
	i.stack = newStack
	atomic.AddUint64(&i.stackSize, 1)
	if i.stackSize >= 8192 {
		panic(NewRuntimeError("Stack oversize: Can't have more than 8192 stack."))
	}

	return func() {
		//todo debug
		// fmt.Printf("return stack: %s\n", i.stack)

		atomic.AddUint64(&i.stackSize, -ONE)
		i.stack = oldStack
	}
}

func (i *Interpreter) Interpret(statemants []parser.Stmt) []interface{} {
	var (
		err     error
		result  interface{}
		results = make([]interface{}, len(statemants))
	)
	for n, stmt := range statemants {
		result, err = i.execute(stmt)
		if err != nil {
			panic(err)
		}
		results[n] = result
	}
	return results
}

func (i *Interpreter) evaluate(expr parser.Expr) (interface{}, error) {
	if expr == nil {
		return nil, nil
	}
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt parser.Stmt) (interface{}, error) {
	if stmt == nil {
		return nil, nil
	}
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(block *parser.Block, environment *Environment) (result interface{}, err error) {
	block = parser.NewBlock(block.Statements, false, i.block)

	defer i.newEnvironmentState(environment)()
	defer i.newBlockState(block)()

	for _, stmt := range block.Statements {
		if block.Stop {
			block.Stop = false
			break
		}
		_, err = i.execute(stmt)
		if err != nil {
			return nil, err
		}

		if i.stack.call != nil && i.stack.call.ReturnValue != nil {
			block.Stop = false
			value := i.stack.call.ReturnValue
			return value, nil
		}
	}

	return nil, nil
}

func (i *Interpreter) Resolve(expr parser.Expr, depth int) {
	// fmt.Printf("debug: set %#v: %d\n", expr, depth)
	i.locals[expr] = depth
}

func (i *Interpreter) lookUpVariable(name *lexer.Token, expr parser.Expr) (interface{}, error) {
	// fmt.Printf("debug: lookup expr: %#v locals:%#v\n", expr, i.locals)
	if distance, ok := i.locals[expr]; !ok {
		return i.globals.get(name.GetValue()), nil
	} else {
		// fmt.Printf("debug: find in %d distance scope: %s\n", distance, name.GetValue())
		return i.environment.getAt(distance, name.GetValue()), nil
	}
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *parser.Logical) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.GetType() == lexer.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else { // AND
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitGroupingExpr(expr *parser.Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *parser.Unary) (result interface{}, err error) {
	// 强制整数转换
	defer func() {
		if err == nil {
			if v, ok := result.(float64); ok && v == math.Trunc(v) {
				result = int(v)
			}
		}
	}()

	var (
		ok    bool
		ve    *parser.Variable
		vi    interface{}
		v     float64
		isVar bool = false
	)

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	if expr.Prefix {
		switch expr.Operator.GetType() {
		case lexer.MINUS:
			if v, ok := interfaceToFloat64(right); ok {
				return -1 * v, nil
			}
			return nil, NewConvertError(right, "float", "")
		case lexer.PLUS:
			if v, ok := interfaceToFloat64(right); ok {
				return v, nil
			}
			return nil, NewConvertError(right, "float", "")
		case lexer.BANG:
			return !isTruthy(right), nil
		case lexer.PLUSPLUS:
			if ve, ok = expr.Right.(*parser.Variable); ok {
				if vi, ok = i.environment.getWithBool(ve.Name.GetValue()); ok {
					isVar = true
				}
			} else {
				vi = right
			}
			if v, ok = interfaceToFloat64(vi); ok {
				if isVar {
					i.environment.assign(ve.Name.GetValue(), v+1)
				}
				return v + 1, nil
			}
			return nil, NewConvertError(vi, "float", "")
		case lexer.MINUSMINUS:
			if ve, ok = expr.Right.(*parser.Variable); ok {
				if vi, ok = i.environment.getWithBool(ve.Name.GetValue()); ok {
					isVar = true
				}
			} else {
				vi = right
			}
			if v, ok = interfaceToFloat64(vi); ok {
				if isVar {
					i.environment.assign(ve.Name.GetValue(), v-1)
				}
				return v - 1, nil
			}
			return nil, NewConvertError(vi, "float", "")
		}
	} else {
		switch expr.Operator.GetType() {
		case lexer.PLUSPLUS:
			if ve, ok = expr.Right.(*parser.Variable); ok {
				if vi, ok = i.environment.getWithBool(ve.Name.GetValue()); ok {
					isVar = true
				}
			} else {
				vi = right
			}
			if v, ok = interfaceToFloat64(vi); ok {
				if isVar {
					i.environment.assign(ve.Name.GetValue(), v+1)
				}
				return v, nil
			}
			return nil, NewConvertError(vi, "float", "")
		case lexer.MINUSMINUS:
			if ve, ok = expr.Right.(*parser.Variable); ok {
				if vi, ok = i.environment.getWithBool(ve.Name.GetValue()); ok {
					isVar = true
				}
			} else {
				vi = right
			}
			if v, ok = interfaceToFloat64(vi); ok {
				if isVar {
					i.environment.assign(ve.Name.GetValue(), v-1)
				}
				return v, nil
			}
			return nil, NewConvertError(vi, "float", "")
		}
	}

	return nil, NewConvertError(right, "float", "invalid unary expression")
}

func (i *Interpreter) VisitBinaryExpr(expr *parser.Binary) (result interface{}, err error) {

	// 强制整数转换
	defer func() {
		if err == nil {
			if v, ok := result.(float64); ok && v == math.Trunc(v) {
				result = int(v)
			}
		}
	}()

	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.GetType() {
	case lexer.MINUS:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v - v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		return nil, NewConvertError(left, "float", "")
	case lexer.SLASH:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v / v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		return nil, NewConvertError(left, "float", "")
	case lexer.STAR:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v * v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		return nil, NewConvertError(left, "float", "")
	case lexer.PLUS:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v + v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		if v, ok := left.(string); ok {
			if v2, ok := right.(string); ok {
				return v + v2, nil
			}
			return nil, NewConvertError(right, "string", "")
		}
		return nil, NewConvertError(left, "string / float", "")
	case lexer.GREATER:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v > v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		return nil, NewConvertError(left, "float", "")
	case lexer.GREATER_EQUAL:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v >= v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		return nil, NewConvertError(left, "float", "")
	case lexer.LESS:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v < v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		return nil, NewConvertError(left, "float", "")
	case lexer.LESS_EQUAL:
		if v, ok := interfaceToFloat64(left); ok {
			if v2, ok := interfaceToFloat64(right); ok {
				return v <= v2, nil
			}
			return nil, NewConvertError(right, "float", "")
		}
		return nil, NewConvertError(left, "float", "")
	case lexer.EQUAL_EQUAL:
		return isEqual(left, right), nil
	case lexer.BANG_EQUAL:
		return !isEqual(left, right), nil
	}

	return nil, NewConvertError(nil, "", "invalid binary expression")
}

func (i *Interpreter) VisitCallExpr(expr *parser.Call) (result interface{}, err error) {
	var (
		callee, v interface{}
	)
	expr = parser.NewCall(expr.Callee, expr.Paren, expr.Arguments, nil)
	defer i.newBlockState(nil)()
	defer i.newWhileState(nil)()

	callee, err = i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	if calleeFunc, ok := callee.(LoxCallable); !ok {
		return nil, NewRuntimeError("Can only call functions and classes.")
	} else {
		defer i.newStackState(expr, calleeFunc)()

		argsLen := len(expr.Arguments)
		if calleeFunc.arity() != argsLen {
			return nil, NewRuntimeError("Excepted %d arguments but got %d.", calleeFunc.arity(), argsLen)
		}
		arguments := make([]interface{}, argsLen)
		for n, argument := range expr.Arguments {
			v, err = i.evaluate(argument)
			if err != nil {
				return nil, err
			}
			arguments[n] = v
		}
		result, err = calleeFunc.call(i, arguments)
		if err != nil {
			panic(NewRuntimeError("%v", err))
		}
		return
	}
}

func (i *Interpreter) VisitVariableExpr(expr *parser.Variable) (interface{}, error) {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) VisitTernaryExpr(expr *parser.Ternary) (interface{}, error) {
	var (
		condResult, result interface{}
		err                error
	)
	if condResult, err = i.evaluate(expr.Condition); err != nil {
		return nil, err
	}

	if isTruthy(condResult) {
		result, err = i.evaluate(expr.ThenExpr)
	} else {
		result, err = i.evaluate(expr.ElseExpr)
	}

	return result, err

}

func (i *Interpreter) VisitAssignExpr(expr *parser.Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	if distance, ok := i.locals[expr]; !ok {
		i.globals.assign(expr.Name.GetValue(), value)
	} else {
		i.environment.assignAt(distance, expr.Name, expr)
	}

	return value, nil
}

func (i *Interpreter) VisitLambdaExpr(expr *parser.Lambda) (interface{}, error) {
	if function, ok := expr.Function.(*parser.Function); ok {
		return NewLoxCustomFunc(function, i.environment), nil
	}
	return nil, NewRuntimeError("Invalid Lmabda in line %d.", expr.Token.GetLine())
}

func (i *Interpreter) VisitExpressionStmt(stmt *parser.Expression) (interface{}, error) {
	result, err := i.evaluate(stmt.Expr)
	return result, err
}

func (i *Interpreter) VisitFunctionStmt(stmt *parser.Function) (interface{}, error) {
	// fmt.Printf("debug: set function: %s\n", stmt.Name.GetValue())
	i.environment.define(stmt.Name.GetValue(), NewLoxCustomFunc(stmt, NewEnvironment(i.environment)))
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *parser.If) (interface{}, error) {
	var (
		condResult interface{}

		err error
	)
	if condResult, err = i.evaluate(stmt.Condition); err != nil {
		return nil, err
	}

	if isTruthy(condResult) {
		_, err = i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		_, err = i.execute(stmt.ElseBranch)
	}

	return nil, err
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.Print) (interface{}, error) {
	result, err := i.evaluate(stmt.Expr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", result)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *parser.Return) (interface{}, error) {
	// if i.block == nil {
	// 	return nil, NewRuntimeError("Invalid return usage in line %d.", stmt.Keyword.GetLine())
	// }

	recursiveBlockStop(i.block)
	recursiveWhileStop(i.while)
	value, err := i.evaluate(stmt.Value)
	if err != nil {
		return nil, err
	}
	if i.stack.call != nil {
		i.stack.call.ReturnValue = value
	}

	return nil, nil
}

func (i *Interpreter) VisitBreakStmt(stmt *parser.Break) (interface{}, error) {
	if stmt.Parent == nil {
		return nil, NewRuntimeError("Invalid break usage in line %d.", stmt.Token.GetLine())
	}
	if while, ok := stmt.Parent.(*parser.While); ok {
		while.Stop = true
	}
	if stmt.Block != nil {
		if block, ok := stmt.Block.(*parser.Block); ok {
			block.Stop = true
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitContinueStmt(stmt *parser.Continue) (interface{}, error) {
	if stmt.Block == nil {
		return nil, NewRuntimeError("Invalid break usage in line %d.", stmt.Token.GetLine())
	}

	if block, ok := stmt.Block.(*parser.Block); ok {
		block.Stop = true
	}
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *parser.While) (interface{}, error) {
	var (
		r   interface{}
		err error
	)
	defer i.newWhileState(stmt)()

	for {
		if stmt.Stop {
			stmt.Stop = false
			break
		}

		r, err = i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !isTruthy(r) {
			break
		}
		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}

	return nil, err
}

func (i *Interpreter) VisitVarStmt(stmt *parser.Var) (interface{}, error) {
	for n := range stmt.Names {
		result, err := i.evaluate(stmt.Initializers[n])

		if err != nil {
			return nil, err
		}
		i.environment.define(stmt.Names[n].GetValue(), result)
	}

	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *parser.Block) (interface{}, error) {
	result, err := i.executeBlock(stmt, NewEnvironment(i.environment))
	return result, err
}
