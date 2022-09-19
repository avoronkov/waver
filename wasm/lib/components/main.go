package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

const (
	ShowCode     = "code"
	ShowMessages = "messages"
)

type Main struct {
	app.Compo

	Log fmt.Stringer

	show string

	codeCompo     *Code
	messagesCompo *Messages
}

func (m *Main) Render() app.UI {
	return app.Main().Role("main").Body(
		app.Section().Class("container").Body(
			app.H3().Class("display-6").Text("Waver Playground (v2.5)"),
			app.P().Body(
				app.A().Href("https://github.com/avoronkov/waver").Text("Source code"),
			),
			app.P().Body(
				app.Div().Class("card").Body(
					app.Select().Class("form-select").Body(
						app.Option().Value(ShowCode).Text("Code").Selected(m.show == "" || m.show == ShowCode),
						app.Option().Value(ShowMessages).Text("Messages").Selected(m.show == ShowMessages),
					).OnChange(m.onChangeView),
				),
			),
			app.
				If(m.show == "" || m.show == ShowCode, m.codeCompo).
				Else(m.messagesCompo),
		),
	)
}

func (m *Main) OnNav(ctx app.Context) {
	m.codeCompo = &Code{}
	m.messagesCompo = &Messages{Log: m.Log}
}

func (m *Main) onChangeView(ctx app.Context, e app.Event) {
	m.codeCompo.Sync()
	m.show = ctx.JSSrc().Get("value").String()
	m.Update()
}
