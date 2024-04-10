package highlight

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"me.weldnor/swede/core/lang/common"
	"me.weldnor/swede/core/lang/swede/parser"
	"me.weldnor/swede/lsp/context"
)

type highlightToken struct {
	startPosition common.Position
	endPosition   common.Position
	tokenType     tokenType
}

type tokenType int

const (
	comment tokenType = 0
	step    tokenType = 1
	tag     tokenType = 3
)

func Highlight() (*protocol.SemanticTokens, error) {
	code := context.GetContext().Code

	parserResult := parser.ParseCode(code)

	return getSemanticTokensByAst(&parserResult.RootNode), nil
}

func getSemanticTokensByAst(rootNode *common.Node) *protocol.SemanticTokens {
	data := make([]uint32, 0)

	var prevToken *highlightToken

	for _, token := range getHighlightTokensByAst(rootNode) {
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

		data = append(data, uint32(token.tokenType))

		tokenModifiers := 0
		data = append(data, uint32(tokenModifiers))

		prevToken = &token
	}

	return &protocol.SemanticTokens{
		Data: data,
	}
}

func getHighlightTokensByAst(rootNode *common.Node) []highlightToken {
	tokens := make([]highlightToken, 0)

	common.VisitNode(rootNode, func(n *common.Node) {
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

		newToken := highlightToken{n.StartPosition, n.EndPosition, tokenType}
		tokens = append(tokens, newToken)
	})

	return tokens
}
