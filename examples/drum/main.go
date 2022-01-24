package main

import (
	"log"
	"math"
	"os"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"
	"gitlab.com/avoronkov/waver/lib/midisynth/player"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

func main() {
	file := "../../samples/4-snare.wav"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	log.Printf("Playing file: %v", file)

	sample, err := waves.ReadSample(file)
	check(err)

	p := player.New(wav.Default)

	ctx := waves.NewNoteCtx(0.0, 0.0, math.Inf(1))

	reader, done := p.PlayContext(sample, ctx)
	check(err)

	c, ready, err := oto.NewContext(wav.Default.SampleRate, wav.Default.ChannelNum, wav.Default.BitDepthInBytes)
	check(err)
	<-ready

	pl := c.NewPlayer(reader)
	pl.Play()

	<-done
	time.Sleep(300 * time.Millisecond)
	runtime.KeepAlive(pl)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
