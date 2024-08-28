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

type FuncDeclStmt struct {
	Name   *Token
	Params []*Token
	Body   Stmt
}

func (e *FuncDeclStmt) Accept(v Visitor) error {
	return v.VisitFunDeclStmt(e)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *IfStmt) Accept(v Visitor) error {
	return v.VisitIfStmt(e)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (e *WhileStmt) Accept(v Visitor) error {
	return v.VisitWhileStmt(e)
}

type ReturnStmt struct {
	// TODO: Keyword *Token for error reporting
	Value Expr
}

func (e *ReturnStmt) Accept(v Visitor) error {
	return v.VisitReturnStmt(e)
}

type BlockStmt struct {
	Stmts []Stmt
}

func (e *BlockStmt) Accept(v Visitor) error {
	return v.VisitBlockStmt(e)
}

type BinaryExpr struct {
	Operator *Token
	Left     Expr
	Right    Expr
}

func (e *BinaryExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinaryExpr(e)
}

// LogicExpr differs from BinaryExpr in that it supports short-circuit
// evaluation on its operands
type LogicExpr struct {
	Operator *Token
	Left     Expr
	Right    Expr
}

func (e *LogicExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitLogicalExpr(e)
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

type CallExpr struct {
	Callee    Expr
	Arguments []Expr
}

func (e *CallExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitCallExpr(e)
}

type Visitor interface {
	VisitVarDeclStmt(stmt *VarDeclStmt) error
	VisitFunDeclStmt(stmt *FuncDeclStmt) error

	VisitInlineExprStmt(stmt *InlineExprStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
	VisitReturnStmt(stmt *ReturnStmt) error

	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitLogicalExpr(expr *LogicExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitAssignExpr(expr *AssignExpr) (interface{}, error)
	VisitCallExpr(expr *CallExpr) (interface{}, error)
	VisitVariableExpr(expr *VariableExpr) (interface{}, error)
}
