package lsp

import (
	"log"
	"net/url"
	"os"

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
		Initialize:             initialize,
		Initialized:            initialized,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentFormatting: textDocumentFormatting,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.DocumentFormattingProvider = true

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
	uri := params.TextDocument.URI
	filepath := uriToFilePath(uri)
	dat, err := os.ReadFile(filepath)

	if err != nil {
		panic("oops: failed to read " + filepath + err.Error())
	}

	code := string(dat)

	lexer := lexer.NewLexer(code)
	parser := parser.NewParser(lexer.Scan())
	parserResult := parser.Parse()

	if len(parserResult.Errors) > 0 {
		return []protocol.TextEdit{}, nil
	}

	rootNode := parserResult.RootNode

	formatter := formatter.NewFormatter(&parserResult.RootNode)
	formattedCode, err := formatter.Format()

	if err != nil {
		panic("oops: failed to format code")
	}

	textEdit := protocol.TextEdit{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(rootNode.StartPosition.Line),
				Character: uint32(rootNode.StartPosition.Column),
			},
			End: protocol.Position{
				Line:      uint32(rootNode.EndPosition.Line),
				Character: uint32(rootNode.StartPosition.Column),
			},
		},
		NewText: formattedCode,
	}

	return []protocol.TextEdit{textEdit}, nil
}

func uriToFilePath(uri string) string {
	// Parse the URI
	parsedURI, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the scheme is 'file'
	if parsedURI.Scheme != "file" {
		log.Fatal("The provided URI does not have a file scheme.")
	}

	return parsedURI.Path[1:]
}
