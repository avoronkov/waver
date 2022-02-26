package seq

import (
	"log"
	"time"

	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type Sequencer struct {
	tempo int

	current []types.Signaler
	next    []types.Signaler

	ch chan<- *signals.Signal
}

var _ signals.Input = (*Sequencer)(nil)

func NewSequencer() *Sequencer {
	s := &Sequencer{
		tempo: 120,
	}
	return s
}

func (s *Sequencer) Start(ch chan<- *signals.Signal) error {
	log.Printf("Starting file sequencer...")
	s.ch = ch
	go func() {
		if err := s.run(); err != nil {
			log.Printf("[ERROR] Sequencer failed: %v", err)
		}
	}()
	return nil
}

func (s *Sequencer) Add(sig types.Signaler) {
	s.next = append(s.next, sig)
}

func (s *Sequencer) Commit() error {
	s.current = s.next
	s.next = nil
	return nil
}

func (s *Sequencer) run() error {
	delay := time.Duration((15.0 / float64(s.tempo)) * float64(time.Second))
	currentDelay := 0 * time.Millisecond
	var bit int64
	for {
		select {
		case <-time.After(currentDelay):
			start := time.Now()
			s.processFuncs(bit, s.current)
			dt := time.Since(start)
			currentDelay = delay - dt
		}
		bit++
	}
}

func (s *Sequencer) processFuncs(bit int64, funcs []types.Signaler) {
	for _, fn := range funcs {
		signals := fn.Eval(bit, types.Context{})
		for _, sig := range signals {
			log.Printf("[Sequencer] sending signal %v ...", sig)
			s.ch <- &sig
			log.Printf("[Sequencer] sending signal %v DONE", sig)
		}
	}
}

func (s *Sequencer) Close() error {
	return nil
}

var DefaultSequencer = NewSequencer()
