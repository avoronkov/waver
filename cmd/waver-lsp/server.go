package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/avoronkov/waver/lib/seq/syntaxgen"
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
	slog.Info("TextDocumentDidChange", "context", context, "params", params)
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

	fpath := params.TextDocument.URI
	if strings.HasPrefix(fpath, "file://") {
		fpath = fpath[6:]
	}
	f, err := os.Open(fpath)
	if err != nil {
		slog.Error("Open file", "file", fpath, "error", err)
		return nil, nil
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		slog.Error("ReadAll", "error", err)
		return nil, nil
	}
	lines := strings.Split(string(data), "\n")

	line := lines[params.Position.Line]

	slog.Info("Line", "string", line, "line", params.Position.Line, "character", params.Position.Character)

	if strings.HasPrefix(line, "%") {
		// complete pragmas
		for _, pragma := range s.params.Pragmas {
			item := pragma
			detail := "[pragma]"
			completionItems = append(completionItems, protocol.CompletionItem{
				Label:      item,
				Detail:     &detail,
				InsertText: &item,
			})
		}
	}

	slog.Info("Completion", "items", completionItems)
	return completionItems, nil
}
