package main

import (
	"log"
	"net/http"

	"github.com/avoronkov/waver/wasm/lib/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &components.Main{})

	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Waver",
		Description: "Waver playground application",
		Styles: []string{
			"https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css",
		},
	})

	log.Printf("Starting Waver Playgound on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func doLog(format string, v ...any) {
	log.Printf(format, v...)
}
