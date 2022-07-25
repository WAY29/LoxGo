package lexer

import "fmt"

type Token struct {
	_type   TokenType
	value   string
	literal interface{}
	line    uint16
}

func (t *Token) GetType() TokenType {
	return t._type
}

func (t *Token) GetValue() string {
	return t.value
}

func (t *Token) GetLiteral() interface{} {
	return t.literal
}

func (t *Token) GetLine() uint16 {
	return t.line
}

func (t *Token) String() string {
	return fmt.Sprintf("%s:%s", t._type, t.value)
}

func NewToken(_type TokenType, value string, literal interface{}, line uint16) *Token {
	return &Token{
		_type:   _type,
		value:   value,
		literal: literal,
		line:    line,
	}
}
