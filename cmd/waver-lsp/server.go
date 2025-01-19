package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"reflect"
	"regexp"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/avoronkov/waver/lib/seq/syntaxgen"
	"github.com/avoronkov/waver/lib/utils"
	"github.com/avoronkov/waver/static"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Server struct {
	params *syntaxgen.Params
	parser *parser.Parser

	docs map[string][]string
}

func NewServer() *Server {

	params := syntaxgen.NewParams()
	return &Server{
		params: params,
		parser: parser.New(),
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

	docUri := params.TextDocument.URI
	posLine := int(params.Position.Line)
	posChar := int(params.Position.Character)

	// complete pragmas
	if s.lineMatchRe(docUri, posLine, posChar, pragmaRe) {
		completionItems = append(completionItems, s.completePragmas()...)
	}

	// complete sample files
	if s.lineMatchRe(docUri, posLine, posChar, sampleFileRe) {
		completionItems = append(completionItems, s.completeSampleFiles()...)
	}

	if s.lineMatchRe(docUri, posLine, posChar, waveNameRe) {
		completionItems = append(completionItems, s.completeWaveNames()...)
	}

	if s.isPragmaOptions(docUri, posLine) {
		if s.lineMatchRe(docUri, posLine, posChar, filterRe) {
			completionItems = append(completionItems, s.completeFilters()...)
		}

		if s.lineMatchRe(docUri, posLine, posChar, filterOptionRe) {
			curFilter := s.currentFilter(docUri, posLine)
			completionItems = append(completionItems, s.completeFilterOptions(curFilter)...)
		}
	}

	return completionItems, nil
}

var pragmaRe = regexp.MustCompile(`^%%?\s*\S+$`)
var sampleFileRe = regexp.MustCompile(`^%%?\s*sample\s+\w+\s+"[^"]*$`)
var waveNameRe = regexp.MustCompile(`^%%?\s*(wave|inst)\s+\w+\s+"[^"]*$`)
var filterRe = regexp.MustCompile(`^-\s*\S*$`)
var filterOptionRe = regexp.MustCompile(`^\s+\S*$`)

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
	res := re.MatchString(str)
	slog.Info("lineMatchRe", "str", str, "result", res)
	return res
}

func (s *Server) completePragmas() (items []protocol.CompletionItem) {
	// kind := protocol.CompletionItemKindProperty
	for pragma, meta := range s.parser.PragmaParsers {
		item := pragma
		detail := meta.Usage
		items = append(items, protocol.CompletionItem{
			Label:      item,
			Detail:     &detail,
			InsertText: &item,
			// Kind:       &kind,
			Deprecated: &meta.Deprecated,
		})
		slog.Info("completePragmas", "label", item)
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

func (s *Server) isPragmaOptions(doc string, line int) (result bool) {
	lines, ok := s.docs[doc]
	if !ok {
		slog.Warn("Document not found", "name", doc)
		return false
	}
	if line >= len(lines) {
		slog.Warn("Line index out of range", "line", line)
		return false
	}
	for i := range line {
		str := lines[i]
		if strings.HasPrefix(str, "%%") {
			result = !result
			slog.Info("isPragmaOptions", "line", i, "res", result)
		}
	}
	if strings.HasPrefix(lines[line], "%%") {
		return false
	}
	slog.Info("isPragmaOptions", "result", result)
	return result
}

func (s *Server) completeFilters() (items []protocol.CompletionItem) {
	kind := protocol.CompletionItemKindFunction
	for name, obj := range filters.Filters {
		item := name
		filterOptions := utils.NewSet[string]()

		v := reflect.TypeOf(obj)
		nField := v.NumField()
		for i := 0; i < nField; i++ {
			fld := v.Field(i)
			tagsRaw := fld.Tag.Get("option")
			if tagsRaw == "" {
				continue
			}
			tags := strings.Split(tagsRaw, ",")
			filterOptions.Add(tags...)
		}
		detail := strings.Join(filterOptions.Values(), ", ")
		items = append(items, protocol.CompletionItem{
			Label:      item,
			Detail:     &detail,
			InsertText: &item,
			Kind:       &kind,
		})
	}
	return items
}

var currentFilterRe = regexp.MustCompile(`^-\s*(\S+):`)

func (s *Server) currentFilter(doc string, line int) string {
	lines, ok := s.docs[doc]
	if !ok {
		slog.Warn("Document not found", "name", doc)
		return ""
	}
	if line >= len(lines) {
		slog.Warn("Line index out of range", "line", line)
		return ""
	}
	res := ""
	for i := range line {
		str := lines[i]
		if str == "%%" {
			res = ""
			continue
		}
		matches := currentFilterRe.FindStringSubmatch(str)
		if len(matches) >= 2 {
			res = matches[1]
		}
	}
	return res
}

func (s *Server) completeFilterOptions(filter string) (items []protocol.CompletionItem) {
	if filter == "" {
		return nil
	}
	kind := protocol.CompletionItemKindProperty
	for name, obj := range filters.Filters {
		if name != filter {
			continue
		}

		v := reflect.TypeOf(obj)
		nField := v.NumField()
		for i := 0; i < nField; i++ {
			fld := v.Field(i)
			tagsRaw := fld.Tag.Get("option")
			if tagsRaw == "" {
				continue
			}
			tags := strings.Split(tagsRaw, ",")
			for _, tag := range tags {
				item := tag
				detail := fld.Type.String()
				items = append(items, protocol.CompletionItem{
					Label:      item,
					Detail:     &detail,
					InsertText: &item,
					Kind:       &kind,
				})
			}
		}
		break
	}
	return items
}
