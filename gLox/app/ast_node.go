package main

type Stmt interface {
	Accept(v StmtVisitor) error
}

type Expr interface {
	Accept(v ExprVisitor) (interface{}, error)
}

type InlineExprStmt struct {
	Child Expr
}

func (e *InlineExprStmt) Accept(v StmtVisitor) error {
	return v.VisitInlineExprStmt(e)
}

type PrintStmt struct {
	Child Expr
}

func (e *PrintStmt) Accept(v StmtVisitor) error {
	return v.VisitPrintStmt(e)
}

type VarDeclStmt struct {
	Name        *Token
	Initializer Expr
}

func (e *VarDeclStmt) Accept(v StmtVisitor) error {
	return v.VisitVarDeclStmt(e)
}

type FuncDeclStmt struct {
	Name   *Token
	Params []*Token
	Body   Stmt
}

func (e *FuncDeclStmt) Accept(v StmtVisitor) error {
	return v.VisitFunDeclStmt(e)
}

type ClassDeclStmt struct {
	Name       *Token
	SuperClass *VariableExpr
	Methods    []*FuncDeclStmt
}

func (e *ClassDeclStmt) Accept(v StmtVisitor) error {
	return v.VisitClassDeclStmt(e)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *IfStmt) Accept(v StmtVisitor) error {
	return v.VisitIfStmt(e)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (e *WhileStmt) Accept(v StmtVisitor) error {
	return v.VisitWhileStmt(e)
}

type ReturnStmt struct {
	Value Expr
}

func (e *ReturnStmt) Accept(v StmtVisitor) error {
	return v.VisitReturnStmt(e)
}

type BreakStmt struct {
}

func (e *BreakStmt) Accept(v StmtVisitor) error {
	return v.VisitBreakStmt(e)
}

type BlockStmt struct {
	Stmts []Stmt
}

func (e *BlockStmt) Accept(v StmtVisitor) error {
	return v.VisitBlockStmt(e)
}

type BinaryExpr struct {
	Operator *Token
	Left     Expr
	Right    Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitBinaryExpr(e)
}

// LogicExpr differs from BinaryExpr in that it supports short-circuit
// evaluation on its operands
type LogicExpr struct {
	Operator *Token
	Left     Expr
	Right    Expr
}

func (e *LogicExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLogicalExpr(e)
}

type UnaryExpr struct {
	Operator *Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitUnaryExpr(e)
}

type GroupingExpr struct {
	Child Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLiteralExpr(e)
}

type VariableExpr struct {
	Name *Token
}

func (e *VariableExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitVariableExpr(e)
}

type AssignExpr struct {
	Name  *Token
	Value Expr
}

func (e *AssignExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitAssignExpr(e)
}

type CallExpr struct {
	Callee    Expr
	Arguments []Expr
}

func (e *CallExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitCallExpr(e)
}

type GetPropertyExpr struct {
	Object   Expr
	Property *Token
}

func (e *GetPropertyExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGetPropertyExpr(e)
}

type SetPropertyExpr struct {
	Object   Expr
	Property *Token
	Value    Expr
}

func (e *SetPropertyExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitSetPropertyExpr(e)
}

type ThisExpr struct{}

func (e *ThisExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitThisExpr(e)
}

type SuperExpr struct {
	Property *Token
}

func (e *SuperExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitSuperExpr(e)
}

type StmtVisitor interface {
	VisitVarDeclStmt(stmt *VarDeclStmt) error
	VisitFunDeclStmt(stmt *FuncDeclStmt) error
	VisitClassDeclStmt(stmt *ClassDeclStmt) error
	VisitInlineExprStmt(stmt *InlineExprStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
	VisitReturnStmt(stmt *ReturnStmt) error
	VisitBreakStmt(stmt *BreakStmt) error
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitLogicalExpr(expr *LogicExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitAssignExpr(expr *AssignExpr) (interface{}, error)
	VisitCallExpr(expr *CallExpr) (interface{}, error)
	VisitGetPropertyExpr(expr *GetPropertyExpr) (interface{}, error)
	VisitSetPropertyExpr(expr *SetPropertyExpr) (interface{}, error)
	VisitVariableExpr(expr *VariableExpr) (interface{}, error)
	VisitThisExpr(expr *ThisExpr) (interface{}, error)
	VisitSuperExpr(expr *SuperExpr) (interface{}, error)
}
