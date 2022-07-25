// package parser

// import (
// 	"bytes"
// 	"fmt"
// )

// type Printer struct{} // impl ExprVisitor

// func NewPrinter() *Printer {
// 	return new(Printer)
// }

// func (p *Printer) parenthesize(name string, exprs ...Expr) string {
// 	var buffer bytes.Buffer

// 	buffer.WriteString("(" + name)
// 	for _, expr := range exprs {
// 		buffer.WriteString(" ")
// 		output, _ := expr.Accept(p)
// 		buffer.WriteString(output.(string))
// 	}
// 	buffer.WriteString(")")

// 	return buffer.String()
// }

// func (p *Printer) Print(expr Expr) string {
// 	if expr == nil {
// 		return ""
// 	}
// 	output, _ := expr.Accept(p)
// 	return output.(string)
// }

// func (p *Printer) VisitUnaryExpr(expr *Unary) (interface{}, error) {
// 	return p.parenthesize(expr.Operator.GetValue(), expr.Right), nil
// }

// func (p *Printer) VisitBinaryExpr(expr *Binary) (interface{}, error) {
// 	return p.parenthesize(expr.Operator.GetValue(), expr.Left, expr.Right), nil
// }

// func (p *Printer) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
// 	return p.parenthesize("group", expr.Expression), nil
// }

// func (p *Printer) VisitLiteralExpr(expr *Literal) (interface{}, error) {
// 	if expr.Value == nil {
// 		return "nil", nil
// 	}
// 	return fmt.Sprintf("%v", expr.Value), nil
// }
package parser
