package interpreter

import (
	"container/list"

	"github.com/WAY29/LoxGo/lexer"
	"github.com/WAY29/LoxGo/parser"
)

type Resolver struct { // impl ExprVisitor, StmtVisitor
	interpreter *Interpreter
	scopes      *list.List

	inFunction bool
}

func NewResolver(i *Interpreter) *Resolver {
	return &Resolver{
		interpreter: i,
		scopes:      list.New(),
		inFunction:  false,
	}
}

func (r *Resolver) newFunctionState() func() {
	oldInFunction := r.inFunction
	r.inFunction = true

	return func() {
		r.inFunction = oldInFunction
	}
}

func (r *Resolver) ResolveStmts(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt parser.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr parser.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr parser.Expr, name *lexer.Token) {
	n := r.scopes.Len() - 1
	for i := r.scopes.Back(); i != nil && n >= 0; i = i.Prev() {
		scope := i.Value.(map[string]bool)
		if _, ok := scope[name.GetValue()]; ok {
			// fmt.Printf("debug: resolve: %s %#v %d\n", name.GetValue(), expr, r.scopes.Len()-1-n)
			r.interpreter.Resolve(expr, r.scopes.Len()-1-n)
			return
		}
		n--
	}
}

func (r *Resolver) resolveFunction(function *parser.Function) {
	defer r.newScope()()
	defer r.newFunctionState()()

	for _, param := range function.Params {
		r.decleare(param)
		r.define(param)
	}
	r.resolveStmt(function.Body)
}

func (r *Resolver) beginScope() {
	r.scopes.PushBack(make(map[string]bool))
}

func (r *Resolver) endScope() {
	v := r.scopes.Back()
	r.scopes.Remove(v)
}

func (r *Resolver) newScope() func() {
	r.beginScope()

	return r.endScope
}

func (r *Resolver) decleare(name *lexer.Token) {
	if r.scopes.Len() == 0 {
		return
	}
	scopeMap := r.scopes.Back().Value.(map[string]bool)
	if _, ok := scopeMap[name.GetValue()]; ok {
		panic(parser.NewParseError(name, "Already variable with this name in this scope."))
	}
	scopeMap[name.GetValue()] = false
}

func (r *Resolver) define(name *lexer.Token) {
	if r.scopes.Len() == 0 {
		return
	}
	scopeMap := r.scopes.Back().Value.(map[string]bool)
	scopeMap[name.GetValue()] = true
}

func (r *Resolver) VisitTernaryExpr(expr *parser.Ternary) (interface{}, error) {
	r.resolveExpr(expr.Condition)
	r.resolveExpr(expr.ThenExpr)
	r.resolveExpr(expr.ElseExpr)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *parser.Assign) (interface{}, error) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *parser.Binary) (interface{}, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *parser.Call) (interface{}, error) {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitArrayExpr(expr *parser.Array) (interface{}, error) {
	for _, element := range expr.Elements {
		r.resolveExpr(element)
	}
	return nil, nil
}

func (r *Resolver) VisitIndexExpr(expr *parser.Index) (interface{}, error) {
	r.resolveExpr(expr.Index)

	scope := r.scopes.Back()
	if scope == nil {
		return nil, nil
	}
	if v, ok := scope.Value.(map[string]bool)[expr.Name.GetValue()]; ok && r.scopes.Len() > 0 && !v {
		panic(NewRuntimeError("Can't read local variable in its own initializer."))
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *parser.Grouping) (interface{}, error) {
	r.resolveExpr(expr.Expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *parser.Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *parser.Logical) (interface{}, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *parser.Unary) (interface{}, error) {
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *parser.Variable) (interface{}, error) {
	scope := r.scopes.Back()
	if scope == nil {
		return nil, nil
	}
	if v, ok := scope.Value.(map[string]bool)[expr.Name.GetValue()]; ok && r.scopes.Len() > 0 && !v {
		panic(NewRuntimeError("Can't read local variable in its own initializer."))
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitLambdaExpr(expr *parser.Lambda) (interface{}, error) {
	r.resolveStmt(expr.Function)
	return nil, nil
}

func (r *Resolver) VisitBlockStmt(stmt *parser.Block) (interface{}, error) {
	defer r.newScope()()

	r.ResolveStmts(stmt.Statements)

	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt *parser.Expression) (interface{}, error) {
	r.resolveExpr(stmt.Expr)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *parser.Function) (interface{}, error) {
	r.decleare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt)
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *parser.If) (interface{}, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *parser.Print) (interface{}, error) {
	r.resolveExpr(stmt.Expr)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *parser.Return) (interface{}, error) {
	if !r.inFunction {
		panic(parser.NewParseError(stmt.Keyword, "Can't return from top-level code."))
	}
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *parser.Var) (interface{}, error) {
	for _, name := range stmt.Names {
		r.decleare(name)
	}

	for _, init := range stmt.Initializers {
		r.resolveExpr(init)
	}

	for _, name := range stmt.Names {
		r.define(name)

	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *parser.While) (interface{}, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil, nil
}

func (r *Resolver) VisitBreakStmt(stmt *parser.Break) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitContinueStmt(stmt *parser.Continue) (interface{}, error) {
	return nil, nil
}
