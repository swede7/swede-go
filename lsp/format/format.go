package format

import (
	"errors"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"me.weldnor/swede/core/formatter"
	"me.weldnor/swede/core/parser"
	"me.weldnor/swede/lsp/context"
)

func Format() ([]protocol.TextEdit, error) {
	code := context.GetContext().Code

	parserResult := parser.ParseCode(code)

	if len(parserResult.Errors) > 0 {
		return nil, errors.New("can't parse file")
	}

	formatter := formatter.NewFormatter(&parserResult.RootNode)
	formattedCode, err := formatter.FormatParallel()
	if err != nil {
		panic("oops: failed to format code")
	}

	textEdit := protocol.TextEdit{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 0,
			},
			End: protocol.Position{
				Line:      uint32(strings.Count(code, "\n") + 2),
				Character: 0,
			},
		},
		NewText: formattedCode,
	}

	return []protocol.TextEdit{textEdit}, nil
}
