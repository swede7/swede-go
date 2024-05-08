package lsp

import (
	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	"errors"
	"github.com/swede7/swede-go/lsp/util"

	"github.com/swede7/swede-go/lsp/context"
	"github.com/swede7/swede-go/lsp/diagnostic"
	"github.com/swede7/swede-go/lsp/format"
	"github.com/swede7/swede-go/lsp/highlight"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
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

	util.Logger.Println("starting lsp server ...")
	server.RunStdio()
}

func textDocumentDidSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	util.Logger.Printf("document saved")

	publishDiagnostic(ctx.Notify)

	return nil
}

func TextDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	util.Logger.Printf("document changed")
	ContentChangeEventWhole, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole)
	if !ok {
		panic(errors.New("cant process document did change error"))
	}
	//fixme race condition?
	context.GetContext().Code = ContentChangeEventWhole.Text
	context.GetContext().URI = params.TextDocument.URI
	context.GetContext().FileExtension = context.GetFileExtensionByURL(params.TextDocument.URI)

	publishDiagnostic(ctx.Notify)

	return nil
}

func textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	util.Logger.Printf("document opened")

	context.GetContext().Code = params.TextDocument.Text
	context.GetContext().URI = params.TextDocument.URI
	context.GetContext().FileExtension = context.GetFileExtensionByURL(params.TextDocument.URI)

	publishDiagnostic(ctx.Notify)

	return nil
}

func initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
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

func initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(ctx *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)

	return nil
}

func setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)

	return nil
}

// region LSP features

func textDocumentFormatting(ctx *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	util.Logger.Printf("request document formatting")
	return format.Format()
}

func textDocumentSemanticTokensFull(ctx *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	util.Logger.Printf("request document semantic tokens")
	return highlight.Highlight()
}

func publishDiagnostic(notifyFunc glsp.NotifyFunc) {
	util.Logger.Printf("request document diagnostic")
	go notifyFunc(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         context.GetContext().URI,
		Diagnostics: diagnostic.Diagnostic(),
	})
}
