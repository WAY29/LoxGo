// Generated code, do not edit.

package parser

import "github.com/WAY29/LoxGo/lexer"

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

type Block struct {
	Statements []Stmt
	Stop       bool
	Parent     *Block
}

func NewBlock(statements []Stmt, stop bool, parent *Block) *Block {
	return &Block{Statements: statements, Stop: stop, Parent: parent}
}
func (n *Block) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBlockStmt(n)
}

type Expression struct {
	Expr Expr
}

func NewExpression(expr Expr) *Expression {
	return &Expression{Expr: expr}
}
func (n *Expression) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitExpressionStmt(n)
}

type Function struct {
	Name   *lexer.Token
	Params []*lexer.Token
	Body   Stmt
}

func NewFunction(name *lexer.Token, params []*lexer.Token, body Stmt) *Function {
	return &Function{Name: name, Params: params, Body: body}
}
func (n *Function) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitFunctionStmt(n)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenbranch Stmt, elsebranch Stmt) *If {
	return &If{Condition: condition, ThenBranch: thenbranch, ElseBranch: elsebranch}
}
func (n *If) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitIfStmt(n)
}

type Print struct {
	Expr Expr
}

func NewPrint(expr Expr) *Print {
	return &Print{Expr: expr}
}
func (n *Print) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitPrintStmt(n)
}

type Return struct {
	Keyword     *lexer.Token
	Initializer Expr
}

func NewReturn(keyword *lexer.Token, initializer Expr) *Return {
	return &Return{Keyword: keyword, Initializer: initializer}
}
func (n *Return) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitReturnStmt(n)
}

type Var struct {
	Names        []*lexer.Token
	Initializers []Expr
}

func NewVar(names []*lexer.Token, initializers []Expr) *Var {
	return &Var{Names: names, Initializers: initializers}
}
func (n *Var) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitVarStmt(n)
}

type While struct {
	Condition Expr
	Body      Stmt
	Stop      bool
	Parent    *While
}

func NewWhile(condition Expr, body Stmt, stop bool, parent *While) *While {
	return &While{Condition: condition, Body: body, Stop: stop, Parent: parent}
}
func (n *While) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitWhileStmt(n)
}

type Break struct {
	Parent Stmt
	Block  Stmt
	Token  *lexer.Token
}

func NewBreak(parent Stmt, block Stmt, token *lexer.Token) *Break {
	return &Break{Parent: parent, Block: block, Token: token}
}
func (n *Break) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBreakStmt(n)
}

type Continue struct {
	Block Stmt
	Token *lexer.Token
}

func NewContinue(block Stmt, token *lexer.Token) *Continue {
	return &Continue{Block: block, Token: token}
}
func (n *Continue) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitContinueStmt(n)
}
