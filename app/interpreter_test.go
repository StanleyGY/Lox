package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	t.Run("Test checking type - float64", func(t *testing.T) {
		p := &Interpreter{}
		v := float64(3.4)
		err := p.checkType(&Token{}, v, []reflect.Kind{reflect.Float64})
		assert.NoError(t, err)
	})
	t.Run("Test checking type - bool", func(t *testing.T) {
		p := &Interpreter{}
		v := true
		err := p.checkType(&Token{}, v, []reflect.Kind{reflect.Float64})
		fmt.Println(err)
		assert.Error(t, err)
	})

	t.Run("Test evaluating unary expr - Bang", func(t *testing.T) {
		p := &Interpreter{}

		res, _ := p.EvaluateExpr(&UnaryExpr{
			Operator: &Token{Type: Bang},
			Right:    &LiteralExpr{Value: nil},
		})
		assert.True(t, res.(bool))

		res, _ = p.EvaluateExpr(&UnaryExpr{
			Operator: &Token{Type: Bang},
			Right:    &LiteralExpr{Value: false},
		})
		assert.True(t, res.(bool))

		res, _ = p.EvaluateExpr(&UnaryExpr{
			Operator: &Token{Type: Bang},
			Right:    &LiteralExpr{Value: "test"},
		})
		assert.False(t, res.(bool))
	})

	t.Run("Test evaluating unary expr - Minus", func(t *testing.T) {
		p := &Interpreter{}

		_, err := p.EvaluateExpr(&UnaryExpr{
			Operator: &Token{Type: Minus},
			Right:    &LiteralExpr{Value: nil},
		})
		assert.Error(t, err)

		res, _ := p.EvaluateExpr(&UnaryExpr{
			Operator: &Token{Type: Minus},
			Right:    &LiteralExpr{Value: 7.0},
		})
		assert.Equal(t, -7.0, res.(float64))
	})

	t.Run("Test evaluating binary expr - Plus", func(t *testing.T) {
		p := &Interpreter{}

		res, _ := p.EvaluateExpr(&BinaryExpr{
			Operator: &Token{Type: Plus},
			Left:     &LiteralExpr{Value: "test"},
			Right:    &LiteralExpr{Value: "string"},
		})
		assert.Equal(t, "teststring", res.(string))

		res, _ = p.EvaluateExpr(&BinaryExpr{
			Operator: &Token{Type: Plus},
			Left:     &LiteralExpr{Value: 1.0},
			Right:    &LiteralExpr{Value: 2.0},
		})
		assert.Equal(t, 3.0, res.(float64))
	})

	t.Run("Test evaluating binary expr - Greater", func(t *testing.T) {
		p := &Interpreter{}

		_, err := p.EvaluateExpr(&BinaryExpr{
			Operator: &Token{Type: Greater},
			Left:     &LiteralExpr{Value: "bad"},
			Right:    &LiteralExpr{Value: "test"},
		})
		assert.Error(t, err)

		res, _ := p.EvaluateExpr(&BinaryExpr{
			Operator: &Token{Type: Greater},
			Left:     &LiteralExpr{Value: 1.0},
			Right:    &LiteralExpr{Value: 2.0},
		})
		assert.False(t, res.(bool))
	})
}
