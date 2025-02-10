package main

import (
	"fmt"
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

	if s.isRegularCode(docUri, posLine) {
		completionItems = append(completionItems, s.completeFunctions()...)
		completionItems = append(completionItems, s.completeModifiers()...)
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

func (s *Server) completePragmas() []protocol.CompletionItem {
	if s.pragmasCompletions != nil {
		return s.pragmasCompletions
	}

	items := []protocol.CompletionItem{}
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

	s.pragmasCompletions = items
	return items
}

func (s *Server) completeSampleFiles() []protocol.CompletionItem {
	if s.sampleFilesCompletions != nil {
		return s.sampleFilesCompletions
	}

	kind := protocol.CompletionItemKindFile
	items := []protocol.CompletionItem{}
	for _, file := range static.ListFiles() {
		items = append(items, protocol.CompletionItem{
			Label: file,
			Kind:  &kind,
		})

	}
	s.sampleFilesCompletions = items
	return items
}

func (s *Server) completeWaveNames() []protocol.CompletionItem {
	if s.waveNamesCompletions != nil {
		return s.waveNamesCompletions
	}

	items := []protocol.CompletionItem{}
	kind := protocol.CompletionItemKindConstant
	for w := range waves.Waves {
		items = append(items, protocol.CompletionItem{
			Label: w,
			Kind:  &kind,
		})
	}

	s.waveNamesCompletions = items
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
	if strings.HasPrefix(lines[line], "%") {
		return false
	}
	slog.Info("isPragmaOptions", "result", result)
	return result
}

func (s *Server) isRegularCode(doc string, line int) bool {
	lines, ok := s.docs[doc]
	if !ok {
		slog.Warn("Document not found", "name", doc)
		return false
	}
	if line >= len(lines) {
		slog.Warn("Line index out of range", "line", line)
		return false
	}
	result := true
	for i := range line {
		str := lines[i]
		if strings.HasPrefix(str, "%%") {
			result = !result
			slog.Info("isRegularCode", "line", i, "res", result)
		}
	}
	if strings.HasPrefix(lines[line], "%") {
		return false
	}
	slog.Info("isRegularCode", "result", result)
	return result
}

func (s *Server) completeFilters() []protocol.CompletionItem {
	if s.filtersCompletions != nil {
		return s.filtersCompletions
	}

	items := []protocol.CompletionItem{}
	s.filterOptionsCompletions = make(map[string][]protocol.CompletionItem)

	kind := protocol.CompletionItemKindFunction
	for name, obj := range filters.Filters {
		item := name
		filterOptions := utils.NewSet[string]()

		v := reflect.TypeOf(obj)
		nField := v.NumField()
		optionItems := []protocol.CompletionItem{}
		for i := 0; i < nField; i++ {
			fld := v.Field(i)
			tagsRaw := fld.Tag.Get("option")
			if tagsRaw == "" {
				continue
			}
			tags := strings.Split(tagsRaw, ",")
			filterOptions.Add(tags...)

			for _, tag := range tags {
				item := tag
				detail := fld.Type.String()
				optionItems = append(optionItems, protocol.CompletionItem{
					Label:      item,
					Detail:     &detail,
					InsertText: &item,
					Kind:       &kind,
				})
			}
		}
		detail := strings.Join(filterOptions.Values(), ", ")
		items = append(items, protocol.CompletionItem{
			Label:      item,
			Detail:     &detail,
			InsertText: &item,
			Kind:       &kind,
		})
		s.filterOptionsCompletions[name] = optionItems
	}

	s.filtersCompletions = items
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

func (s *Server) completeFilterOptions(filter string) []protocol.CompletionItem {
	if filter == "" {
		return nil
	}
	if s.filterOptionsCompletions == nil {
		_ = s.completeFilters()
	}
	return s.filterOptionsCompletions[filter]
}

func (s *Server) completeFunctions() []protocol.CompletionItem {
	if s.functionsCompletions != nil {
		return s.functionsCompletions
	}

	items := []protocol.CompletionItem{}
	kind := protocol.CompletionItemKindFunction
	for token, meta := range s.parser.FuncParsers {
		item := token.String()
		detail := meta.Usage
		items = append(items, protocol.CompletionItem{
			Label:      item,
			Detail:     &detail,
			Kind:       &kind,
			InsertText: &item,
		})
	}

	s.functionsCompletions = items
	return items
}

func (s *Server) completeModifiers() []protocol.CompletionItem {
	if s.modifiersCompletions != nil {
		return s.modifiersCompletions
	}

	items := []protocol.CompletionItem{}
	kind := protocol.CompletionItemKindFunction
	for token, meta := range s.parser.ModParsers {
		item := token.String()
		detail := meta.Usage
		items = append(items, protocol.CompletionItem{
			Label:      item,
			Detail:     &detail,
			Kind:       &kind,
			InsertText: &item,
		})
	}

	s.modifiersCompletions = items
	return items
}
