package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/common"
)

const code string = `
@all
Feature: Basic calculator operations

This feature defines a set of operations that the calculator must support.

# Comment example

@pass @automated
Scenario: Addition
- Enter "2 + 2"
- Click on calculation button
- Check that the answer is "5"

@fail
Scenario: Division by zero
- Enter "5 / 0"
- Click on calculation button
- Аn exception must be thrown

`

func TestLexerForCodeExample(t *testing.T) {
	lexer := NewLexer(code)
	lexemes := lexer.Scan()

	expectedLexemes := []Lexeme{
		{AT_CHR, common.Position{Offset: 1, Line: 1, Column: 0}, common.Position{Offset: 1, Line: 1, Column: 0}, "@"},
		{WORD, common.Position{Offset: 2, Line: 1, Column: 1}, common.Position{Offset: 4, Line: 1, Column: 3}, "all"},
	}

	for _, expectedLexeme := range expectedLexemes {
		assert.Contains(t, lexemes, expectedLexeme)
	}
}
