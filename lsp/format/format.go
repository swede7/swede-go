package format

import (
	"errors"
	"strings"

	"github.com/swede7/swede-go/core/lang/swede/formatter"
	"github.com/swede7/swede-go/core/lang/swede/parser"
	"github.com/swede7/swede-go/lsp/context"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Format() ([]protocol.TextEdit, error) {
	if context.GetContext().FileExtension == context.Swede {
		return FormatSwede()
	}

	return nil, errors.New("unknown file extension")
}

func FormatSwede() ([]protocol.TextEdit, error) {
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
