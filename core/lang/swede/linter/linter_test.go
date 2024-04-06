package linter_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lang/swede/linter"
	"me.weldnor/swede/core/lang/swede/parser"
)

func Test_emptyFeatureTextRule(t *testing.T) {
	code := "@tag1 @tag2\nFeature:\n\n"

	parserResult := parser.ParseCode(code)

	rootNode := parserResult.RootNode

	linter := linter.NewLinter(&rootNode)
	errors := linter.Lint()

	fmt.Print(errors)

	assert.Len(t, errors, 1)
}

func Test_emptyScenarioTextRule(t *testing.T) {
	code := "@tag1 @tag2\nFeature: example\n\nScenario:\n\n"

	parserResult := parser.ParseCode(code)

	rootNode := parserResult.RootNode

	linter := linter.NewLinter(&rootNode)
	errors := linter.Lint()

	fmt.Print(errors)

	assert.Len(t, errors, 1)
}

func Test_featureNodeInAnotherPosition(t *testing.T) {
	code := "Scenario:example scenario\n- step1\n\nFeature: example\n"

	parserResult := parser.ParseCode(code)

	rootNode := parserResult.RootNode

	linter := linter.NewLinter(&rootNode)
	errors := linter.Lint()

	fmt.Print(errors)

	assert.Len(t, errors, 1)
}
