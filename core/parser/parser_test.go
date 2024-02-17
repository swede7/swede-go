package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lexer"
)

const code string = `
@all
Feature: Basic calculator operations

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

func TestParserWithValidCode(t *testing.T) {
	lexer := lexer.NewLexer(code)
	lexemes := lexer.Scan()

	parser := NewParser(lexemes)
	result := parser.Parse()

	fmt.Print(result.RootNode)
	fmt.Print(result.Errors)
}

func TestParserAddTagToFeatureRule(t *testing.T) {
	lexer := lexer.NewLexer("@example\nFeature: hello world")
	lexemes := lexer.Scan()

	parser := NewParser(lexemes)
	result := parser.Parse()

	assert.Empty(t, result.Errors)

	rootNode := result.RootNode
	assert.NotNil(t, rootNode)

	assert.Len(t, rootNode.Children, 1)

	assert.Equal(t, FEATURE, rootNode.Children[0].Type)

	featureNode := rootNode.Children[0]

	assert.Len(t, featureNode.Children, 1)

	assert.Equal(t, TAG, featureNode.Children[0].Type)
}
