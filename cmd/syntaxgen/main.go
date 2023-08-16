package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/avoronkov/waver/etc/std"
	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/avoronkov/waver/lib/utils"
)

//go:embed pelia.vim
var peliaVim string

//go:embed surfer.vim
var surferVim string

//go:embed init-codemirror.js
var initCodemirrorJs string

func processTemplate(params any, name, tpl, file string) error {
	t := template.Must(template.New(name).Funcs(template.FuncMap{
		"stringsJoin": strings.Join,
	}).Parse(tpl))

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	return t.Execute(f, params)
}

func findFile(path string) string {
	for _, sub := range []string{"../..", "."} {
		p := filepath.Join(sub, path)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	panic(fmt.Errorf("File not found: %v", path))
}

func main() {
	params := NewParams()

	log.Println("Generate pelia.vim...")
	err := processTemplate(params, "pelia.vim", peliaVim, findFile("tools/vim/syntax/pelia.vim"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Generate init-codemirror.js...")
	err = processTemplate(params, "init-codemirror.js", initCodemirrorJs, findFile("wasm/web/init-codemirror.js"))
	if err != nil {
		log.Fatal(err)
	}

	surferParams := InitSurferParams()
	log.Printf("Generate surfer.vim...")
	err = processTemplate(surferParams, "surfer.vim", surferVim, findFile("tools/vim/syntax/surfer.vim"))

	log.Printf("OK")
}

type Params struct {
	parser            *parser.Parser
	Pragmas           []string
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
	p.initStdFunctions()
	p.initFunctions()
	p.initModifiers()
	p.initIdentifiers()
	p.initFilters()
	return p
}

func (p *Params) initPragmas() {
	for pragma := range p.parser.PragmaParsers {
		p.Pragmas = append(p.Pragmas, pragma)
	}
	sort.Strings(p.Pragmas)
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
		for i := 0; i < nField; i++ {
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
