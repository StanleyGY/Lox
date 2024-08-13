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

func (p *AstPrinter) PrettyPrint(expr Expr) string {
	p.buf.Reset()
	expr.Accept(p)
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

func (p *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) {
	p.parenthesis(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) {
	p.parenthesis(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *GroupingExpr) {
	p.parenthesis("Group", expr.Child)
}

func (p *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) {
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
}
