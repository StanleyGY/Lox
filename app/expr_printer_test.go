package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExprPrinter(t *testing.T) {
	expr := &BinaryExpr{
		Operator: &Token{
			Type:   Plus,
			Lexeme: "+",
		},
		Left: &LiteralExpr{
			Value: 1,
		},
		Right: &LiteralExpr{
			Value: 2,
		},
	}
	printer := &AstPrinter{}
	res := printer.PrettyPrint(expr)
	assert.Equal(t, "(+ 1 2)", res)
}
