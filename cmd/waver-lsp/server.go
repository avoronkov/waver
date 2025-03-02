package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/avoronkov/waver/lib/seq/syntaxgen"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Server struct {
	params *syntaxgen.Params
	parser *parser.Parser

	docs            map[string][]string
	hoverInfo       map[string]string
	definitionsInfo map[string]map[string]protocol.Location
	languageDefs    map[string]protocol.Location

	// Completion caches
	pragmasCompletions       []protocol.CompletionItem
	sampleFilesCompletions   []protocol.CompletionItem
	waveNamesCompletions     []protocol.CompletionItem
	filtersCompletions       []protocol.CompletionItem
	filterOptionsCompletions map[string][]protocol.CompletionItem
	functionsCompletions     []protocol.CompletionItem
	modifiersCompletions     []protocol.CompletionItem
}

func NewServer() *Server {
	params := syntaxgen.NewParams()
	s := &Server{
		params:          params,
		parser:          parser.New(),
		docs:            make(map[string][]string),
		hoverInfo:       make(map[string]string),
		definitionsInfo: make(map[string]map[string]protocol.Location),
		languageDefs:    make(map[string]protocol.Location),
	}
	s.initHoverInfo()
	s.initLanguageDefinitions()
	return s
}

func (s *Server) Initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "Initializing server...")

	capabilities := handler.CreateServerCapabilities()

	capabilities.CompletionProvider = &protocol.CompletionOptions{}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func (s *Server) Shutdown(context *glsp.Context) error {
	return nil
}

func (s *Server) TextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.docs[params.TextDocument.URI] = strings.Split(params.TextDocument.Text, "\n")
	s.updateDocumentDefinitions(params.TextDocument.URI)
	return nil
}

func (s *Server) TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if wc, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			s.docs[params.TextDocument.URI] = strings.Split(wc.Text, "\n")
		} else {
			slog.Error("Unsupported change type", "type", fmt.Sprintf("%T", params.ContentChanges[0]))
		}
	}
	s.updateDocumentDefinitions(params.TextDocument.URI)
	return nil
}
