// Generated code, do not edit.

package parser

import "github.com/WAY29/LoxGo/lexer"

type Expr interface {
	Accept(v ExprVisitor) (interface{}, error)
}

type Ternary struct {
	Condition Expr
	ThenExpr  Expr
	ElseExpr  Expr
}

func NewTernary(condition Expr, thenexpr Expr, elseexpr Expr) *Ternary {
	return &Ternary{Condition: condition, ThenExpr: thenexpr, ElseExpr: elseexpr}
}
func (n *Ternary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitTernaryExpr(n)
}

type Assign struct {
	Name  *lexer.Token
	Value Expr
}

func NewAssign(name *lexer.Token, value Expr) *Assign {
	return &Assign{Name: name, Value: value}
}
func (n *Assign) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitAssignExpr(n)
}

type Binary struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

func NewBinary(left Expr, operator *lexer.Token, right Expr) *Binary {
	return &Binary{Left: left, Operator: operator, Right: right}
}
func (n *Binary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitBinaryExpr(n)
}

type Call struct {
	Callee      Expr
	Paren       *lexer.Token
	Arguments   []Expr
	ReturnValue interface{}
}

func NewCall(callee Expr, paren *lexer.Token, arguments []Expr, returnvalue interface{}) *Call {
	return &Call{Callee: callee, Paren: paren, Arguments: arguments, ReturnValue: returnvalue}
}
func (n *Call) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitCallExpr(n)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}
func (n *Grouping) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGroupingExpr(n)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{Value: value}
}
func (n *Literal) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLiteralExpr(n)
}

type Logical struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

func NewLogical(left Expr, operator *lexer.Token, right Expr) *Logical {
	return &Logical{Left: left, Operator: operator, Right: right}
}
func (n *Logical) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLogicalExpr(n)
}

type Unary struct {
	Operator *lexer.Token
	Right    Expr
	Prefix   bool
}

func NewUnary(operator *lexer.Token, right Expr, prefix bool) *Unary {
	return &Unary{Operator: operator, Right: right, Prefix: prefix}
}
func (n *Unary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitUnaryExpr(n)
}

type Variable struct {
	Name *lexer.Token
}

func NewVariable(name *lexer.Token) *Variable {
	return &Variable{Name: name}
}
func (n *Variable) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitVariableExpr(n)
}

type Lambda struct {
	Token    *lexer.Token
	Function Stmt
}

func NewLambda(token *lexer.Token, function Stmt) *Lambda {
	return &Lambda{Token: token, Function: function}
}
func (n *Lambda) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLambdaExpr(n)
}
