package main

type Expr interface {
	Accept(v Visitor)
}

type BinaryExpr struct {
	Operator *Token
	Left     Expr
	Right    Expr
}

func (e *BinaryExpr) Accept(v Visitor) {
	v.VisitBinaryExpr(e)
}

type UnaryExpr struct {
	Operator *Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v Visitor) {
	v.VisitUnaryExpr(e)
}

type GroupingExpr struct {
	Child Expr
}

func (e *GroupingExpr) Accept(v Visitor) {
	v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(v Visitor) {
	v.VisitLiteralExpr(e)
}

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr)
	VisitUnaryExpr(expr *UnaryExpr)
	VisitGroupingExpr(expr *GroupingExpr)
	VisitLiteralExpr(expr *LiteralExpr)
}
