package linter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lexer"
	"me.weldnor/swede/core/parser"
)

func Test_emptyFeatureTextRule(t *testing.T) {
	code := "@tag1 @tag2\nFeature:\n\n"

	lexer := lexer.NewLexer(code)
	parser := parser.NewParser(lexer.Scan())

	rootNode := parser.Parse().RootNode

	linter := NewLinter(&rootNode)
	errors := linter.Lint()

	fmt.Print(errors)

	assert.Len(t, errors, 1)
}


func Test_emptyScenarioTextRule(t *testing.T) {
	code := "@tag1 @tag2\nFeature: example\n\nScenario:\n\n"

	lexer := lexer.NewLexer(code)
	parser := parser.NewParser(lexer.Scan())

	rootNode := parser.Parse().RootNode

	linter := NewLinter(&rootNode)
	errors := linter.Lint()

	fmt.Print(errors)

	assert.Len(t, errors, 1)
}