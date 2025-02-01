package main

import (
	"log/slog"
	"regexp"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var instDefinitionRe = regexp.MustCompile(`^%%?\s*(?:wave|inst|sample)\s+(\S+)\s`)
var varDefinitionRe = regexp.MustCompile(`^(\S+)\s*=`)
var funcDefinitionRe = regexp.MustCompile(`^(\S+)\s+\S+\s*=`)

func (s *Server) updateDocumentDefinitions(doc string) error {
	lines, ok := s.docs[doc]
	if !ok {
		slog.Warn("Document not found", "name", doc)
		return nil
	}
	info := map[string]protocol.Location{}
	for i, line := range lines {
		for _, re := range []*regexp.Regexp{instDefinitionRe, varDefinitionRe, funcDefinitionRe} {
			matches := re.FindStringSubmatchIndex(line)
			if len(matches) < 4 {
				continue
			}
			name := line[matches[2]:matches[3]]
			info[name] = protocol.Location{
				URI: doc,
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      protocol.UInteger(i),
						Character: protocol.UInteger(matches[2]),
					},
					End: protocol.Position{
						Line:      protocol.UInteger(i),
						Character: protocol.UInteger(matches[3]),
					},
				},
			}
		}
	}
	s.definitionsInfo[doc] = info
	return nil
}

// Returns: Location | []Location | []LocationLink | nil
func (s *Server) TextDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	docUri := params.TextDocument.URI
	posLine := int(params.Position.Line)
	posChar := int(params.Position.Character)

	defs := s.definitionsInfo[docUri]
	word := s.findWordUnderCursor(docUri, posLine, posChar)
	loc, ok := defs[word]
	if !ok {
		return nil, nil
	}
	return loc, nil
}

// TODO
