// Generated code, do not edit.

package parser

type ExprVisitor interface {
	VisitTernaryExpr(expr *Ternary) (interface{}, error)
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitCallExpr(expr *Call) (interface{}, error)
	VisitGetExpr(expr *Get) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitLogicalExpr(expr *Logical) (interface{}, error)
	VisitSetExpr(expr *Set) (interface{}, error)
	VisitThisExpr(expr *This) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
	VisitArrayExpr(expr *Array) (interface{}, error)
	VisitIndexExpr(expr *Index) (interface{}, error)
	VisitLambdaExpr(expr *Lambda) (interface{}, error)
}
