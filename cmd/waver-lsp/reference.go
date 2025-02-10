package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func (s *Server) WriteReference(out io.Writer) {
	fmt.Fprintf(out, "# Waver (pelia) language reference\n\n")

	s.pragmasReference(out)
}

func (s *Server) pragmasReference(out io.Writer) {
	fmt.Fprintf(out, "## Pragmas\n\n")
	pragmas := []string{}
	for pragma := range s.parser.PragmaParsers {
		pragmas = append(pragmas, pragma)
	}
	sort.Strings(pragmas)

	for _, pragma := range pragmas {
		meta := s.parser.PragmaParsers[pragma]
		fmt.Fprintf(out, "### %v\n", pragma)
		if meta.Deprecated {
			fmt.Fprintf(out, "**Deprecated**\n")
		}
		if meta.Usage != "" {
			fmt.Fprintf(out, "Usage: `%v`\n", meta.Usage)
		}
		if meta.Desc != "" {
			fmt.Fprintf(out, "\n%v\n", trimLeadingSpaces(meta.Desc))
		}
		if len(meta.Examples) > 0 {
			fmt.Fprintf(out, "\nExamples: `%v`\n", strings.Join(meta.Examples, "`, `"))
		}
		fmt.Fprintln(out)
	}
}
