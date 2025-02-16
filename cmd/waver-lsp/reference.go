package main

import (
	"fmt"
	"io"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/utils"
)

type descer interface {
	Desc() string
}

func (s *Server) WriteReference(out io.Writer) {
	fmt.Fprintf(out, "# Waver (pelia) language reference\n\n")

	s.pragmasReference(out)
	s.filtersReference(out)
}

func (s *Server) filtersReference(out io.Writer) {
	fmt.Fprintf(out, "## Filters\n\n")

	filtersNames := slices.Sorted(maps.Keys(filters.Filters))

	for _, name := range filtersNames {
		obj := filters.Filters[name]
		filterOptions := utils.NewSet[string]()

		fmt.Fprintf(out, "### %v\n\n", name)

		if d, ok := obj.(descer); ok {
			fmt.Fprintf(out, "%v\n\n", d.Desc())
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
			optInfo := tags[0]
			if otherTags := tags[1:]; len(otherTags) > 0 {
				optInfo += fmt.Sprintf(" (%v)", strings.Join(otherTags, ", "))
			}
			optInfo += ": " + fld.Type.String()
			filterOptions.Add(optInfo)
		}
		for _, opt := range slices.Sorted(slices.Values(filterOptions.Values())) {
			fmt.Fprintf(out, "- %v\n", opt)
		}
		fmt.Fprintln(out)
	}
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
		if len(meta.Examples) > 9 {
			fmt.Fprintf(out, "\nExamples:\n")
			for _, ex := range meta.Examples {

				fmt.Fprintf(out, "- `%v`\n", ex)
			}

		} else if len(meta.Examples) > 0 {
			fmt.Fprintf(out, "\nExamples: `%v`\n", strings.Join(meta.Examples, "`, `"))
		}
		fmt.Fprintln(out)
	}
}
