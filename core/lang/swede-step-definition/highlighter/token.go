package highlight

import (
	"me.weldnor/swede/core/lang/common"
	"me.weldnor/swede/core/lang/swede-step-definition/parser"
)

type TokenType string

const (
	Variable TokenType = "variable"
	Type     TokenType = "type"
)

type Token struct {
	Start     common.Position
	Length    int
	Type      TokenType
	modifiers []string
}

func Highlight(code string) []Token {
	parserResult, err := parser.Parse(code)

	if err != nil {
		return nil
	}

	if len(parserResult.Errors) > 0 {
		return nil
	}

	rootNode := parserResult.RootNode

	tokens := make([]Token, 0)

	common.VisitNode(&rootNode, func(node *common.Node) {
		switch node.Type {
		case parser.VARIABLE_NAME:
			tokens = append(tokens, Token{
				Start:  node.StartPosition,
				Length: len(node.Value),
				Type:   Variable,
			})
		case parser.VARIABLE_TYPE:
			tokens = append(tokens, Token{
				Start:  node.StartPosition,
				Length: len(node.Value),
				Type:   Type,
			})
		}
	})

	return tokens
}
