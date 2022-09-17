package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Main struct {
	app.Compo
}

func (m *Main) Render() app.UI {
	return app.Main().Role("main").Body(
		app.Section().Class("container").Body(
			app.H3().Class("display-6").Text("Waver Playground (v2.5)"),
			app.P().Body(
				app.A().Href("https://github.com/avoronkov/waver").Text("Source code"),
			),
			&Code{},
		),
	)
}
