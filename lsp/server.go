package lsp

import (
	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"me.weldnor/swede/lsp/context"
	"me.weldnor/swede/lsp/format"
	"me.weldnor/swede/lsp/highlight"
)

const lsName = "swede"

var (
	version string = "0.0.1"
	handler protocol.Handler
)

type LspServer struct{}

func NewLspServer() *LspServer {
	return &LspServer{}
}

func (l *LspServer) Start() {
	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		SetTrace:                       setTrace,
		TextDocumentFormatting:         textDocumentFormatting,
		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
		TextDocumentDidOpen:            textDocumentDidOpen,
		TextDocumentDidSave:            textDocumentDidSave,
		TextDocumentDidChange:          TextDocumentDidChange,
	}
	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func textDocumentDidSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	return nil
}

func TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	contentChangeEvent, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEvent)

	if ok {
		updateCode(contentChangeEvent.Text)

		return nil
	}

	ContentChangeEventWhole, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole)
	if ok {
		updateCode(ContentChangeEventWhole.Text)

		return nil
	}

	panic("can't process documentDidChange event")
}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	updateCode(params.TextDocument.Text)

	return nil
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
	return format.Format()
}

func textDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	return highlight.Highlight()
}

func updateCode(newCode string) {
	context.GetContext().Code = newCode
}

func getCode() string {
	return context.GetContext().Code
}
