package unisynth

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/avoronkov/waver/lib/midisynth/multiplayer"
	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/notes"
)

type PlayCloser interface {
	Play(io.Reader) error
	Close() error
}

type Output struct {
	settings *wav.Settings
	play     *multiplayer.MultiPlayer

	player PlayCloser

	scale notes.Scale

	tempo int

	instruments InstrumentSet

	// Octave -> Note -> Release fn()
	notesReleases map[notes.Note]func()

	startTime time.Time

	wavFilename   string
	wavSpaceLeft  float64
	wavSpaceRight float64
	saver         *WavDataSaver
}

var _ signals.Output = (*Output)(nil)

type InstrumentSet interface {
	Wave(inst string) (waves.Wave, bool)
	// Sample(name string) (waves.Wave, bool)
	WaveControlled(inst string) (waves.WaveControlled, bool)
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

	// Init Player
	output.play = multiplayer.New(output.settings)

	// output.context = c
	var reader io.Reader
	if output.wavFilename != "" {
		log.Printf("Unisynth: using WavDataSaver")
		leftPad := secondsToFrames(output.settings, output.wavSpaceLeft)
		rightPad := secondsToFrames(output.settings, output.wavSpaceRight)
		output.saver = NewWavDataSaver(output.play, output.wavFilename, leftPad, rightPad)
		reader = output.saver
	} else {
		reader = output.play
	}

	go func() {
		if err := output.player.Play(reader); err != nil {
			slog.Error("Play failed", "error", err)
		}
	}()

	log.Printf("Unisynth initialized!")
	return output, nil
}

func (o *Output) ProcessAsync(tm float64, s signals.Interface) {
	switch a := s.(type) {
	case *signals.Signal:
		o.processSignal(tm, a)
	case *signals.Tempo:
		o.tempo = a.Tempo
	case *signals.StartTime:
		o.startTime = a.Start
	default:
		panic(fmt.Errorf("Unknown signal type: %v (%T)", s, s))
	}
}

func (o *Output) processSignal(tm float64, s *signals.Signal) {
	_ = tm
	at := s.Time.Add(1 * time.Second)

	absTime := float64(s.Time.Sub(o.startTime))/float64(time.Second) - 1.0
	var err error
	if !s.Manual {
		// Play note
		err = o.PlayNoteAt(at, s.Instrument, s.Note, s.DurationBits, s.Amp, absTime)
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
	if o.saver != nil {
		log.Printf("Unisynth: closing WavDataSaver")
		if err := o.saver.Close(); err != nil {
			return err
		}
	}
	if o.player != nil {
		log.Printf("Unisynth: [skipped] closing player")
		// TODO this hangs at the momemt.
		// if err := o.player.Close(); err != nil {
		// 	return err
		// }
	}
	return nil
}

func (o *Output) PlaySampleAt(at time.Time, name string, duration, amp, absTime float64) error {
	in, ok := o.instruments.Wave(name)
	if !ok {
		return fmt.Errorf("Unknown sample: %q", name)
	}
	o.play.AddWaveAt(at, in, waves.NewNoteCtx(0, amp, duration, absTime))

	return nil
}

func (o *Output) PlayNoteAt(at time.Time, instr string, note notes.Note, durationBits int, amp, absTime float64) error {
	freq := note.Freq
	dur := 15.0 * float64(durationBits) / float64(o.tempo)
	o.playNoteAt(at, instr, freq, dur, amp, absTime)
	return nil
}

func (o *Output) playNoteAt(at time.Time, inst string, hz, dur, amp, absTime float64) {
	in, ok := o.instruments.Wave(inst)
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	o.play.AddWaveAt(at, in, waves.NewNoteCtx(hz, amp, dur, absTime))
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

func (o *Output) PlayNoteControlled(inst string, note notes.Note, amp float64) (stop func(), err error) {
	panic("NIY")
	/*
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
	*/
}

/*
func (o *Output) SetTempo(tempo int) {
	o.tempo = tempo
}
*/

func secondsToFrames(settings *wav.Settings, seconds float64) int {
	return int(float64(settings.SampleRate) * seconds)
}
