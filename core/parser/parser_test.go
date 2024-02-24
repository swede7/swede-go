package parser_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/parser"
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
	parserResult := parser.ParseCode(code)

	fmt.Print(parserResult.RootNode)
	fmt.Print(parserResult.Errors)
}

func TestParserAddTagToFeatureRule(t *testing.T) {
	parserResult := parser.ParseCode("@example\nFeature: hello world")

	assert.Empty(t, parserResult.Errors)

	rootNode := parserResult.RootNode
	assert.NotNil(t, rootNode)

	assert.Len(t, rootNode.Children, 1)

	assert.Equal(t, parser.FEATURE, rootNode.Children[0].Type)

	featureNode := rootNode.Children[0]

	assert.Len(t, featureNode.Children, 1)

	assert.Equal(t, parser.TAG, featureNode.Children[0].Type)
}
