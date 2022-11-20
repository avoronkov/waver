package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

const (
	ShowCode     = "code"
	ShowMessages = "messages"
	ShowExamples = "examples"
)

type Main struct {
	app.Compo

	Log fmt.Stringer

	show string

	codeCompo     *Code
	messagesCompo *Messages
	examplesCompo *Examples
}

func (m *Main) Render() app.UI {
	return app.Main().Role("main").Body(
		app.Section().Class("container").Body(
			app.H3().Text("Waver Playground (v2.6)"),
			app.P().Body(
				app.A().Href("https://github.com/avoronkov/waver").Text("Source code"),
			),
			app.P().Body(
				app.Div().Class("card").Body(
					app.Select().Class("form-select").Body(
						app.Option().Value(ShowCode).Text("Code").Selected(m.show == "" || m.show == ShowCode),
						app.Option().Value(ShowMessages).Text("Messages").Selected(m.show == ShowMessages),
						app.Option().Value(ShowExamples).Text("Examples").Selected(m.show == ShowExamples),
					).OnChange(m.onChangeView),
				),
			),
			app.
				If(m.show == "" || m.show == ShowCode, m.codeCompo).
				ElseIf(m.show == ShowExamples, m.examplesCompo).
				Else(m.messagesCompo),
		),
	)
}

func (m *Main) OnMount(ctx app.Context) {
	ctx.Handle("onChangeView", m.handleChangeView)
}

func (m *Main) OnNav(ctx app.Context) {
	m.codeCompo = &Code{}
	m.messagesCompo = &Messages{Log: m.Log}
	m.examplesCompo = &Examples{}
}

func (m *Main) onChangeView(ctx app.Context, e app.Event) {
	m.codeCompo.Sync(ctx.LocalStorage())
	m.show = ctx.JSSrc().Get("value").String()
	m.Update()
}

func (m *Main) handleChangeView(ctx app.Context, a app.Action) {
	m.show = a.Value.(string)
	m.Update()
}
