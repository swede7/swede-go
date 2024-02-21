package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lexer"
	"me.weldnor/swede/core/parser"
)

const code string = `@example 
Feature:          calculator



# comment

@positive 
Scenario: Test     addition
-      Add "2" and "2"
- Check that result is "4"



       @negative 
Scenario: Test     addition? but result is not correct
-                Add "2" and "2"


- Check that result is "5"
`

const formattedCode string = `@example 
Feature: calculator

# comment

@positive 
Scenario: Test     addition
- Add "2" and "2"
- Check that result is "4"

@negative 
Scenario: Test     addition? but result is not correct
- Add "2" and "2"
- Check that result is "5"

`

func TestFormatter_FormatParallel(t *testing.T) {
	lexer := lexer.NewLexer(code)
	parser := parser.NewParser(lexer.Scan())

	rootNode := parser.Parse().RootNode
	formatter := NewFormatter(&rootNode)

	result, _ := formatter.FormatParallel()
	assert.Equal(t, formattedCode, result)
}
