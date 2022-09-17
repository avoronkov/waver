package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Messages struct {
	app.Compo
}

func (m *Messages) Render() app.UI {
	return app.Div().Class("card arert alert-success").Body(
		app.P().Body(
			app.Span().Text("Some messages"),
		),
	)
}
