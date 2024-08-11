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
		scanner.Scan(fmt.Sprintf("\"%s\"", literal))

		token := scanner.tokens[0]
		assert.Equal(t, String, token.Type)
		assert.Equal(t, token.Literal, literal)

		err := scanner.Scan("\"unterminated")
		assert.Error(t, err)
	})

	t.Run("Test numbers", func(t *testing.T) {
		scanner := ScannerImpl{}

		scanner.Scan("132")
		token := scanner.tokens[0]
		assert.Equal(t, Number, token.Type)
		assert.Equal(t, float64(132), token.Literal)

		scanner.Scan("132.45")
		token = scanner.tokens[0]
		assert.Equal(t, Number, token.Type)
		assert.Equal(t, 132.45, token.Literal)

		scanner.Scan("132.")
		token = scanner.tokens[0]
		assert.Equal(t, Number, token.Type)
		assert.Equal(t, float64(132), token.Literal)
	})

	t.Run("Test identifiers", func(t *testing.T) {
		scanner := ScannerImpl{}

		scanner.Scan("var _my_var")
		assert.Equal(t, Var, scanner.tokens[0].Type)
		assert.Equal(t, Identifier, scanner.tokens[1].Type)
		assert.Equal(t, "_my_var", scanner.tokens[1].Lexeme)

		scanner.Scan("if else")
		assert.Equal(t, If, scanner.tokens[0].Type)
		assert.Equal(t, Else, scanner.tokens[1].Type)
	})
}
