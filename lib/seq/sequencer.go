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

	currentVars assignments
	nextVars    assignments

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
	s.currentVars = s.nextVars
	s.nextVars = nil
	return nil
}

func (s *Sequencer) Assign(name string, value types.ValueFn) {
	s.nextVars = append(s.nextVars, assignment{name, value})
}

func (s *Sequencer) run() error {
	delay := time.Duration((15.0 / float64(s.tempo)) * float64(time.Second))
	var bit int64
	for {
		start := time.Now()
		ok, err := s.processFuncs(bit)
		if err != nil {
			log.Printf("File processing failed: %v", err)
		}
		dt := time.Since(start)
		currentDelay := delay - dt
		time.Sleep(currentDelay)

		if ok || bit > 0 {
			bit++
		}
	}
}

func (s *Sequencer) processFuncs(bit int64) (bool, error) {
	if len(s.current) == 0 {
		return false, nil
	}

	// eval variables first
	ctx := types.NewContext()
	for _, as := range s.currentVars {
		if err := ctx.Put(as.name, as.valueFn); err != nil {
			return false, err
		}
	}

	for _, fn := range s.current {
		ct := ctx.Copy()
		signals := fn.Eval(bit, ct)
		for _, sig := range signals {
			sg := sig
			s.ch <- &sg
		}
	}
	return true, nil
}

func (s *Sequencer) Close() error {
	return nil
}

var DefaultSequencer = NewSequencer()
