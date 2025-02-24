package syntaxgen

import (
	"bufio"
	"bytes"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/avoronkov/waver/etc/std"
	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/avoronkov/waver/lib/utils"
)

type Params struct {
	parser            *parser.Parser
	Pragmas           []string
	Waves             []string
	StdFunctions      []string
	Functions         []string
	FunctionOperators []string
	Modifiers         []string
	ModifierOperators []string
	Identifiers       []string
	Filters           []string
	FilterOptions     []string
}

func NewParams() *Params {
	p := &Params{
		parser: parser.New(),
	}
	p.initPragmas()
	p.initWaves()
	p.initStdFunctions()
	p.initFunctions()
	p.initModifiers()
	p.initIdentifiers()
	p.initFilters()
	return p
}

func (p *Params) initPragmas() {
	p.Pragmas = slices.Sorted(maps.Keys(p.parser.PragmaParsers))
}

func (p *Params) initWaves() {
	p.Waves = slices.Sorted(maps.Keys(waves.Waves))
}

func (p *Params) initStdFunctions() {
	s := bufio.NewScanner(bytes.NewReader(std.StdEdo12))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 {
			p.StdFunctions = append(p.StdFunctions, fields[0])
		}
	}
	sort.Strings(p.StdFunctions)
}

func (p *Params) initFunctions() {
	for fn := range p.parser.FuncParsers {
		if _, ok := fn.(lexer.IdentToken); ok {
			p.Functions = append(p.Functions, fn.String())
		} else {
			p.FunctionOperators = append(p.FunctionOperators, fn.String())
		}
	}
	sort.Strings(p.Functions)
	sort.Strings(p.FunctionOperators)
}

func (p *Params) initModifiers() {
	for fn := range p.parser.ModParsers {
		if _, ok := fn.(lexer.IdentToken); ok {
			p.Modifiers = append(p.Modifiers, fn.String())
		} else {
			p.ModifierOperators = append(p.ModifierOperators, fn.String())
		}
	}
	sort.Strings(p.Modifiers)
	sort.Strings(p.ModifierOperators)
}

func (p *Params) initIdentifiers() {
	p.Identifiers = []string{"_", "_dur", "true", "false"}
}

func (p *Params) initFilters() {
	filterOptions := utils.NewSet[string]()
	for name, obj := range filters.Filters {
		p.Filters = append(p.Filters, name)

		v := reflect.TypeOf(obj)
		nField := v.NumField()
		for i := range nField {
			fld := v.Field(i)
			filterOptions.Add(strings.ToLower(fld.Name))
			tagsRaw := fld.Tag.Get("option")
			if tagsRaw == "" {
				continue
			}
			tags := strings.Split(tagsRaw, ",")
			filterOptions.Add(tags...)
		}
	}
	sort.Strings(p.Filters)
	p.FilterOptions = filterOptions.Values()
	sort.Strings(p.FilterOptions)
}
