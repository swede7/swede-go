package lsp

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"me.weldnor/swede/core/common"
	"me.weldnor/swede/core/parser"
)

type token struct {
	startPosition common.Position
	endPosition   common.Position
	tokenType     tokenType
}

type tokenType string

const (
	comment tokenType = "comment"
	step    tokenType = "step"
	tag     tokenType = "tag"
)

func getTokens(rootNode *parser.Node) []token {
	tokens := make([]token, 0)

	parser.VisitNode(rootNode, func(n *parser.Node) {
		var tokenType tokenType
		switch n.Type {
		case parser.COMMENT:
			tokenType = comment
		case parser.STEP:
			tokenType = step
		case parser.TAG:
			tokenType = tag
		default:
			return
		}

		currentToken := token{n.StartPosition, n.EndPosition, tokenType}
		tokens = append(tokens, currentToken)
	})

	return tokens
}

func highlight(rootNode *parser.Node) *protocol.SemanticTokens {
	data := make([]uint32, 0)

	var prevToken *token

	for _, token := range getTokens(rootNode) {
		startPosition := token.startPosition
		endPosition := token.endPosition

		var deltaLine int
		if prevToken == nil {
			deltaLine = startPosition.Line
		} else {
			deltaLine = startPosition.Line - prevToken.startPosition.Line
		}
		data = append(data, uint32(deltaLine))

		var deltaStart int
		if deltaLine == 0 {
			if prevToken == nil {
				deltaStart = startPosition.Column
			} else {
				deltaStart = startPosition.Column - prevToken.startPosition.Column
			}
		} else {
			deltaStart = startPosition.Column
		}
		data = append(data, uint32(deltaStart))

		length := endPosition.Offset - startPosition.Offset + 1
		data = append(data, uint32(length))

		tokenType := mapTokenType(token.tokenType)
		data = append(data, uint32(tokenType))

		// todo
		tokenModifiers := 0
		data = append(data, uint32(tokenModifiers))

		prevToken = &token
	}

	return &protocol.SemanticTokens{
		Data: data,
	}
}

func mapTokenType(tokenType tokenType) int {
	switch tokenType {
	case comment:
		return 0
	case step:
		return 1
	case tag:
		return 3
	default:
		panic("oops")
	}
}
