// Generated code, do not edit.

package parser

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) (interface{}, error)
	VisitClassStmt(stmt *Class) (interface{}, error)
	VisitExpressionStmt(stmt *Expression) (interface{}, error)
	VisitFunctionStmt(stmt *Function) (interface{}, error)
	VisitIfStmt(stmt *If) (interface{}, error)
	VisitPrintStmt(stmt *Print) (interface{}, error)
	VisitReturnStmt(stmt *Return) (interface{}, error)
	VisitVarStmt(stmt *Var) (interface{}, error)
	VisitWhileStmt(stmt *While) (interface{}, error)
	VisitBreakStmt(stmt *Break) (interface{}, error)
	VisitContinueStmt(stmt *Continue) (interface{}, error)
}
