package parser

import (
	"fmt"
	"testing"

	"me.weldnor/swede/core/lexer"
)

const code string = `
@all
Feature: Basic calculator operations

This feature defines a set of operations that the calculator must support.

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

func TestParser(t *testing.T) {
	lexer := lexer.NewLexer(code)
	lexemes := lexer.Scan()

	parser := NewParser(lexemes)
	result := parser.Parse()

	fmt.Print(result.Errors)
}
