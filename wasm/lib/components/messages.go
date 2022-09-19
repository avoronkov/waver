package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Messages struct {
	app.Compo

	Log fmt.Stringer
}

func (m *Messages) Render() app.UI {
	return app.Div().Class("card arert alert-success").Body(
		app.P().Body(
			app.Textarea().
				ReadOnly(true).
				Class("form-control").
				Text(m.Log.String()),
		),
	)
}
