package formatter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"me.weldnor/swede/core/lexer"
	"me.weldnor/swede/core/parser"
)

const code string = `
       @all
Feature:           Basic calculator operations


# Comment example




     @pass              @automated
Scenario:     Addition
- Enter "2 + 2"

- Click on calculation button


- Check that the answer is "5"

        @fail
Scenario: Division by zero

- Enter "5 / 0"
- Click on calculation button



- An exception must be thrown

`

func TestFormatter_FormatParallel(t *testing.T) {
	lexer := lexer.NewLexer(code)
	parser := parser.NewParser(lexer.Scan())

	rootNode := parser.Parse().RootNode
	formatter := NewFormatter(&rootNode)

	formattedCode, err := formatter.FormatParallel()
	assert.Nil(t, err)

	fmt.Print(formattedCode)
}
