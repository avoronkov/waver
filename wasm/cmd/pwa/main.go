package main

import (
	"github.com/avoronkov/waver/wasm/lib/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func initializeComponents() {
	app.Route("/", &components.App{})

	app.RunWhenOnBrowser()
}
