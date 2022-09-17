package components

import (
	"log"

	"github.com/avoronkov/waver/lib/midisynth"
	"github.com/avoronkov/waver/lib/seq"
	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Main struct {
	app.Compo

	mParser    *parser.Parser
	mSequencer *seq.Sequencer
	mMidiSynth *midisynth.MidiSynth
}

func (m *Main) Render() app.UI {
	return app.Main().Role("main").Body(
		app.Section().Class("container").Body(
			app.H3().Class("display-6").Text("Waver Playground (v2.5)"),
			app.P().Body(
				app.A().Href("https://github.com/avoronkov/waver").Text("Source code"),
			),
			&Code{},
		),
	)
}

func (m *Main) OnMount(ctx app.Context) {
	ctx.Handle("play", m.handlePlay)
}

func (m *Main) handlePlay(ctx app.Context, a app.Action) {
	log.Printf("handlePlay: %v", a.Value)
	m.mSequencer.Pause(false)
	m.mParser.ParseData([]byte(a.Value.(string)))
	log.Printf("handlePlay: OK")
}
