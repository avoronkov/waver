package components

import (
	"log"
	"strings"

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

	log    strings.Builder
	logger *log.Logger
}

func (a *App) Render() app.UI {
	return &Main{
		Log: &a.log,
	}
}

func (a *App) OnMount1(ctx app.Context) {
	ctx.Handle("play", a.handlePlay)
	ctx.Handle("pause", a.handlePause)
}

func (ap *App) handlePlay(ctx app.Context, a app.Action) {
	ap.mSequencer.Pause(false)
	if err := ap.mParser.ParseData([]byte(a.Value.(string))); err != nil {
		ap.doLog("Parse data failed: %v", err)
	}
}

func (ap *App) handlePause(ctx app.Context, a app.Action) {
	ap.mSequencer.Pause(true)
}

func (ap *App) doLog(format string, v ...any) {
	if ap.logger == nil {
		ap.logger = log.New(&ap.log, "", log.LstdFlags)
	}
	log.Printf(format, v...)
	ap.logger.Printf(format, v...)
}
