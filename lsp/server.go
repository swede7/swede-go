package lsp

import (
	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	"errors"

	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"me.weldnor/swede/lsp/context"
	"me.weldnor/swede/lsp/diagnostic"
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

func TextDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	ContentChangeEventWhole, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole)
	if !ok {
		panic(errors.New("cant process document did change error"))
	}

	context.GetContext().Code = ContentChangeEventWhole.Text
	context.GetContext().URI = params.TextDocument.URI

	publishDiagnostic(ctx.Notify)

	return nil
}

func textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	context.GetContext().Code = params.TextDocument.Text
	context.GetContext().URI = params.TextDocument.URI

	publishDiagnostic(ctx.Notify)

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

func publishDiagnostic(notifyFunc glsp.NotifyFunc) {
	go notifyFunc(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         context.GetContext().URI,
		Diagnostics: diagnostic.Diagnostic(),
	})
}

func updateCode(newCode string) {
	context.GetContext().Code = newCode
}

func getCode() string {
	return context.GetContext().Code
}
