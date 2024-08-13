package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("Test binary expr", func(t *testing.T) {
		parser := &RDParser{}
		expr, err := parser.Parse([]*Token{
			{Type: Number, Literal: 1},
			{Type: Plus, Lexeme: "+"},
			{Type: Number, Literal: 1},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrint(expr)
		assert.Equal(t, "(+ 1 1)", res)
	})

	t.Run("Test unary expr", func(t *testing.T) {
		parser := &RDParser{}
		expr, err := parser.Parse([]*Token{
			{Type: Bang, Lexeme: "!"},
			{Type: Bang, Lexeme: "!"},
			{Type: Number, Literal: 1},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrint(expr)
		assert.Equal(t, "(! (! 1))", res)
	})

	t.Run("Test comparison expr", func(t *testing.T) {
		parser := &RDParser{}
		expr, err := parser.Parse([]*Token{
			{Type: String, Literal: "3"},
			{Type: BangEqual, Lexeme: "!="},
			{Type: Number, Literal: 1},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrint(expr)
		assert.Equal(t, "(!= \"3\" 1)", res)
	})

	t.Run("Test factor expr", func(t *testing.T) {
		parser := &RDParser{}
		expr, err := parser.Parse([]*Token{
			{Type: Number, Literal: 3},
			{Type: Star, Lexeme: "*"},
			{Type: Number, Literal: 1.6},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrint(expr)
		assert.Equal(t, "(* 3 1.6)", res)
	})

	t.Run("Test primary expr", func(t *testing.T) {
		parser := &RDParser{}
		printer := &AstPrinter{}

		expr, _ := parser.Parse([]*Token{
			{Type: False, Lexeme: "false"},
		})
		res := printer.PrettyPrint(expr)
		assert.Equal(t, "false", res)

		expr, _ = parser.Parse([]*Token{
			{Type: True, Lexeme: "true"},
		})
		res = printer.PrettyPrint(expr)
		assert.Equal(t, "true", res)

		expr, _ = parser.Parse([]*Token{
			{Type: Nil, Lexeme: "nil"},
		})
		res = printer.PrettyPrint(expr)
		assert.Equal(t, "nil", res)
	})

	t.Run("Test error expr with no matching rule", func(t *testing.T) {
		parser := &RDParser{}
		_, err := parser.Parse([]*Token{
			{Type: Number, Literal: 3},
			{Type: Number, Literal: 1.6},
		})
		assert.Error(t, err)
	})
}
