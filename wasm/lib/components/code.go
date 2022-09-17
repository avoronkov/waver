package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Code struct {
	app.Compo

	code string
}

func (c *Code) Render() app.UI {
	return app.Div().Class("card arert alert-success").Body(
		app.H3().Text("Code"),
		app.P().Body(
			app.Textarea().
				ID("code-story").
				Class("form-control").
				Name("code-story").
				OnInput(c.onInput),
		),
		app.P().Body(
			app.Button().
				ID("update-code").
				Type("button").
				Class("btn badge badge-success").
				OnClick(c.onPlay).
				Text("Play"),
		),
	)
}

func (c *Code) onInput(ctx app.Context, e app.Event) {
	c.code = ctx.JSSrc().Get("value").String()
}

func (c *Code) onPlay(ctx app.Context, e app.Event) {
	ctx.NewActionWithValue("play", c.code)
}
