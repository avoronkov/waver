package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/avoronkov/waver/lib/seq/syntaxgen"
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
	params := syntaxgen.NewParams()

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
