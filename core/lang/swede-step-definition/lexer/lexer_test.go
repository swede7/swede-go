package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lang/swede-step-definition/lexer"
)

func TestLexer(t *testing.T) {
	var testCases = []struct {
		source  string
		lexemes []lexer.Lexeme
	}{
		{
			"Hello world",
			[]lexer.Lexeme{
				{0, 10, lexer.TEXT, "Hello world"},
			},
		},
		{
			"Add <first:int> and <second:int>",
			[]lexer.Lexeme{
				{0, 10, lexer.TEXT, "Add "},
				{0, 10, lexer.TEXT, "first:int"},
				{0, 10, lexer.TEXT, " amd "},
				{0, 10, lexer.TEXT, "second:int"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.source, func(t *testing.T) {
			lexemes := lexer.Lex(testCase.source)
			assert.Equal(t, testCase.lexemes, lexemes)
		})
	}
}
