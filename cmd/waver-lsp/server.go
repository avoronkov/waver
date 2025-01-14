package main

import (
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/avoronkov/waver/lib/seq/syntaxgen"
	"github.com/avoronkov/waver/static"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Server struct {
	params *syntaxgen.Params
}

func NewServer() *Server {
	params := syntaxgen.NewParams()
	return &Server{
		params: params,
	}
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

func (s *Server) TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	// slog.Info("TextDocumentDidChange", "context", context, "params", params)
	return nil
}

func (s *Server) TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered", "error", r)
		}
	}()

	var completionItems []protocol.CompletionItem

	slog.Info("TDC", "context.Params", context.Params, "params", params)
	slog.Info(fmt.Sprintf("Params: %v (%v)", params.PartialResultToken, params.WorkDoneProgressParams))

	// complete pragmas
	completionItems = append(completionItems, s.completePragmas()...)

	// complete sample files
	completionItems = append(completionItems, s.completeSampleFiles()...)

	slog.Info("Completion", "items", completionItems)
	return completionItems, nil
}

func (s *Server) completePragmas() (items []protocol.CompletionItem) {
	for _, pragma := range s.params.Pragmas {
		item := pragma
		// detail := "[pragma]"
		kind := protocol.CompletionItemKindProperty
		items = append(items, protocol.CompletionItem{
			Label: item,
			// Detail:     &detail,
			InsertText: &item,
			Kind:       &kind,
		})
	}
	return items
}

func (s *Server) completeSampleFiles() (items []protocol.CompletionItem) {
	subdir := "samples"
	subdirlen := len(subdir) + 1
	err := fs.WalkDir(static.Files, subdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		kind := protocol.CompletionItemKindFile

		items = append(items, protocol.CompletionItem{
			Label: path[subdirlen:],
			// Detail:     &detail,
			Kind: &kind,
		})
		slog.Info("Sample", "file", path)
		return nil
	})
	if err != nil {
		slog.Error("WalkDir failed", "error", err)
	}
	return items
}
