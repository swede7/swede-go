package lsp

import (
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"me.weldnor/swede/core/formatter"
	"me.weldnor/swede/core/lexer"
	"me.weldnor/swede/core/parser"

	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "swede"

var version string = "0.0.1"
var handler protocol.Handler

type LspServer struct {
}

func NewLspServer() *LspServer {
	return &LspServer{}
}

func (l *LspServer) Start() {
	// commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		SetTrace:                       setTrace,
		TextDocumentFormatting:         textDocumentFormatting,
		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
		TextDocumentDidOpen: func(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
			CODE = params.TextDocument.Text
			return nil
		},
		TextDocumentDidSave: func(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error { return nil },
		TextDocumentDidChange: func(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
			var event, ok = params.ContentChanges[0].(protocol.TextDocumentContentChangeEvent)
			if ok {
				CODE = event.Text
				return nil
			}

			var event1, ok1 = params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole)
			if ok1 {
				CODE = event1.Text
				return nil
			}

			panic("oops")
		},
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.DocumentFormattingProvider = true
	capabilities.SemanticTokensProvider = protocol.SemanticTokensOptions{
		Full:   true,
		Legend: protocol.SemanticTokensLegend{TokenTypes: []string{"comment", "string", "keyword", "parameter"}},
		Range:  false,
	}
	capabilities.TextDocumentSync = 1
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func textDocumentFormatting(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	code := CODE

	lexer := lexer.NewLexer(code)
	parser := parser.NewParser(lexer.Scan())
	parserResult := parser.Parse()

	if len(parserResult.Errors) > 0 {
		return []protocol.TextEdit{}, nil
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

func textDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	code := CODE

	lexer := lexer.NewLexer(string(code))
	lexemes := lexer.Scan()

	parserResult := parser.NewParser(lexemes).Parse()

	semanticTokens := highlight(&parserResult.RootNode)

	return semanticTokens, nil
}
