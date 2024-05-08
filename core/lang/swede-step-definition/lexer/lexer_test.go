package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swede7/swede-go/core/lang/common"
	"github.com/swede7/swede-go/core/lang/swede-step-definition/lexer"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		source  string
		lexemes []common.Lexeme
	}{
		{
			"Hello world",
			[]common.Lexeme{
				makeLexeme(lexer.TEXT, 0, 10, "Hello world"),
			},
		},
		{
			"Add <first:int> and <second:int>",
			[]common.Lexeme{
				makeLexeme(lexer.TEXT, 0, 3, "Add "),
				makeLexeme(lexer.L_BRACKET, 4, 4, "<"),
				makeLexeme(lexer.TEXT, 5, 9, "first"),
				makeLexeme(lexer.COLON, 10, 10, ":"),
				makeLexeme(lexer.TEXT, 11, 13, "int"),
				makeLexeme(lexer.R_BRACKET, 14, 14, ">"),
				makeLexeme(lexer.TEXT, 15, 19, " and "),
				makeLexeme(lexer.L_BRACKET, 20, 20, "<"),
				makeLexeme(lexer.TEXT, 21, 26, "second"),
				makeLexeme(lexer.COLON, 27, 27, ":"),
				makeLexeme(lexer.TEXT, 28, 30, "int"),
				makeLexeme(lexer.R_BRACKET, 31, 31, ">"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.source, func(t *testing.T) {
			lexemes, err := lexer.Lex(testCase.source)

			assert.Nil(t, err)
			assert.Equal(t, testCase.lexemes, lexemes)
		})
	}
}

func makeLexeme(lexemeType common.LexemeType, start int, end int, value string) common.Lexeme {
	return common.Lexeme{
		Type:          lexemeType,
		StartPosition: common.Position{Offset: start, Line: 0, Column: start},
		EndPosition:   common.Position{Offset: end, Line: 0, Column: end},
		Value:         value,
	}
}
