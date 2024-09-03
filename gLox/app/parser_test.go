package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("Test single stmt with binary expr", func(t *testing.T) {
		parser := &RDParser{}
		stmts, err := parser.Parse([]*Token{
			{Type: Number, Literal: 1},
			{Type: Plus, Lexeme: "+"},
			{Type: Number, Literal: 1},
			{Type: SemiColon},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrintStmt(stmts[0])
		assert.Equal(t, "(+ 1 1)", res)
	})

	t.Run("Test single stmt with unary expr", func(t *testing.T) {
		parser := &RDParser{}
		stmts, err := parser.Parse([]*Token{
			{Type: Bang, Lexeme: "!"},
			{Type: Bang, Lexeme: "!"},
			{Type: Number, Literal: 1},
			{Type: SemiColon},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrintStmt(stmts[0])
		assert.Equal(t, "(! (! 1))", res)
	})

	t.Run("Test single stmt with comparison expr", func(t *testing.T) {
		parser := &RDParser{}
		stmts, err := parser.Parse([]*Token{
			{Type: String, Literal: "3"},
			{Type: BangEqual, Lexeme: "!="},
			{Type: Number, Literal: 1},
			{Type: SemiColon},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrintStmt(stmts[0])
		assert.Equal(t, "(!= \"3\" 1)", res)
	})

	t.Run("Test single stmt with factor expr", func(t *testing.T) {
		parser := &RDParser{}
		stmts, err := parser.Parse([]*Token{
			{Type: Number, Literal: 3},
			{Type: Star, Lexeme: "*"},
			{Type: Number, Literal: 1.6},
			{Type: SemiColon},
		})
		assert.NoError(t, err)

		printer := &AstPrinter{}
		res := printer.PrettyPrintStmt(stmts[0])
		assert.Equal(t, "(* 3 1.6)", res)
	})

	t.Run("Test single stmt with primary expr", func(t *testing.T) {
		parser := &RDParser{}
		printer := &AstPrinter{}

		stmts, _ := parser.Parse([]*Token{
			{Type: False, Lexeme: "false"},
			{Type: SemiColon},
		})
		res := printer.PrettyPrintStmt(stmts[0])
		assert.Equal(t, "false", res)

		stmts, _ = parser.Parse([]*Token{
			{Type: True, Lexeme: "true"},
			{Type: SemiColon},
		})
		res = printer.PrettyPrintStmt(stmts[0])
		assert.Equal(t, "true", res)

		stmts, _ = parser.Parse([]*Token{
			{Type: Nil, Lexeme: "nil"},
			{Type: SemiColon},
		})
		res = printer.PrettyPrintStmt(stmts[0])
		assert.Equal(t, "nil", res)
	})

	t.Run("Test multiple stmts", func(t *testing.T) {
		parser := &RDParser{}
		stmts, _ := parser.Parse([]*Token{
			{Type: Number, Literal: 3},
			{Type: SemiColon},
			{Type: Number, Literal: 4},
			{Type: SemiColon},
		})
		assert.Equal(t, len(stmts), 2)
	})

	t.Run("Test single stmt with error expr not having matched rule", func(t *testing.T) {
		parser := &RDParser{}
		_, err := parser.Parse([]*Token{
			{Type: Number, Literal: 3},
			{Type: Number, Literal: 1.6},
			{Type: SemiColon},
		})
		assert.Error(t, err)
	})

	t.Run("Test single stmt with missing semicolon", func(t *testing.T) {
		parser := &RDParser{}
		_, err := parser.Parse([]*Token{
			{Type: Number, Literal: 3},
		})
		assert.Error(t, err)
	})

}
