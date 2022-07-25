// Generated code, do not edit.

package parser

type ExprVisitor interface {
	VisitTernaryExpr(expr *Ternary) (interface{}, error)
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitCallExpr(expr *Call) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitLogicalExpr(expr *Logical) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
	VisitLambdaExpr(expr *Lambda) (interface{}, error)
}
