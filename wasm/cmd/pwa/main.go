package main

import (
	"log"

	"github.com/avoronkov/waver/wasm/lib/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func initializeComponents() {
	app.Route("/", &components.App{})

	app.RunWhenOnBrowser()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func doLog(format string, v ...any) {
	log.Printf(format, v...)
}
