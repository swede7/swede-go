package diagnostic

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"me.weldnor/swede/core/lang/swede/parser"
	"me.weldnor/swede/lsp/context"
	"me.weldnor/swede/lsp/util"
)

func Diagnostic() []protocol.Diagnostic {
	if context.GetContext().FileExtension == context.Swede {
		return DiagnosticSwede()
	}

	return []protocol.Diagnostic{}
}

func DiagnosticSwede() []protocol.Diagnostic {
	util.Logger.Println("starting swede diagnostic")

	code := context.GetContext().Code

	parserResult := parser.ParseCode(code)

	diagnostics := make([]protocol.Diagnostic, 0)

	for _, _error := range parserResult.Errors {
		diagnostic := protocol.Diagnostic{
			Message: _error.Message,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(_error.StartPosition.Line),
					Character: uint32(_error.StartPosition.Column),
				},
				End: protocol.Position{
					Line:      uint32(_error.EndPosition.Line),
					Character: uint32(_error.EndPosition.Column) + 1,
				},
			},
		}

		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}
