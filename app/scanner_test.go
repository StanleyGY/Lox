package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	t.Run("Test non-alpha lexeme", func(t *testing.T) {
		scanner := ScannerImpl{}
		scanner.Scan("!*+-/=<>!=")

		expectedTypes := []int{
			Bang,
			Star,
			Plus,
			Minus,
			Slash,
			Equal,
			Less,
			Greater,
			BangEqual,
			EOF,
		}
		for idx := range scanner.tokens {
			assert.Equal(t, expectedTypes[idx], scanner.tokens[idx].Type)
		}
	})

	t.Run("Test comments", func(t *testing.T) {
		scanner := ScannerImpl{}
		scanner.Scan("// this is a comment")

		expectedTypes := []int{
			EOF,
		}
		for idx := range scanner.tokens {
			assert.Equal(t, expectedTypes[idx], scanner.tokens[idx].Type)
		}
	})

	t.Run("Test multi-line", func(t *testing.T) {
		scanner := ScannerImpl{}
		scanner.Scan("/ \n // this is a comment")

		expectedTypes := []int{
			Slash,
			EOF,
		}
		for idx := range scanner.tokens {
			assert.Equal(t, expectedTypes[idx], scanner.tokens[idx].Type)
		}
	})

	t.Run("Test string literal", func(t *testing.T) {
		literal := "this is a string literal"

		scanner := ScannerImpl{}
		tokens, _ := scanner.Scan(fmt.Sprintf("\"%s\"", literal))

		assert.Equal(t, String, tokens[0].Type)
		assert.Equal(t, tokens[0].Literal, literal)

		_, err := scanner.Scan("\"unterminated")
		assert.Error(t, err)
	})

	t.Run("Test numbers", func(t *testing.T) {
		scanner := ScannerImpl{}

		tokens, _ := scanner.Scan("132")
		assert.Equal(t, Number, tokens[0].Type)
		assert.Equal(t, float64(132), tokens[0].Literal)

		tokens, _ = scanner.Scan("132.45")
		assert.Equal(t, Number, tokens[0].Type)
		assert.Equal(t, 132.45, tokens[0].Literal)

		tokens, _ = scanner.Scan("132.")
		assert.Equal(t, Number, tokens[0].Type)
		assert.Equal(t, float64(132), tokens[0].Literal)
	})

	t.Run("Test identifiers", func(t *testing.T) {
		scanner := ScannerImpl{}

		tokens, _ := scanner.Scan("var _my_var")
		assert.Equal(t, Var, tokens[0].Type)
		assert.Equal(t, Identifier, tokens[1].Type)
		assert.Equal(t, "_my_var", tokens[1].Lexeme)

		tokens, _ = scanner.Scan("if else")
		assert.Equal(t, If, tokens[0].Type)
		assert.Equal(t, Else, tokens[1].Type)
	})
}
