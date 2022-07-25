// package parser

// import (
// 	"testing"

// 	"github.com/WAY29/LoxGo/lexer"
// )

// func TestPrinter(t *testing.T) {
// 	printer := new(Printer)

// 	testCases := map[string]Expr{
// 		"(+ 1 2)": NewBinary(
// 			NewLiteral(1),
// 			lexer.NewToken(lexer.PLUS, "+", nil, 1),
// 			NewLiteral(2),
// 		),
// 		"(* (- 123) (group 45.67))": NewBinary(
// 			NewUnary(
// 				lexer.NewToken(lexer.MINUS, "-", nil, 1),
// 				NewLiteral(123),
// 			),
// 			lexer.NewToken(lexer.STAR, "*", nil, 1),
// 			NewGrouping(NewLiteral(45.67)),
// 		),
// 	}

// 	for expected, expr := range testCases {
// 		if result := printer.Print(expr); result != expected {
// 			t.Fatalf("expected %s, but got %s", expected, result)
// 		}
// 	}

// }
package parser
