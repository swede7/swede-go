package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lang/swede-step-definition/lexer"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		source  string
		lexemes []lexer.Lexeme
	}{
		{
			"Hello world",
			[]lexer.Lexeme{
				{lexer.TEXT, 0, 10, "Hello world"},
			},
		},
		{
			"Add <first:int> and <second:int>",
			[]lexer.Lexeme{
				{lexer.TEXT, 0, 3, "Add "},
				{lexer.L_BRACKET, 4, 4, "<"},
				{lexer.TEXT, 5, 9, "first"},
				{lexer.COLON, 10, 10, ":"},
				{lexer.TEXT, 11, 13, "int"},
				{lexer.R_BRACKET, 14, 14, ">"},
				{lexer.TEXT, 15, 19, " and "},
				{lexer.L_BRACKET, 20, 20, "<"},
				{lexer.TEXT, 21, 26, "second"},
				{lexer.COLON, 27, 27, ":"},
				{lexer.TEXT, 28, 30, "int"},
				{lexer.R_BRACKET, 31, 31, ">"},
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
