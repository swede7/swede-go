package diagnostic

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"me.weldnor/swede/core/lexer"
	"me.weldnor/swede/core/parser"
	"me.weldnor/swede/lsp/context"
)

func Diagnostic() []protocol.Diagnostic {
	code := context.GetContext().Code

	lexer := lexer.NewLexer(code)
	parser := parser.NewParser(lexer.Scan())
	parserResult := parser.Parse()

	diagnostics := make([]protocol.Diagnostic, 0)

	for _, error := range parserResult.Errors {
		diagnostic := protocol.Diagnostic{
			Message: error.Message,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(error.StartPosition.Line),
					Character: uint32(error.StartPosition.Column),
				},
				End: protocol.Position{
					Line:      uint32(error.EndPosition.Line),
					Character: uint32(error.EndPosition.Column) + 1,
				},
			},
		}

		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}
