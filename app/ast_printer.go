package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type AstPrinter struct {
	buf bytes.Buffer
}

func (p *AstPrinter) PrettyPrintExpr(expr Expr) string {
	p.buf.Reset()
	expr.Accept(p)
	return p.buf.String()
}

func (p *AstPrinter) PrettyPrintStmt(stmt Stmt) string {
	p.buf.Reset()
	stmt.Accept(p)
	return p.buf.String()
}

func (p *AstPrinter) parenthesis(name string, exprs ...Expr) {
	p.buf.WriteString("(")
	p.buf.WriteString(name)
	p.buf.WriteString(" ")
	for idx, expr := range exprs {
		expr.Accept(p)
		if idx < len(exprs)-1 {
			p.buf.WriteString(" ")
		}
	}
	p.buf.WriteString(")")
}

func (p *AstPrinter) VisitInlineExprStmt(stmt *InlineExprStmt) error {
	stmt.Child.Accept(p)
	return nil
}

func (p *AstPrinter) VisitPrintStmt(stmt *PrintStmt) error {
	p.parenthesis("print", stmt.Child)
	return nil
}

func (p *AstPrinter) VisitVarDeclStmt(stmt *VarDeclStmt) error {
	p.parenthesis("let", stmt.Initializer)
	return nil
}

func (p *AstPrinter) VisitFunDeclStmt(stmt *FuncDeclStmt) error {
	// TODO:
	return nil
}

func (p *AstPrinter) VisitIfStmt(stmt *IfStmt) error {
	p.buf.WriteString("If")
	stmt.Condition.Accept(p)
	p.buf.WriteString("Then")
	stmt.ThenBranch.Accept(p)
	p.buf.WriteString("Else")
	stmt.ElseBranch.Accept(p)
	return nil
}

func (p *AstPrinter) VisitWhileStmt(stmt *WhileStmt) error {
	// TODO: add ast printer
	return nil
}

func (p *AstPrinter) VisitReturnStmt(stmt *ReturnStmt) error {
	// TODO: add ast printer
	return nil
}

func (p *AstPrinter) VisitBlockStmt(stmt *BlockStmt) error {
	for _, s := range stmt.Stmts {
		s.Accept(p)
	}
	return nil
}

func (p *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	p.parenthesis(expr.Operator.Lexeme, expr.Left, expr.Right)
	return nil, nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	p.parenthesis(expr.Operator.Lexeme, expr.Right)
	return nil, nil
}

func (p *AstPrinter) VisitLogicalExpr(expr *LogicExpr) (interface{}, error) {
	p.parenthesis(expr.Operator.Lexeme, expr.Left, expr.Right)
	return nil, nil
}

func (p *AstPrinter) VisitAssignExpr(expr *AssignExpr) (interface{}, error) {
	p.parenthesis("let", &LiteralExpr{expr.Name}, expr.Value)
	return nil, nil
}

func (p *AstPrinter) VisitCallExpr(expr *CallExpr) (interface{}, error) {
	return nil, nil
}

func (p *AstPrinter) VisitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	p.parenthesis("Group", expr.Child)
	return nil, nil
}

func (p *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	switch expr.Value.(type) {
	case string:
		p.buf.WriteString(fmt.Sprintf("\"%s\"", expr.Value.(string)))
	case float64:
		p.buf.WriteString(strconv.FormatFloat(expr.Value.(float64), 'f', -1, 64))
	case int:
		p.buf.WriteString(strconv.Itoa(expr.Value.(int)))
	case bool:
		p.buf.WriteString(strconv.FormatBool(expr.Value.(bool)))
	default:
		reflVal := reflect.ValueOf(expr.Value)
		if !reflVal.IsValid() {
			p.buf.WriteString("nil")
		}
	}
	return nil, nil
}

func (p *AstPrinter) VisitVariableExpr(expr *VariableExpr) (interface{}, error) {
	p.buf.WriteString(expr.Name.Lexeme)
	return nil, nil
}
