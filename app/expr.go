package main

type Expr interface {
	Accept(v Visitor) (interface{}, error)
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

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
}
