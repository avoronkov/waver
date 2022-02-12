package synth

import (
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"gitlab.com/avoronkov/waver/lib/midisynth/player"
	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
	"gitlab.com/avoronkov/waver/lib/notes"
)

type Output struct {
	settings *wav.Settings
	play     *player.Player

	context *oto.Context

	scale notes.Scale

	tempo int

	instruments InstrumentSet

	// Octave -> Note -> Release fn()
	notesReleases map[int]map[string]func()
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
		notesReleases: make(map[int]map[string]func()),
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

func (o *Output) ProcessAsync(s *signals.Signal) {
	var err error
	if s.Sample != "" {
		// Play sample
		dur := 15.0 * float64(s.DurationBits) / float64(o.tempo)
		err = o.PlaySample(s.Sample, dur, s.Amp)
	} else if !s.Manual {
		// Play note
		err = o.PlayNote(s.Instrument, s.Octave, s.Note, s.DurationBits, s.Amp)
	} else if s.Stop {
		// Stop manual note
		o.releaseNote(s.Octave, s.Note)
	} else {
		// Play manual note
		stop, err := o.PlayNoteControlled(
			s.Instrument,
			s.Octave,
			s.Note,
			s.Amp,
		)
		if err != nil {
			log.Printf("[Manual] error: %v", err)
			return
		}
		o.storeNoteReleaseFn(s.Octave, s.Note, stop)
	}
	if err != nil {
		log.Printf("Error: %v", err)
	}

}

func (o *Output) Close() error {
	return errors.New("NIY")
}

func (o *Output) PlaySample(name string, duration float64, amp float64) error {
	in, ok := o.instruments.Sample(name)
	if !ok {
		return fmt.Errorf("Unknown sample: %q", name)
	}
	data, done := o.play.PlayContext(in, waves.NewNoteCtx(0, amp, duration))

	p := o.context.NewPlayer(data)
	p.Play()

	<-done
	time.Sleep(1 * time.Second)
	runtime.KeepAlive(p)
	return nil
}

func (o *Output) PlayNote(instr int, octave int, note string, durationBits int, amp float64) error {
	freq, ok := o.scale.Note(octave, note)
	if !ok {
		return fmt.Errorf("Unknown note: %v%v", octave, note)
	}
	dur := 15.0 * float64(durationBits) / float64(o.tempo)
	o.playNote(instr, freq, dur, amp)
	return nil
}

func (o *Output) playNote(inst int, hz float64, dur float64, amp float64) {
	in, ok := o.instruments.Wave(inst)
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	data, done := o.play.PlayContext(in, waves.NewNoteCtx(hz, amp, dur))

	p := o.context.NewPlayer(data)
	p.Play()

	<-done
	time.Sleep(1 * time.Second)
	runtime.KeepAlive(p)
}

func (o *Output) storeNoteReleaseFn(octave int, note string, release func()) {
	if _, ok := o.notesReleases[octave]; ok {
		o.notesReleases[octave][note] = release
	} else {
		o.notesReleases[octave] = map[string]func(){
			note: release,
		}
	}
}

func (o *Output) releaseNote(octave int, note string) {
	log.Printf("releaseNote(%v, %v) : %+v", octave, note, o.notesReleases)
	if notes, ok := o.notesReleases[octave]; ok {
		if release, ok := notes[note]; ok {
			release()
			delete(notes, note)
		}
	}
}

func (o *Output) PlayNoteControlled(instr int, octave int, note string, amp float64) (stop func(), err error) {
	freq, ok := o.scale.Note(octave, note)
	if !ok {
		return nil, fmt.Errorf("Unknown note: %v%v", octave, note)
	}

	return o.playNoteControlled(instr, freq, amp)
}

func (o *Output) playNoteControlled(inst int, hz float64, amp float64) (stop func(), err error) {
	wave, ok := o.instruments.WaveControlled(inst)
	if !ok {
		return nil, fmt.Errorf("Unknown instrument: %v", inst)
	}
	log.Printf("playNoteControlled: wave, ok = %v, %v", wave, ok)

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
