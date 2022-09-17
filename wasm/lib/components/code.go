package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Code struct {
	app.Compo

	text string
}

func (c *Code) Render() app.UI {
	return app.Div().Class("card arert alert-success").Body(
		app.P().Body(
			app.Textarea().
				ID("code-story").
				Class("form-control").
				Name("code-story"),
		),
		app.P().Body(
			app.Button().
				ID("update-code").
				Type("button").
				Class("btn badge badge-success").
				OnClick(c.onPlay).
				Text("Play"),
			app.Text(" "),
			app.Button().
				ID("stop-code").
				Type("button").
				Class("btn badge badge-danger").
				OnClick(c.onPause).
				Text("Stop"),
			app.Text(" "),
			app.Button().
				ID("clear-code").
				Type("button").
				Class("btn badge badge-danger").
				OnClick(c.onClear).
				Text("Clear"),
		),
	)
}

func (c *Code) OnMount(ctx app.Context) {
	app.Window().Call("initCodeMirror")
	app.Window().Call("setCodeMirrorCode", c.text)
}

func (c *Code) onPlay(ctx app.Context, e app.Event) {
	c.Sync()
	ctx.NewActionWithValue("play", c.text)
}

func (c *Code) onPause(ctx app.Context, e app.Event) {
	ctx.NewAction("pause")
}

func (c *Code) onClear(ctx app.Context, e app.Event) {
	app.Window().Call("setCodeMirrorCode", "")
}

func (c *Code) Sync() {
	c.text = app.Window().Call("getCodeMirrorCode").String()
}
