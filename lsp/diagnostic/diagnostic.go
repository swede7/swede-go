package diagnostic

import (
	"strings"

	stepDefinitionParser "github.com/swede7/swede-go/core/lang/swede-step-definition/parser"
	parser "github.com/swede7/swede-go/core/lang/swede/parser"
	"github.com/swede7/swede-go/lsp/context"
	"github.com/swede7/swede-go/lsp/util"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Diagnostic() []protocol.Diagnostic {
	if context.GetContext().FileExtension == context.Swede {
		return DiagnosticSwede()
	}

	if context.GetContext().FileExtension == context.Go {
		return DiagnosticGo()
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

func DiagnosticGo() []protocol.Diagnostic {
	diagnostics := make([]protocol.Diagnostic, 0)

	code := context.GetContext().Code

	for lineNumber, line := range strings.Split(code, "\n") {
		if !isSwedeStepDefinitionComment(line) {
			continue
		}

		util.Logger.Println("found step definition comment")
		util.Logger.Println(line)

		lineOffset := len(stepDefinitionPrefix)

		stepDefinitionString := line[lineOffset:]
		util.Logger.Printf("step definition string: %s\n", stepDefinitionString)

		parserResult, err := stepDefinitionParser.Parse(stepDefinitionString)
		if err != nil {
			util.Logger.Println("cant parse step definition")
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      protocol.UInteger(lineNumber),
						Character: protocol.UInteger(0),
					},
					End: protocol.Position{
						Line:      protocol.UInteger(lineNumber),
						Character: protocol.UInteger(len(line)),
					},
				},
				Message: "cannot parse step definition comment",
			})

			continue
		}
		util.Logger.Printf("hi here")
		parserErrors := parserResult.Errors

		util.Logger.Printf("found %d errors in step definition\n", len(parserErrors))

		for _, parserError := range parserErrors {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      protocol.UInteger(lineNumber),
						Character: protocol.UInteger(lineOffset + parserError.StartPosition.Column),
					},
					End: protocol.Position{
						Line:      protocol.UInteger(lineNumber),
						Character: protocol.UInteger(lineOffset + parserError.EndPosition.Column),
					},
				},
				Message: parserError.Message,
			})
		}
	}

	return diagnostics
}

var stepDefinitionPrefix = "// swede:step"

func isSwedeStepDefinitionComment(line string) bool {
	return strings.HasPrefix(line, stepDefinitionPrefix)
}
