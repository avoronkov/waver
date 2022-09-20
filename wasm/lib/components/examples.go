package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Examples struct {
	app.Compo

	// title -> url
	list []Pair[string, string]
}

func (e *Examples) Render() app.UI {
	return app.Div().Class("card arert alert-success").Body(
		app.Range(e.list).Slice(func(i int) app.UI {
			return app.P().Body(
				app.A().Href(e.list[i].Second).Text(e.list[i].First),
			)
		},
		),
	)
}

func (e *Examples) OnMount(ctx app.Context) {
	ctx.LocalStorage().Get("examples", &e.list)
	e.list = append(e.list, Pair[string, string]{"foobar", "http://example.com"})
	e.Update()
}

type Pair[T any, U any] struct {
	First  T
	Second U
}
