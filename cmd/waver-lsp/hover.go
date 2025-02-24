package main

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) TextDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	defer func() {
		if e := recover(); e != nil {
			slog.Error("Recovered", "error", e, "stack", debug.Stack())
		}
	}()
	docUri := params.TextDocument.URI
	posLine := int(params.Position.Line)
	posChar := int(params.Position.Character)
	hover := s.findSymbolUnderCursor(docUri, posLine, posChar)
	if hover == nil {
		return nil, nil
	}
	return &protocol.Hover{
		Contents: hover,
	}, nil
}

func trimLeadingSpaces(s string) string {
	ls := strings.Split(s, "\n")
	for i, l := range ls {
		ls[i] = strings.TrimSpace(l)
	}
	return strings.Join(ls, "\n")
}

func (s *Server) initHoverInfo() {
	for pragma, meta := range s.parser.PragmaParsers {
		info := fmt.Sprintf("[pragma] %v: %v", pragma, meta.Usage)
		if meta.Desc != "" {
			info += "\n" + trimLeadingSpaces(meta.Desc)
		}
		s.hoverInfo[pragma] = info
	}

	for token, meta := range s.parser.FuncParsers {
		t := token.String()
		info := fmt.Sprintf("[func] %v: %v", t, meta.Usage)
		if meta.Desc != "" {
			info += "\n" + trimLeadingSpaces(meta.Desc)
		}
		s.hoverInfo[token.String()] = info
	}

	for token, meta := range s.parser.ModParsers {
		if meta.Deprecated {
			continue
		}
		t := token.String()
		info := fmt.Sprintf("[mod] %v: %v", t, meta.Usage)
		s.hoverInfo[token.String()] = info
	}

	for name, obj := range filters.Filters {
		filterOptions := utils.NewSet[string]()

		v := reflect.TypeOf(obj)
		nField := v.NumField()
		for i := range nField {
			fld := v.Field(i)
			tagsRaw := fld.Tag.Get("option")
			if tagsRaw == "" {
				continue
			}
			tags := strings.Split(tagsRaw, ",")
			filterOptions.Add(tags...)
		}
		detail := strings.Join(filterOptions.Values(), ", ")
		desc := ""
		if d, ok := obj.(descer); ok {
			desc = "\n" + d.Desc()
		}
		info := fmt.Sprintf("[filter] %v: %v%v", name, detail, desc)
		s.hoverInfo[name] = info
	}
}

func (s *Server) findSymbolUnderCursor(doc string, line, pos int) *protocol.MarkupContent {
	lines, ok := s.docs[doc]
	if !ok {
		slog.Warn("Document not found", "name", doc)
		return nil
	}
	if line >= len(lines) {
		slog.Warn("Line index out of range", "line", line)
		return nil
	}
	str := lines[line]

	symbol := string(str[pos])
	slog.Info("findSymbolUnderCursor", "line", line, "pos", pos, "symbol", symbol)
	if meta, ok := s.hoverInfo[symbol]; ok {
		return &protocol.MarkupContent{
			Kind:  protocol.MarkupKindPlainText,
			Value: meta,
		}
	}

	word := s.findWordUnderCursor(doc, line, pos)
	if meta, ok := s.hoverInfo[word]; ok {
		return &protocol.MarkupContent{
			Kind:  protocol.MarkupKindPlainText,
			Value: meta,
		}
	}
	return nil
}

func (s *Server) findWordUnderCursor(doc string, line, pos int) string {
	lines, ok := s.docs[doc]
	if !ok {
		slog.Warn("Document not found", "name", doc)
		return ""
	}
	if line >= len(lines) {
		slog.Warn("Line index out of range", "line", line)
		return ""
	}
	str := lines[line]
	i := pos - 1
	if i >= 0 {
		for ; i >= 0; i-- {
			if !isWordSymbol(str[i]) {
				i++
				break
			}
		}
		if i < 0 {
			i = 0
		}
	} else {
		i = 0
	}
	slog.Debug("findSymbolUnderCursor", "i", i)
	j := pos + 1
	if j < len(str) {
		for ; j < len(str); j++ {
			if !isWordSymbol(str[j]) {
				break
			}
		}
		if j >= len(str) {
			j = len(str)
		}
	} else {
		j = len(str)
	}
	slog.Debug("findSymbolUnderCursor", "j", j)
	word := str[i:j]
	return word
}

func isWordSymbol(r byte) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '\''
}
