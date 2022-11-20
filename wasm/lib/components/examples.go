package components

import (
	"log"

	"github.com/avoronkov/waver/wasm/lib/storage"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Examples struct {
	app.Compo

	// title -> url
	list []storage.Example
}

func (e *Examples) Render() app.UI {
	return app.Div().Class("card arert alert-success").Body(
		app.Range(e.list).Slice(func(i int) app.UI {
			return app.P().Body(
				app.A().
					Href("/?code="+e.list[i].Second).
					OnClick(e.onClickExample).
					Text(e.list[i].First),
				app.Text(" "),
				app.A().OnClick(e.onDelete(i), i).Text("[del]"),
			)
		},
		),
	)
}

func (e *Examples) OnMount(ctx app.Context) {
	var err error
	e.list, err = storage.From(ctx.LocalStorage()).GetExamples()
	if err != nil {
		log.Printf("Failed to get examples: %v", err)
		return
	}
	e.Update()
}

func (ex *Examples) onDelete(index int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		ok := app.Window().Call("confirm", "Delete example?").Bool()
		if !ok {
			return
		}

		err := storage.From(ctx.LocalStorage()).DelExample(index)
		if err != nil {
			log.Printf("Failed to delete example: %v", err)
			return
		}
		ex.list, err = storage.From(ctx.LocalStorage()).GetExamples()
		if err != nil {
			log.Printf("Failed to get examples: %v", err)
			return
		}
		ex.Update()
	}
}

func (e *Examples) onClickExample(ctx app.Context, event app.Event) {
	ctx.NewActionWithValue("onChangeView", "")
}
