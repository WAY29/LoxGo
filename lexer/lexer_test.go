package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	testCases := []string{
		`print "hello world"`,
		`a = 5`,
		`b = 6// test comment`,
		`c = 8 / 2`,
	}
	for _, testCase := range testCases {
		l := NewLexer(strings.NewReader(testCase))
		l.ScanTokens()
		fmt.Println(l.tokens)
	}
}
