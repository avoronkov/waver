package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"regexp"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/seq/syntaxgen"
	"github.com/avoronkov/waver/static"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Server struct {
	params *syntaxgen.Params

	docs map[string][]string
}

func NewServer() *Server {
	params := syntaxgen.NewParams()
	return &Server{
		params: params,
		docs:   make(map[string][]string),
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
	slog.Info("TextDocumentDidChange", "document", params.TextDocument.URI, "changes", params.ContentChanges, "type", fmt.Sprintf("%T", params.ContentChanges[0]))
	for _, change := range params.ContentChanges {
		if wc, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			s.docs[params.TextDocument.URI] = strings.Split(wc.Text, "\n")
		} else {
			slog.Error("Unsupported change type", "type", fmt.Sprintf("%T", params.ContentChanges[0]))
		}
	}
	return nil
}

func (s *Server) TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered", "error", r)
		}
	}()

	var completionItems []protocol.CompletionItem

	// complete pragmas
	if s.lineMatchRe(params.TextDocument.URI, int(params.Position.Line), int(params.Position.Character), pragmaRe) {
		completionItems = append(completionItems, s.completePragmas()...)
	}

	// complete sample files
	if s.lineMatchRe(params.TextDocument.URI, int(params.Position.Line), int(params.Position.Character), sampleFileRe) {
		completionItems = append(completionItems, s.completeSampleFiles()...)
	}

	if s.lineMatchRe(params.TextDocument.URI, int(params.Position.Line), int(params.Position.Character), waveNameRe) {
		completionItems = append(completionItems, s.completeWaveNames()...)
	}

	return completionItems, nil
}

var pragmaRe = regexp.MustCompile(`^%%?\s*\w+$`)
var sampleFileRe = regexp.MustCompile(`^%%?\s*sample\s+\w+\s+"\S*$`)
var waveNameRe = regexp.MustCompile(`^%%?\s*(wave|inst)\s+\w+\s+"\S*$`)

func (s *Server) lineMatchRe(doc string, line, pos int, re *regexp.Regexp) bool {
	lines, ok := s.docs[doc]
	if !ok {
		slog.Warn("Document not found", "name", doc)
		return false
	}
	if line >= len(lines) {
		slog.Warn("Line index out of range", "line", line)
		return false
	}
	str := lines[line]
	if pos > len(str) {
		slog.Warn("Position index out of range", "pos", pos)
		return false
	}
	str = str[0:pos]
	slog.Info("lineMatchRe", "str", str)
	return re.MatchString(str)
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
			Kind:  &kind,
		})
		slog.Info("Sample", "file", path)
		return nil
	})
	if err != nil {
		slog.Error("WalkDir failed", "error", err)
	}
	return items
}

func (s *Server) completeWaveNames() (items []protocol.CompletionItem) {
	kind := protocol.CompletionItemKindConstant
	for w := range waves.Waves {
		items = append(items, protocol.CompletionItem{
			Label: w,
			Kind:  &kind,
		})
	}
	return items
}
