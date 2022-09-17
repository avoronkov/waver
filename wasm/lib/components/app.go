package components

import (
	"log"

	"github.com/avoronkov/waver/lib/midisynth"
	"github.com/avoronkov/waver/lib/seq"
	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type App struct {
	app.Compo

	mParser    *parser.Parser
	mSequencer *seq.Sequencer
	mMidiSynth *midisynth.MidiSynth
}

func (a *App) Render() app.UI {
	return &Main{}
}

func (a *App) OnMount(ctx app.Context) {
	ctx.Handle("play", a.handlePlay)
}

func (ap *App) handlePlay(ctx app.Context, a app.Action) {
	log.Printf("handlePlay: %v", a.Value)
	ap.mSequencer.Pause(false)
	ap.mParser.ParseData([]byte(a.Value.(string)))
	log.Printf("handlePlay: OK")
}
