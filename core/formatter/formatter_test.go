package formatter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestFormatter_Format(t *testing.T) {
	formatter := NewFormatter(code)

	formattedCode, err := formatter.Format()
	assert.Nil(t, err)

	fmt.Print(formattedCode)
}
