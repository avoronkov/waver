package synth

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"github.com/avoronkov/waver/lib/midisynth/player"
	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/notes"
)

type Output struct {
	settings *wav.Settings
	play     *player.Player

	context *oto.Context

	scale notes.Scale

	tempo int

	instruments InstrumentSet

	// Octave -> Note -> Release fn()
	notesReleases map[notes.Note]func()
}

var _ signals.Output = (*Output)(nil)

type InstrumentSet interface {
	Wave(inst int) (waves.Wave, bool)
	Sample(name string) (waves.Wave, bool)
	WaveControlled(inst int) (waves.WaveControlled, bool)
}

func New(opts ...func(*Output)) (*Output, error) {
	output := &Output{
		settings:      wav.Default,
		tempo:         120,
		notesReleases: make(map[notes.Note]func()),
	}
	for _, opt := range opts {
		opt(output)
	}

	// Init scale
	if output.scale == nil {
		output.scale = notes.NewStandard()
	}

	// Init oto.Context
	c, ready, err := oto.NewContext(
		output.settings.SampleRate,
		output.settings.ChannelNum,
		output.settings.BitDepthInBytes,
	)
	if err != nil {
		return nil, err
	}
	<-ready

	// Init Player
	output.play = player.New(output.settings)

	output.context = c

	return output, nil
}

func (o *Output) ProcessAsync(tm float64, s *signals.Signal) {
	at := s.Time.Add(20 * time.Millisecond)
	var err error
	if s.Sample != "" {
		// Play sample
		dur := 15.0 * float64(s.DurationBits) / float64(o.tempo)
		err = o.PlaySampleAt(at, s.Sample, dur, s.Amp)
	} else if !s.Manual {
		// Play note
		err = o.PlayNoteAt(at, s.Instrument, s.Note, s.DurationBits, s.Amp)
	} else if s.Stop {
		// Stop manual note
		o.releaseNote(s.Note)
	} else {
		// Play manual note
		stop, err := o.PlayNoteControlled(
			s.Instrument,
			s.Note,
			s.Amp,
		)
		if err != nil {
			log.Printf("[Manual] error: %v", err)
			return
		}
		o.storeNoteReleaseFn(s.Note, stop)
	}
	if err != nil {
		log.Printf("Error: %v", err)
	}

}

func (o *Output) Close() error {
	// oto.Context does not support method Close()
	return nil
}

func (o *Output) PlaySampleAt(at time.Time, name string, duration float64, amp float64) error {
	in, ok := o.instruments.Sample(name)
	if !ok {
		return fmt.Errorf("Unknown sample: %q", name)
	}
	data, done := o.play.PlayContext(in, waves.NewNoteCtx(0, amp, duration))

	p := o.context.NewPlayer(data)

	time.Sleep(at.Sub(time.Now()))
	p.Play()

	<-done
	time.Sleep(1 * time.Second)
	runtime.KeepAlive(p)
	return nil
}

func (o *Output) PlayNoteAt(at time.Time, instr int, note notes.Note, durationBits int, amp float64) error {
	freq := note.Freq
	dur := 15.0 * float64(durationBits) / float64(o.tempo)
	o.playNoteAt(at, instr, freq, dur, amp)
	return nil
}

func (o *Output) playNoteAt(at time.Time, inst int, hz float64, dur float64, amp float64) {
	in, ok := o.instruments.Wave(inst)
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	data, done := o.play.PlayContext(in, waves.NewNoteCtx(hz, amp, dur))

	p := o.context.NewPlayer(data)

	time.Sleep(at.Sub(time.Now()))
	p.Play()

	<-done
	time.Sleep(1 * time.Second)
	runtime.KeepAlive(p)
}

func (o *Output) storeNoteReleaseFn(note notes.Note, release func()) {
	o.notesReleases[note] = release
}

func (o *Output) releaseNote(note notes.Note) {
	if release, ok := o.notesReleases[note]; ok {
		release()
		delete(o.notesReleases, note)
	}
}

func (o *Output) PlayNoteControlled(inst int, note notes.Note, amp float64) (stop func(), err error) {
	wave, ok := o.instruments.WaveControlled(inst)
	if !ok {
		return nil, fmt.Errorf("Unknown instrument: %v", inst)
	}

	hz := note.Freq
	data, done := o.play.PlayContext(wave, waves.NewNoteCtx(hz, amp, math.Inf(+1)))

	go func() {
		p := o.context.NewPlayer(data)
		p.Play()

		<-done
		time.Sleep(1 * time.Second)
		runtime.KeepAlive(p)

	}()
	return wave.Release, nil
}
