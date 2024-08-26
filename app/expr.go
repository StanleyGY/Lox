package main

type Stmt interface {
	Accept(v Visitor) error
}

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type InlineExprStmt struct {
	Child Expr
}

func (e *InlineExprStmt) Accept(v Visitor) error {
	return v.VisitInlineExprStmt(e)
}

type PrintStmt struct {
	Child Expr
}

func (e *PrintStmt) Accept(v Visitor) error {
	return v.VisitPrintStmt(e)
}

type VarDeclStmt struct {
	Name        *Token
	Initializer Expr
}

func (e *VarDeclStmt) Accept(v Visitor) error {
	return v.VisitVarDeclStmt(e)
}

type BinaryExpr struct {
	Operator *Token
	Left     Expr
	Right    Expr
}

func (e *BinaryExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinaryExpr(e)
}

type UnaryExpr struct {
	Operator *Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnaryExpr(e)
}

type GroupingExpr struct {
	Child Expr
}

func (e *GroupingExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitLiteralExpr(e)
}

type VariableExpr struct {
	Name *Token
}

func (e *VariableExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitVariableExpr(e)
}

type AssignExpr struct {
	Name  *Token
	Value Expr
}

func (e *AssignExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitAssignExpr(e)
}

type Visitor interface {
	VisitInlineExprStmt(stmt *InlineExprStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitVarDeclStmt(stmt *VarDeclStmt) error

	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitAssignExpr(expr *AssignExpr) (interface{}, error)
	VisitVariableExpr(expr *VariableExpr) (interface{}, error)
}
