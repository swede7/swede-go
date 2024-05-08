package parser_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lang/swede-step-definition/model"
	"me.weldnor/swede/core/lang/swede-step-definition/parser"
)

func TestParse(t *testing.T) {
	code := "Add <first:int> and <second:string>"

	parserResult, err := parser.Parse(code)

	assert.Nil(t, err)

	rootNode := parserResult.RootNode

	assert.Nil(t, err)
	assert.NotNil(t, rootNode)

	assert.Len(t, rootNode.Children, 4)

	firstTextNode := rootNode.Children[0]
	assert.Equal(t, parser.TEXT, firstTextNode.Type)
	assert.Equal(t, "Add ", firstTextNode.Value)

	firstVariableNode := rootNode.Children[1]
	assert.Equal(t, parser.VARIABLE, firstVariableNode.Type)

	firstVariableNameNode := firstVariableNode.GetChildByType(parser.VARIABLE_NAME)
	assert.Equal(t, "first", firstVariableNameNode.Value)
	firstVariableTypeNode := firstVariableNode.GetChildByType(parser.VARIABLE_TYPE)
	assert.Equal(t, "int", firstVariableTypeNode.Value)

	secondTextNode := rootNode.Children[2]
	assert.Equal(t, parser.TEXT, secondTextNode.Type)
	assert.Equal(t, " and ", secondTextNode.Value)

	secondVariableNode := rootNode.Children[3]
	assert.Equal(t, parser.VARIABLE, secondVariableNode.Type)

	secondVariableNameNode := secondVariableNode.GetChildByType(parser.VARIABLE_NAME)
	assert.Equal(t, "second", secondVariableNameNode.Value)
	secondVariableTypeNode := secondVariableNode.GetChildByType(parser.VARIABLE_TYPE)
	assert.Equal(t, "string", secondVariableTypeNode.Value)

	// validate _model

	_model := parserResult.StepDefinition

	assert.Equal(t, 2, len(_model.Variables))

	firstVariableModel := _model.Variables[0]
	assert.Equal(t, "first", firstVariableModel.Name)
	assert.Equal(t, model.Int, firstVariableModel.Type)

	secondVariableModel := _model.Variables[1]
	assert.Equal(t, "second", secondVariableModel.Name)
	assert.Equal(t, model.String, secondVariableModel.Type)
}

func TestParseWithError(t *testing.T) {
	//todo remove it
	code := "Add second <asdasd:"

	parserResult, err := parser.Parse(code)

	if err != nil {
		panic(err)
	}

	fmt.Println(parserResult)
}
