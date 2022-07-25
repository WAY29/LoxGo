package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

var (
	lexerPool = &sync.Pool{
		New: func() interface{} {
			return new(Lexer)
		},
	}
)

type Lexer struct {
	tokens []*Token

	line   uint16
	buffer *bufio.Reader
	err    error
	end    bool
}

func NewLexer(reader io.Reader) *Lexer {
	l := lexerPool.Get().(*Lexer)
	l.buffer = bufio.NewReader(reader)
	l.line = 1
	l.end = false

	return l
}

func (l *Lexer) GetTokens() []*Token {
	return l.tokens
}

func (l *Lexer) GetError() error {
	return l.err
}

func (l *Lexer) peekOne() (rune, error) {
	chars, err := l.buffer.Peek(1)
	if err != nil {
		return 0, err
	}

	return rune(chars[0]), err
}

func (l *Lexer) peek(n int) ([]byte, error) {
	chars, err := l.buffer.Peek(n)
	if err != nil {
		return nil, err
	}

	return chars, err
}

func (l *Lexer) read(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := l.buffer.Read(bytes)
	if err == io.EOF {
		err = nil
		l.end = true
	}

	return bytes, err
}

func (l *Lexer) unreadByte() error {
	if l.end == true {
		return nil
	}
	return l.buffer.UnreadByte()
}

func (l *Lexer) readUntil(delim byte) ([]byte, error) {
	if l.end == true {
		return nil, nil
	}

	bytes, err := l.buffer.ReadBytes(delim)
	if err == io.EOF {
		err = nil
		l.end = true
	}

	if bytes[len(bytes)-1] == delim {
		l.buffer.UnreadByte()
		return bytes[:len(bytes)-1], err
	}
	return bytes, err
}

func (l *Lexer) readRune() (rune, error) {
	ch, _, err := l.buffer.ReadRune()
	if err == io.EOF {
		err = nil
		l.end = true
	}
	return ch, err
}

func (l *Lexer) unreadRune() error {
	if l.end == true {
		return nil
	}
	return l.buffer.UnreadRune()
}

func (l *Lexer) match(c rune) (bool, error) {
	pc, err := l.peekOne()
	if l.end == true {
		return false, nil
	}
	ok := pc == c

	if ok {
		_, err = l.readRune()
		if err != nil {
			return false, err
		}
	}
	return ok, nil
}

func (l *Lexer) scanString() error {
	var (
		c       rune
		err     error = nil
		builder       = &strings.Builder{}
	)

	for {
		c, err = l.readRune()
		if l.end {
			err = fmt.Errorf("unterminated string")
			break
		}
		if err != nil {
			break
		}

		if c == '"' {
			l.addToken(STRING, builder.String(), builder.String())
			break
		} else if c == '\n' {
			err = l.unreadByte()
			if err != nil {
				break
			}
			err = fmt.Errorf("unexpect '\\n' in string")
		} else {
			_, err = builder.WriteRune(c)
			if err != nil {
				break
			}
		}
	}
	return err
}

func (l *Lexer) scanNumber(hc rune) error {
	var (
		c            rune
		err          error = nil
		builder            = &strings.Builder{}
		hasScanPoint bool  = false
		v            interface{}
		vString      string
	)
	builder.WriteRune(hc)

	for {
		c, err = l.readRune()
		if err != nil {
			break
		}

		if isDigit(c) {
			_, err = builder.WriteRune(c)
			if err != nil {
				break
			}
		} else {
			if hasScanPoint || c != '.' {
				err = l.unreadByte()
				vString = builder.String()
				if strings.Contains(vString, ".") {
					v, err = strconv.ParseFloat(vString, 64)
				} else {
					v, err = strconv.Atoi(vString)
				}
				l.addToken(NUMBER, vString, v)
				break
			} else if c == '.' {
				hasScanPoint = true
				_, err = builder.WriteRune(c)
				if err != nil {
					break
				}
				continue
			}

		}
	}
	return err
}

func (l *Lexer) scanIdentifier(hc rune) error {
	var (
		c         rune
		err       error = nil
		builder         = &strings.Builder{}
		value     string
		tokenType TokenType
		ok        bool
	)
	builder.WriteRune(hc)

	for {
		c, err = l.readRune()
		if err != nil {
			return err
		}

		if isAlphaNumeric(c) {
			_, err = builder.WriteRune(c)
			if err != nil {
				return err
			}
		} else {
			value = builder.String()
			if tokenType, ok = KEYWORDS[value]; ok {
				l.addToken(tokenType, value)
			} else {
				l.addToken(IDENTIFIER, value)
			}

			err = l.unreadRune()
			if err != nil {
				return err
			}
			break
		}
	}

	return err
}

func (l *Lexer) addToken(_type TokenType, value string, literals ...interface{}) {
	var literal interface{}
	if len(literals) > 0 {
		literal = literals[0]
	}

	l.tokens = append(l.tokens, NewToken(_type, value, literal, l.line))
}

func (l *Lexer) scanToken() error {
	var (
		c   rune
		err error

		ok bool
	)

	c, err = l.readRune()
	if err != nil {
		return err
	}
	if l.end {
		return nil
	}

	switch c {
	case '(':
		l.addToken(LEFT_PAREN, "(")
	case ')':
		l.addToken(RIGHT_PAREN, ")")
	case '{':
		l.addToken(LEFT_BRACE, "{")
	case '}':
		l.addToken(RIGHT_BRACE, "}")
	case ',':
		l.addToken(COMMA, ",")
	case '.':
		l.addToken(DOT, ".")
	case '?':
		l.addToken(QUESTION, "?")
	case ':':
		l.addToken(COLON, ":")
	case '-':
		if ok, err = l.match('-'); ok {
			l.addToken(MINUSMINUS, "--")
		} else {
			l.addToken(MINUS, "-")
		}
	case '+':
		if ok, err = l.match('+'); ok {
			l.addToken(PLUSPLUS, "++")
		} else {
			l.addToken(PLUS, "+")
		}
	case ';':
		l.addToken(SEMICOLON, ";")
	case '*':
		l.addToken(STAR, "*")
	case '/':
		if ok, err = l.match('/'); ok {
			_, err = l.readUntil('\n')
		} else {
			l.addToken(SLASH, "/")
		}
	case '!':
		if ok, err = l.match('='); ok {
			l.addToken(BANG_EQUAL, "!=")
		} else {
			l.addToken(BANG, "!")
		}
	case '=':
		if ok, err = l.match('='); ok {
			l.addToken(EQUAL_EQUAL, "==")
		} else {
			l.addToken(EQUAL, "=")
		}
	case '<':
		if ok, err = l.match('='); ok {
			l.addToken(LESS_EQUAL, "<=")
		} else {
			l.addToken(LESS, "<")
		}
	case '>':
		if ok, err = l.match('='); ok {
			l.addToken(GREATER_EQUAL, ">=")
		} else {
			l.addToken(GREATER, ">")
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		l.line++
	case '"': // string
		err = l.scanString()
	default:
		if isDigit(c) {
			err = l.scanNumber(c)
		} else if isAlpha(c) {
			err = l.scanIdentifier(c)
		} else if c != '\u0000' {
			err = fmt.Errorf("line %d error: unexpected character: %c", l.line, c)
		}
	}

	return err
}

func (l *Lexer) ScanTokens() {
	var (
		err error
	)
	for !l.end {
		err = l.scanToken()
		if err != nil {
			break
		}
	}

	l.addToken(EOF, "", nil)

	if err == io.EOF {
		return
	} else if err != nil {
		l.err = fmt.Errorf("lexer error: %v", err)
		return
	}
}
