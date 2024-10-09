package main

import (
	"log"

	"github.com/avoronkov/waver/wasm/lib/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &components.App{})

	app.RunWhenOnBrowser()

	err := app.GenerateStaticWebsite("./waver", &app.Handler{
		Name:        "Waver",
		Description: "Waver playground application",
		Styles: []string{
			"https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css",
			"/web/lib/codemirror.css",
		},
		Scripts: []string{
			"/web/lib/codemirror.js",
			"/web/addon/mode/simple.js",
			"/web/init-codemirror.js",
			"/web/clipboard.js",
		},
		Resources: app.GitHubPages("waver"),
	})
	if err != nil {
		log.Fatal(err)
	}
}
