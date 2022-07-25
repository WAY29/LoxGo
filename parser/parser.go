/*
grammer:

program        → declaration* EOF ;
declaration    → funDecl
			   | varDecl
               | statement ;
funDecl        → "fun" function ;
function       → IDENTIFIER "(" parameters? ")" block ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
statement      → exprStmt | forStmt | ifStmt | printStmt | returnStmt | whileStmt | block ;
forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
                 expression? ";"
                 expression? ")" statement ;
ifStmt         → "if" "(" expression ")" statement
               ( "else" statement )? ;
whileStmt      → "while" "(" expression ")" statement ;
exprStmt       → expression";" ;
printStmt      → "print" expression ";" ;
returnStmt     → "return" expression? ";"
block          → "{" declaration* "}" ;

expression     → assignment;
assignment     → IDENTIFIER "=" assignment
				 | ternary ;
ternary        -> logic_or ("?": ternary ":" ternary)? ;
logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ("++" | "--")? | call;
call           → primary ( "(" arguments? ")" )* ;
arguments      → expression ( "," expression )* ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" | IDENTIFIER | lambda;
lambda         → "fun" "(" parameters? ")" block;
*/
package parser

import (
	"fmt"

	"github.com/WAY29/LoxGo/lexer"
)

type Parser struct {
	tokens    []*lexer.Token
	tokensLen int
	current   int

	blockStmt Stmt
	whileStmt Stmt
}

func NewParaser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens:    tokens,
		tokensLen: len(tokens),
		current:   0,
		blockStmt: nil,
		whileStmt: nil,
	}
}

func (p *Parser) newBlockStmtState(newStmt Stmt) func() {
	oldStmt := p.blockStmt
	p.blockStmt = newStmt

	return func() {
		p.blockStmt = oldStmt
	}
}

func (p *Parser) newWhileStmtState(newStmt Stmt) func() {
	oldStmt := p.whileStmt
	p.whileStmt = newStmt

	return func() {
		p.whileStmt = oldStmt
	}
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().GetType() == tokenType
}

func (p *Parser) advance() *lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().GetType() == lexer.EOF
}

func (p *Parser) peek() *lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *lexer.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType lexer.TokenType, message string) *lexer.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(NewParseError(p.peek(), message))
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().GetType() == lexer.SEMICOLON {
			return
		}

		switch p.peek().GetType() {
		case lexer.CLASS:
			fallthrough
		case lexer.FUN:
			fallthrough
		case lexer.VAR:
			fallthrough
		case lexer.IF:
			fallthrough
		case lexer.WHILE:
			fallthrough
		case lexer.PRINT:
			fallthrough
		case lexer.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) Parse() (statements []Stmt) {
	statements = make([]Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() (stmt Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); ok {
				p.synchronize()
				stmt = nil
			} else {
				panic(r)
			}
		}
	}()

	if p.match(lexer.FUN) {
		return p.function("function", true)
	} else if p.match(lexer.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) function(kind string, hasName bool) (stmt Stmt) {
	var name *lexer.Token = nil
	var funcName string = ""
	if hasName {
		name = p.consume(lexer.IDENTIFIER, fmt.Sprintf("Excepted %s name.", kind))
		funcName = name.GetValue()
	}

	p.consume(lexer.LEFT_PAREN, fmt.Sprintf("Excepted '(' after fun %s.", funcName))
	parameters := make([]*lexer.Token, 0)
	if !p.check(lexer.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				panic(NewParseError(p.peek(), "Can't have more than 255 parameters."))
			}
			parameters = append(parameters, p.consume(lexer.IDENTIFIER, "Except parameter name."))

			if !p.match(lexer.COMMA) {
				break
			}
		}
	}
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(lexer.LEFT_BRACE, fmt.Sprintf("Excepted %s body.", kind))
	body := p.block()
	return NewFunction(name, parameters, body)
}

func (p *Parser) varDeclaration() Stmt {
	var (
		name         *lexer.Token
		initializer  Expr
		names        = make([]*lexer.Token, 0)
		initializers = make([]Expr, 0)
	)
	for {
		name = p.consume(lexer.IDENTIFIER, "Expect variable name.")
		if p.match(lexer.EQUAL) {
			initializer = p.expression()
		}
		names = append(names, name)
		initializers = append(initializers, initializer)
		if !p.match(lexer.COMMA) {
			break
		}
	}
	p.consume(lexer.SEMICOLON, "Expect ';' after variable declaration.")
	return NewVar(names, initializers)
}

func (p *Parser) statement() Stmt {
	if p.match(lexer.FOR) {
		return p.forStatement()
	} else if p.match(lexer.BREAK) {
		return p.breakStatement()
	} else if p.match(lexer.CONTINUE) {
		return p.continueStatement()
	} else if p.match(lexer.IF) {
		return p.ifStatement()
	} else if p.match(lexer.PRINT) {
		return p.printStatement()
	} else if p.match(lexer.RETURN) {
		return p.returnStatement()
	} else if p.match(lexer.WHILE) {
		return p.whileStatement()
	} else if p.match(lexer.LEFT_BRACE) {
		return p.block()
	}

	return p.expressionStatement()
}

func (p *Parser) continueStatement() Stmt {
	p.consume(lexer.SEMICOLON, "Expect ';' after continue.")

	return NewContinue(p.blockStmt, p.previous())
}

func (p *Parser) breakStatement() Stmt {
	p.consume(lexer.SEMICOLON, "Expect ';' after break.")

	return NewBreak(p.whileStmt, p.blockStmt, p.previous())
}

func (p *Parser) forStatement() Stmt {
	var (
		initializer, body    Stmt
		condition, increment Expr
		while                *While = NewWhile(nil, nil, false, nil)
	)
	defer p.newWhileStmtState(while)()

	p.consume(lexer.LEFT_PAREN, "Expect '(' after 'for'.")
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	if !p.check(lexer.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(lexer.SEMICOLON, "Expect ';' after loop condition.")

	if !p.check(lexer.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after for clauses.")

	body = p.statement()
	if increment != nil {
		body = NewBlock([]Stmt{body, NewExpression(increment)}, false, nil)
	}

	if condition == nil {
		condition = NewLiteral(true)
	}
	while.Condition = condition
	while.Body = body
	body = while

	if initializer != nil {
		body = NewBlock([]Stmt{initializer, body}, false, nil)
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	var (
		thenBranch, elseBranch Stmt
	)
	p.consume(lexer.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after 'if'.")

	thenBranch = p.statement()
	if p.match(lexer.ELSE) {
		elseBranch = p.statement()
	}

	return NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	return NewPrint(expr)
}

func (p *Parser) returnStatement() Stmt {
	var value Expr = nil
	keyword := p.previous()
	if !p.check(lexer.SEMICOLON) {
		value = p.expression()
	}
	p.consume(lexer.SEMICOLON, "Expect ';' after return value.")
	return NewReturn(keyword, value)
}

func (p *Parser) whileStatement() Stmt {
	var while *While = NewWhile(nil, nil, false, nil)
	defer p.newWhileStmtState(while)()

	p.consume(lexer.LEFT_PAREN, "Expect '(' after 'while'.")
	while.Condition = p.expression()
	p.consume(lexer.RIGHT_PAREN, "Expect ')' after 'while'.")
	while.Body = p.statement()
	return while
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	return NewExpression(expr)
}

func (p *Parser) block() Stmt {
	block := NewBlock(nil, false, nil)
	stmts := make([]Stmt, 0)
	defer p.newBlockStmtState(block)()

	for !p.isAtEnd() && !p.check(lexer.RIGHT_BRACE) {
		stmts = append(stmts, p.declaration())
	}
	p.consume(lexer.RIGHT_BRACE, "Expect '}' after block.")

	block.Statements = stmts
	return block
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.ternary()

	if p.match(lexer.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if variable, ok := expr.(*Variable); ok {
			name := variable.Name
			return NewAssign(name, value)
		}
		panic(NewParseError(equals, "Invalid assignment target."))
	}

	return expr
}

func (p *Parser) ternary() Expr {
	expr := p.or()

	if p.match(lexer.QUESTION) {
		thenExpr := p.ternary()
		p.consume(lexer.COLON, "Expect ':' after ternary '?'.")
		elseExpr := p.ternary()
		expr = NewTernary(expr, thenExpr, elseExpr)
	}

	return expr
}

func (p *Parser) or() Expr {
	var (
		expr, right Expr
		operator    *lexer.Token
	)
	expr = p.and()

	for p.match(lexer.OR) {
		operator = p.previous()
		right = p.and()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) and() Expr {
	var (
		expr, right Expr
		operator    *lexer.Token
	)
	expr = p.equality()

	for p.match(lexer.AND) {
		operator = p.previous()
		right = p.equality()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(lexer.EQUAL_EQUAL, lexer.BANG_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(lexer.BANG, lexer.MINUS, lexer.PLUSPLUS, lexer.MINUSMINUS) {
		return NewUnary(p.previous(), p.unary(), true)
	} else {
		expr := p.primary()
		if p.match(lexer.PLUSPLUS, lexer.MINUSMINUS) {
			return NewUnary(p.previous(), expr, false)
		} else if p.match(lexer.LEFT_PAREN) {
			return p.call(expr)
		} else {
			return expr
		}
	}
}

func (p *Parser) call(expr Expr) Expr {
	for {
		expr = p.finishCall(expr)
		if !p.match(lexer.LEFT_PAREN) {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0)
	if !p.check(lexer.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				panic(NewParseError(p.peek(), "Function can't have more than 255 arguments."))
			}
			arguments = append(arguments, p.expression())
			if !p.match(lexer.COMMA) {
				break
			}
		}
	}
	paren := p.consume(lexer.RIGHT_PAREN, "Expect ')' after arguments.")
	return NewCall(callee, paren, arguments, nil)
}

func (p *Parser) primary() Expr {
	if p.match(lexer.FALSE) {
		return NewLiteral(false)
	} else if p.match(lexer.TRUE) {
		return NewLiteral(true)
	} else if p.match(lexer.NIL) {
		return NewLiteral(nil)
	} else if p.match(lexer.NUMBER, lexer.STRING) {
		return NewLiteral(p.previous().GetLiteral())
	} else if p.match(lexer.IDENTIFIER) {
		return NewVariable(p.previous())
	} else if p.match(lexer.LEFT_PAREN) {
		expr := p.expression()
		p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	} else if p.match(lexer.FUN) {
		function := p.function("function", false)
		return NewLambda(p.previous(), function)
	}

	panic(NewParseError(p.peek(), "Except expression."))
}
