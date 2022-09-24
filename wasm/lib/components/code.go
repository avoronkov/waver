package components

import (
	"log"

	"github.com/avoronkov/waver/lib/share"
	"github.com/avoronkov/waver/wasm/lib/storage"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Code struct {
	app.Compo

	text string

	saving   bool
	saveName string
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
			app.Text(" "),
			app.Button().
				ID("save").
				Type("button").
				Class("btn badge badge-primary").
				OnClick(c.onSave).
				Text("Save..."),

			app.If(c.saving, app.P().Body(
				app.Input().
					ID("save-input").
					Placeholder("Example name...").
					OnChange(c.onSaveInputChange),
				app.Text(" "),
				app.Button().
					ID("save-submit").
					Class("btn badge badge-primary").
					OnClick(c.onSaveSubmit).
					Text("Save"),
			)),
		),
	)
}

func (c *Code) OnMount(ctx app.Context) {
	queryCode, ok := ctx.Page().URL().Query()["code"]
	if ok {
		var err error
		c.text, err = share.Decode(queryCode[0])
		if err != nil {
			log.Printf("Failed to decode query code: %v", err)
		}
	}
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

func (c *Code) onSave(ctx app.Context, e app.Event) {
	c.saving = true
	c.Update()
}

func (c *Code) onSaveInputChange(ctx app.Context, e app.Event) {
	c.saveName = ctx.JSSrc().Get("value").String()
}

func (c *Code) onSaveSubmit(ctx app.Context, e app.Event) {
	defer func() {
		c.saveName = ""
		c.saving = false
		c.Update()
	}()

	c.Sync()
	encoded, err := share.Encode(c.text)
	if err != nil {
		log.Printf("Failed to encode example: %v", err)
		return
	}
	example := storage.Example{First: c.saveName, Second: encoded}
	err = storage.From(ctx.LocalStorage()).AddExample(example)
	if err != nil {
		log.Printf("Failed to save example: %v", err)
		return
	}
}
