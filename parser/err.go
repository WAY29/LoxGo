package parser

import (
	"fmt"

	"github.com/WAY29/LoxGo/lexer"
)

type ParseError struct {
	token *lexer.Token

	extraMsg string
}

func NewParseError(token *lexer.Token, msg string) *ParseError {
	return &ParseError{
		token:    token,
		extraMsg: msg,
	}
}

func (e *ParseError) Error() string {
	var (
		where string
		token = e.token
	)

	if token.GetType() == lexer.EOF {
		where = " at end"
	} else {
		where = fmt.Sprintf("at '%s'", token.GetValue())
	}

	return fmt.Sprintf("Parse error in line %d: Error %s: %s", token.GetLine(), where, e.extraMsg)
}
