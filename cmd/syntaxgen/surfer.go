package main

import (
	"sort"
	"strings"

	"github.com/avoronkov/waver/lib/forth/parser"
	"github.com/avoronkov/waver/lib/surfer"
)

type SurferParams struct {
	Pragmas      []string
	StdFunctions []string
	Functions    []string
}

func InitSurferParams() *SurferParams {
	p := &SurferParams{}

	p.Pragmas = []string{"define"}

	p.initStdFunctions()
	p.initFunctions()
	return p
}

func (p *SurferParams) initStdFunctions() {
	for fn := range parser.Funcs {
		p.StdFunctions = append(p.StdFunctions, fn)
	}
	sort.Strings(p.StdFunctions)
}

func (p *SurferParams) initFunctions() {
	in := surfer.NewInterpreter()
	for fn := range in.Functions {
		if strings.ContainsAny(fn, "|") {
			continue
		}
		p.Functions = append(p.Functions, fn)
	}
	sort.Strings(p.Functions)
}
