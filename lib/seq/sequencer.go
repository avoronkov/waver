package seq

import (
	"fmt"
	"log"
	"time"

	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

type Sequencer struct {
	tempo int
	bit   int64

	current []types.Signaler
	next    []types.Signaler

	currentVars assignments
	nextVars    assignments

	pause bool

	ch chan<- signals.Interface

	startingBit int64
	showBits    int64

	globalContext map[string]any
}

var _ signals.Input = (*Sequencer)(nil)

func NewSequencer(opts ...func(*Sequencer)) *Sequencer {
	s := &Sequencer{
		tempo:         120,
		globalContext: make(map[string]any),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *Sequencer) Run(ch chan<- signals.Interface) error {
	if s.ch == nil {
		panic("Channel should be set before Running sequencer")
	}
	log.Printf("Starting file sequencer...")
	if s.ch != ch {
		panic(fmt.Errorf("Cannels are differ: %v != %v", s.ch, ch))
	}
	return s.run()
}

func (s *Sequencer) Pause(v bool) {
	s.pause = v
}

func (s *Sequencer) Add(sig types.Signaler) {
	s.next = append(s.next, sig)
}

func (s *Sequencer) Commit() error {
	s.current = s.next
	s.next = nil
	s.currentVars = s.nextVars
	s.nextVars = nil
	return nil
}

func (s *Sequencer) Assign(name string, value types.ValueFn) {
	s.nextVars = append(s.nextVars, assignment{name, value})
}

func (s *Sequencer) SetTempo(tempo int) {
	s.tempo = tempo
	// Send to channel
	s.ch <- &signals.Tempo{
		Tempo: tempo,
	}
}

func (s *Sequencer) delay() time.Duration {
	return time.Duration((15.0 / float64(s.tempo)) * float64(time.Second))
}

func (s *Sequencer) run() error {
	// Skip until starting bit
	for s.bit < s.startingBit {
		_, err := s.processFuncs(time.Time{}, s.bit, true)
		if err != nil {
			log.Printf("File processing failed: %v", err)
		}
		s.bit++
	}

	start := time.Now()
	frameTime := start

	s.ch <- &signals.StartTime{Start: start}

	// Main loop
	for {
		var ok bool
		if !s.pause {
			if s.showBits > 0 && s.bit%s.showBits == 0 {
				go func(bit int64) {
					time.Sleep(1 * time.Second)
					log.Printf("Bit: %v", bit)
				}(s.bit)
			}

			var err error
			ok, err = s.processFuncs(frameTime, s.bit, false)
			if err != nil {
				log.Printf("File processing failed: %v", err)
			}
		}
		dt := time.Since(frameTime)
		currentDelay := s.delay() - dt
		time.Sleep(currentDelay)
		frameTime = frameTime.Add(s.delay())

		if !s.pause && (ok || s.bit > 0) {
			s.bit++
		}
	}
}

func (s *Sequencer) processFuncs(tm time.Time, bit int64, dryRun bool) (bool, error) {
	if len(s.current) == 0 {
		return false, nil
	}

	// eval variables first
	ctx := types.NewContext(types.WithGlobalContext(s.globalContext))
	// set default duration
	_ = ctx.Put("_dur", common.Const(4))
	for _, as := range s.currentVars {
		if err := ctx.Put(as.name, as.valueFn); err != nil {
			return false, err
		}
	}

	for _, fn := range s.current {
		ct := ctx.Copy()
		signals := fn.Eval(bit, ct)
		if !dryRun {
			for _, sig := range signals {
				sg := sig
				sg.Time = tm
				s.ch <- &sg
			}
		}
	}
	return true, nil
}

func (s *Sequencer) Close() error {
	return nil
}

var DefaultSequencer = NewSequencer()
