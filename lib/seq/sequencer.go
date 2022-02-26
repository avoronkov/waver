package seq

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type Sequencer struct {
	tempo int

	current []types.Signaler
	next    []types.Signaler

	vars assignments

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

func (s *Sequencer) Assign(name string, value types.ValueFn) {
	s.vars = append(s.vars, assignment{name, value})
}

func (s *Sequencer) run() error {
	delay := time.Duration((15.0 / float64(s.tempo)) * float64(time.Second))
	currentDelay := 0 * time.Millisecond
	var bit int64
	for {
		select {
		case <-time.After(currentDelay):
			start := time.Now()
			if err := s.processFuncs(bit); err != nil {
				log.Printf("File processing failed: %v", err)
			}
			dt := time.Since(start)
			currentDelay = delay - dt
		}
		bit++
	}
}

func (s *Sequencer) processFuncs(bit int64) error {
	// eval variables first
	ctx := types.Context{}
	for _, as := range s.vars {
		if _, exists := ctx[as.name]; exists {
			return fmt.Errorf("Cannot re-assign variable: %v", as.name)
		}
		value := as.valueFn.Val(bit, ctx)
		ctx[as.name] = value
	}

	for _, fn := range s.current {

		signals := fn.Eval(bit, ctx)
		for _, sig := range signals {
			s.ch <- &sig
		}
	}
	return nil
}

func (s *Sequencer) Close() error {
	return nil
}

var DefaultSequencer = NewSequencer()
