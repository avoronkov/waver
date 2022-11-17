//go:build !js

package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	initializeComponents()

	http.Handle("/", &app.Handler{
		Name:        "Waver",
		Description: "Waver playground application",
		Styles: []string{
			// "https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css",
			"/web/pico.min.css",
			"/web/lib/codemirror.css",
		},
		Scripts: []string{
			"/web/lib/codemirror.js",
			"/web/addon/mode/simple.js",
			"/web/init-codemirror.js",
		},
	})

	log.Printf("Starting Waver Playgound on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
