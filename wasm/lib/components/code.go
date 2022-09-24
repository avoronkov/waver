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
				Class("btn badge badge-warning").
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
				Class("btn badge badge-secondary").
				OnClick(c.onSave).
				Text("Save..."),
		),
	)
}

func (c *Code) OnMount(ctx app.Context) {
	log.Printf("OnMount")
	queryCode, ok := ctx.Page().URL().Query()["code"]
	if ok {
		var err error
		c.text, err = share.Decode(queryCode[0])
		if err != nil {
			log.Printf("Failed to decode query code: %v", err)
		}
	} else {
		var err error
		c.text, err = storage.From(ctx.LocalStorage()).GetCode()
		if err != nil {
			log.Printf("Failed to get code from local storage: %v", err)
		}
		log.Printf("Got code from localstorage: %v", c.text)
	}
	app.Window().Call("initCodeMirror")
	app.Window().Call("setCodeMirrorCode", c.text)
}

func (c *Code) OnNav(ctx app.Context) {
	log.Printf("OnNav")
}

func (c *Code) onPlay(ctx app.Context, e app.Event) {
	c.Sync(ctx.LocalStorage())
	ctx.NewActionWithValue("play", c.text)
}

func (c *Code) onPause(ctx app.Context, e app.Event) {
	ctx.NewAction("pause")
}

func (c *Code) onClear(ctx app.Context, e app.Event) {
	ok := app.Window().Call("confirm", "Clear code?").Bool()
	if !ok {
		return
	}
	app.Window().Call("setCodeMirrorCode", "")
	c.text = ""
}

func (c *Code) Sync(st app.BrowserStorage) {
	c.text = app.Window().Call("getCodeMirrorCode").String()
	log.Printf("Saving code into localstorage: %v", c.text)
	if err := storage.From(st).SaveCode(c.text); err != nil {
		log.Printf("Failed to save code into local storage: %v", err)
	}
}

func (c *Code) onSave(ctx app.Context, e app.Event) {
	name := app.Window().Call("prompt", "Saving example...", "Example name...")
	if name.IsNull() {
		return
	}

	defer func() {
		c.Update()
	}()

	c.Sync(ctx.LocalStorage())
	encoded, err := share.Encode(c.text)
	if err != nil {
		log.Printf("Failed to encode example: %v", err)
		return
	}
	example := storage.Example{First: name.String(), Second: encoded}
	err = storage.From(ctx.LocalStorage()).AddExample(example)
	if err != nil {
		log.Printf("Failed to save example: %v", err)
		return
	}
}
